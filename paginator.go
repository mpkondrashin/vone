package vone

import (
	"context"
	"iter"
)

type paginatedRequest[T any] interface {
	Do(ctx context.Context) (*T, error)
	nextLink() string
	resetPagination()
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
	return &Paginator[T, Item]{req: req, items: items}
}

func (p *Paginator[T, Item]) Range(
	ctx context.Context,
) iter.Seq2[*Item, error] {
	return func(yield func(*Item, error) bool) {
		for {
			resp, err := p.req.Do(ctx)
			if err != nil {
				yield(nil, err)
				return
			}

			for i := range p.items(resp) {
				if !yield(&p.items(resp)[i], nil) {
					return
				}
			}

			if p.req.nextLink() == "" {
				return
			}
		}
	}
}
