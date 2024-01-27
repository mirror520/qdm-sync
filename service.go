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
	Sync(start time.Time, end time.Time) (<-chan Progress, int64, error)
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

func (svc *service) Sync(start time.Time, end time.Time) (<-chan Progress, int64, error) {
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
				orders, err := it.Fetch(10)
				if err != nil {
					if errors.Is(err, qdm.EOF) {
						it.Close(nil)

						log.Info(err.Error())
						return
					}

					log.Error(err.Error())
					return
				}

				if err := svc.orders.Store(orders); err != nil {
					log.Error(err.Error())
					return
				}

				progress.Current += int64(len(orders))

				ch <- progress
			}
		}
	}(svc.ctx, it, ch)

	return ch, progress.Total, nil
}

func (svc *service) Close() {
	if svc.cancel != nil {
		svc.cancel()
	}

	svc.cancel = nil
}
