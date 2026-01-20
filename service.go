package sync

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/mirror520/qdm-sync/orders"
	"github.com/mirror520/qdm-sync/qdm"
)

type Service interface {
	SyncOrders(start time.Time, end time.Time) (<-chan Progress, int64, error)
	SyncCustomers(start time.Time, end time.Time) (<-chan Progress, int64, error)
	SyncCustomerGroups() (int, error)
	Close()
}

func NewService(qdm qdm.Service, repo orders.Repository) Service {
	ctx, cancel := context.WithCancel(context.Background())
	return &service{
		log: zap.L().With(
			zap.String("service", "qdm-sync"),
		),
		qdm:    qdm,
		orders: repo,
		ctx:    ctx,
		cancel: cancel,
	}
}

type service struct {
	log    *zap.Logger
	qdm    qdm.Service
	orders orders.Repository
	ctx    context.Context
	cancel context.CancelFunc
}

type Progress struct {
	Total   int64
	Current int64
}

func (svc *service) SyncOrders(start time.Time, end time.Time) (<-chan Progress, int64, error) {
	it, err := svc.qdm.FindOrders(start, end)
	if err != nil {
		return nil, 0, err
	}

	progress := Progress{
		Total:   it.Count(),
		Current: 0,
	}

	ch := make(chan Progress)
	go func(ctx context.Context, it qdm.Iterator, ch chan<- Progress) {
		log := svc.log.With(
			zap.String("action", "sync"),
			zap.String("entity", "orders"),
			zap.Int64("count", it.Count()),
		)

		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				log.Info("done")
				return

			case <-it.Done():
				log.Info("done")
				return

			case <-ticker.C:
				items, err := it.Fetch(10)
				if err != nil {
					if errors.Is(err, qdm.EOF) {
						it.Close(nil)

						log.Info(err.Error())
						return
					}

					log.Error(err.Error())
					return
				}

				newOrders := make([]orders.Order, len(items))
				for i, item := range items {
					order, ok := item.(orders.Order)
					if !ok {
						log.Error("type assertion failed")
						return
					}

					newOrders[i] = order
				}

				if err := svc.orders.Store(newOrders); err != nil {
					log.Error(err.Error())
					return
				}

				progress.Current += int64(len(newOrders))

				ch <- progress
			}
		}
	}(svc.ctx, it, ch)

	return ch, progress.Total, nil
}

func (svc *service) SyncCustomers(start time.Time, end time.Time) (<-chan Progress, int64, error) {
	it, err := svc.qdm.FindCustomers(start, end)
	if err != nil {
		return nil, 0, err
	}

	progress := Progress{
		Total:   it.Count(),
		Current: 0,
	}

	ch := make(chan Progress)
	go func(ctx context.Context, it qdm.Iterator, ch chan<- Progress) {
		log := svc.log.With(
			zap.String("action", "sync"),
			zap.String("entity", "customers"),
			zap.Int64("count", it.Count()),
		)

		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				log.Info("done")
				return

			case <-it.Done():
				log.Info("done")
				return

			case <-ticker.C:
				items, err := it.Fetch(10)
				if err != nil {
					if errors.Is(err, qdm.EOF) {
						it.Close(nil)

						log.Info(err.Error())
						return
					}

					log.Error(err.Error())
					return
				}

				newCustomers := make([]orders.Customer, len(items))
				for i, item := range items {
					customer, ok := item.(orders.Customer)
					if !ok {
						log.Error("type assertion failed")
						return
					}

					newCustomers[i] = customer
				}

				if err := svc.orders.StoreCustomers(newCustomers); err != nil {
					log.Error(err.Error())
					return
				}

				progress.Current += int64(len(newCustomers))

				ch <- progress
			}
		}
	}(svc.ctx, it, ch)

	return ch, progress.Total, nil
}

func (svc *service) SyncCustomerGroups() (int, error) {
	groups, err := svc.qdm.FindCustomerGroups()
	if err != nil {
		return 0, err
	}

	if err := svc.orders.StoreCustomerGroups(groups); err != nil {
		return 0, err
	}

	return len(groups), nil
}

func (svc *service) Close() {
	if svc.cancel != nil {
		svc.cancel()
	}

	svc.cancel = nil
}
