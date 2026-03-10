package vone

import (
	"context"
	"errors"
	"iter"
	"time"
)

type paginatedRequest[T any] interface {
	Do(ctx context.Context) (*T, error)
	nextLink() string
	resetPagination()
	isDone(*T) bool
}

type Paginator[T any, Item any] struct {
	req   paginatedRequest[T]
	items func(*T) []Item
}

func NewPaginator[T any, Item any](
	req paginatedRequest[T],
	items func(*T) []Item,
) *Paginator[T, Item] {
	req.resetPagination()
	return &Paginator[T, Item]{
		req:   req,
		items: items,
	}
}

func (p *Paginator[T, Item]) Range(ctx context.Context) iter.Seq2[*Item, error] {
	return func(yield func(*Item, error) bool) {
		for {
			resp, err := p.req.Do(ctx)
			if err != nil {
				if rl, ok := IsRateLimit(err); ok {
					timer := time.NewTimer(time.Duration(rl.Reset) * time.Second)
					select {
					case <-ctx.Done():
						timer.Stop()
						yield(nil, ctx.Err())
						return
					case <-timer.C:
					}
					continue
				}

				yield(nil, err)
				return
			}

			items := p.items(resp)
			for i := range items {
				item := items[i]
				if !yield(&item, nil) {
					return
				}
			}

			if p.req.isDone(resp) {
				return
			}

			// защита от кривого ответа API:
			// если ещё не done, но nextLink пустой, это тупик
			if p.req.nextLink() == "" {
				yield(nil, errors.New("pagination stalled: not done but nextLink is empty"))
				return
			}
		}
	}
}
