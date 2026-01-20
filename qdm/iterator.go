package qdm

import (
	"context"
	"sync"
)

type Iterator interface {
	Fetch(batch int) ([]any, error)
	Count() int64
	Close(err error)
	Done() <-chan struct{}
	Error() error
}

type iterator struct {
	count     int64
	cursor    int64
	ch        chan any
	errCh     <-chan error
	ctx       context.Context
	cancel    context.CancelCauseFunc
	closeOnce sync.Once
}

func (it *iterator) Fetch(batch int) ([]any, error) {
	if it.cursor >= it.count {
		return nil, EOF
	}

	items := make([]any, 0)
	for item := range it.ch {
		items = append(items, item)
		it.cursor++

		if it.cursor >= it.count {
			return items, nil
		}

		if len(items) == batch {
			return items, nil
		}
	}

	// channel closed
	if len(items) == 0 {
		return nil, EOF
	}

	return items, nil
}

func (it *iterator) Count() int64 {
	return it.count
}

func (it *iterator) Close(err error) {
	it.closeOnce.Do(func() {
		if it.cancel != nil {
			it.cancel(err)
		}

		if it.ch != nil {
			close(it.ch)
		}
	})
}

func (it *iterator) Done() <-chan struct{} {
	return it.ctx.Done()
}

func (it *iterator) Error() error {
	return it.ctx.Err()
}

func (it *iterator) handle(ctx context.Context, errCh <-chan error) {
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
