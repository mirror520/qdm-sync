package qdm

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/mirror520/qdm-sync/orders"
)

type Service interface {
	Authorize(id string, secret string) (*AuthData, error)
	OrderCount(start time.Time, end time.Time, opts ...OrderOption) (int64, error)
	FindOrders(start time.Time, end time.Time, opts ...OrderOption) (Iterator, error)
	FindCustomerGroups() ([]*orders.CustomerGroup, error)
	Close()
}

func NewService(cfg Config) (Service, error) {
	log := zap.L().With(
		zap.String("service", "qdm"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	svc := &service{
		cfg: cfg,
		log: log,
		client: resty.New().
			SetBaseURL("https://" + cfg.BaseURL + "/api/v1").
			SetAllowGetMethodPayload(true),
		ctx:    ctx,
		cancel: cancel,
	}

	auth, err := svc.Authorize(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, err
	}

	go svc.refreshToken(ctx, auth)

	svc.token = auth.AccessToken
	log.Info("authorized")

	return svc, nil
}

type service struct {
	cfg    Config
	log    *zap.Logger
	client *resty.Client
	token  string
	ctx    context.Context
	cancel context.CancelFunc
}

func (svc *service) Authorize(id string, secret string) (*AuthData, error) {
	var result Result
	resp, err := svc.client.R().
		SetBasicAuth(id, secret).
		SetResult(&result).
		SetError(&Result{}).
		ForceContentType("application/json").
		Post("/token/authorize")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		result, ok := resp.Error().(*Result)
		if !ok {
			return nil, errors.New(resp.String())
		}

		return nil, result.Error()
	}

	return result.AuthData()
}

const renewAfter float64 = 0.75

func (svc *service) refreshToken(ctx context.Context, auth *AuthData) {
	log := svc.log.With(
		zap.String("action", "check_token"),
		zap.Time("expired_at", auth.ExpiresIn),
	)

	renewTS := time.Until(auth.ExpiresIn).Seconds() * renewAfter
	renewDuration := time.Duration(renewTS) * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Info("done")
			return

		case <-time.After(renewDuration):
			for {
				auth, err := svc.Authorize(svc.cfg.ClientID, svc.cfg.ClientSecret)
				if err != nil {
					log.Error(err.Error())

					time.Sleep(3000 * time.Millisecond)
					continue
				}

				go svc.refreshToken(ctx, auth)

				svc.token = auth.AccessToken
				log.Info("token refreshed")
				return
			}
		}
	}
}

func (svc *service) OrderCount(start time.Time, end time.Time, opts ...OrderOption) (int64, error) {
	params := &OrderParams{
		CreatedAtMin: start,
		CreatedAtMax: end,
	}

	for _, opt := range opts {
		opt.apply(params)
	}

	var result Result

	resp, err := svc.client.R().
		SetAuthToken(svc.token).
		SetFormDataFromValues(params.Values()).
		SetResult(&result).
		SetError(&Result{}).
		ForceContentType("application/json").
		Get("/orders/count")

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() != http.StatusOK {
		result, ok := resp.Error().(*Result)
		if !ok {
			return 0, errors.New(resp.String())
		}

		return 0, result.Error()
	}

	data, err := result.OrderCountData()
	if err != nil {
		return 0, err
	}

	return data.Count, nil
}

func (svc *service) FindOrders(start time.Time, end time.Time, opts ...OrderOption) (Iterator, error) {
	params := &OrderParams{
		CreatedAtMin: start,
		CreatedAtMax: end,
		PageSize:     300,
		PageNumber:   1,
	}

	for _, opt := range opts {
		opt.apply(params)
	}

	count, err := svc.OrderCount(start, end, opts...)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("empty data")
	}

	ch := make(chan orders.Order, params.PageNumber/2)
	errCh := make(chan error)
	go func(ch chan<- orders.Order, errCh chan<- error) {
		defer func() {
			recover()
			errCh <- errors.New("force close")
		}()

		defer func(ch chan<- orders.Order) {
			for {
				if len(ch) == 0 {
					close(ch)
					return
				}

				time.Sleep(500 * time.Millisecond)
			}
		}(ch)

		for {
			var result Result

			resp, err := svc.client.R().
				SetAuthToken(svc.token).
				SetFormDataFromValues(params.Values()).
				SetResult(&result).
				SetError(&Result{}).
				ForceContentType("application/json").
				Get("/orders")

			if err != nil {
				errCh <- err
				return
			}

			if resp.StatusCode() != http.StatusOK {
				result, ok := resp.Error().(*Result)
				if !ok {
					errCh <- errors.New(resp.String())
					return
				}

				errCh <- result.Error()
				return
			}

			data, err := result.OrderData()
			if err != nil {
				errCh <- err
				return
			}

			if data.Count == 0 {
				errCh <- EOF
				return
			}

			for _, o := range data.Result {
				ch <- o
			}

			sc := data.SearchCriteria
			if sc.PageNumber == sc.PageCount {
				errCh <- EOF
				return
			}

			params.PageNumber++
		}
	}(ch, errCh)

	ctx, cancel := context.WithCancelCause(svc.ctx)
	it := &iterator{
		count:  count,
		ch:     ch,
		errCh:  errCh,
		ctx:    ctx,
		cancel: cancel,
	}

	go it.handle(ctx, errCh)

	return it, nil
}

func (svc *service) Close() {
	if svc.cancel != nil {
		svc.cancel()
	}

	svc.cancel = nil
}

type Iterator interface {
	Fetch(batch int) ([]orders.Order, error)
	Count() int64
	Close(err error)
	Done() <-chan struct{}
	Error() error
}

type iterator struct {
	count  int64
	ch     chan orders.Order
	errCh  <-chan error
	ctx    context.Context
	cancel context.CancelCauseFunc
}

func (it *iterator) Fetch(batch int) ([]orders.Order, error) {
	orders := make([]orders.Order, 0)

	for order := range it.ch {
		orders = append(orders, order)

		if len(orders) == batch {
			return orders, nil
		}
	}

	// channel closed
	if len(orders) == 0 {
		return nil, EOF
	}

	return orders, nil
}

func (it *iterator) Count() int64 {
	return it.count
}

func (it *iterator) Close(err error) {
	if it.cancel != nil {
		it.cancel(err)
	}

	if it.ch != nil && len(it.ch) > 0 {
		close(it.ch)
	}

	it.ch = nil
	it.cancel = nil
}

func (it *iterator) Done() <-chan struct{} {
	return it.ctx.Done()
}

func (it *iterator) Error() error {
	return it.ctx.Err()
}

func (it iterator) handle(ctx context.Context, errCh <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return

		case err := <-errCh:
			it.Close(err)
			return
		}
	}
}

func (svc *service) FindCustomerGroups() ([]*orders.CustomerGroup, error) {
	var result Result

	resp, err := svc.client.R().
		SetAuthToken(svc.token).
		SetResult(&result).
		SetError(&Result{}).
		ForceContentType("application/json").
		Get("/customers/group")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		result, ok := resp.Error().(*Result)
		if !ok {
			return nil, errors.New(resp.String())
		}

		return nil, result.Error()
	}

	data, err := result.CustomerGroupData()
	if err != nil {
		return nil, err
	}

	return data.Result, nil
}
