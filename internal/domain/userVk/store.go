package userVk

import (
	"context"
)

type Store interface {
	Create(ctx context.Context, from *UserVk) error
}

type Filters struct {
	VkIDs []VkID
}

type GetOptFunc func(*GetOptions)

type GetOptions struct {
	Filters Filters
}

func NewGetOptions(opts ...GetOptFunc) *GetOptions {
	o := &GetOptions{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithFilters(filters Filters) GetOptFunc {
	return func(o *GetOptions) {
		o.Filters = filters
	}
}
