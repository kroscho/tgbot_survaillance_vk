package user

import (
	"context"
)

type Store interface {
	Create(context.Context, *User) (*User, error)
	Get(ctx context.Context, opts ...GetOptFunc) ([]*User, error)
	Update(ctx context.Context, user *User) error
	GetUserByTgID(ctx context.Context, tgID TelegramID) (*User, error)
	Delete(ctx context.Context, user *User) error
}

type Filters struct {
	IDs   []ID
	TgIDs []TelegramID
}

type Lock struct {
	ForUpdate  bool
	SkipLocked bool
}

type GetOptFunc func(*GetOptions)

type GetOptions struct {
	Filters Filters
	Sort    []Sort
	Lock    Lock
	Limit   uint
	Offset  *uint
}

func NewGetOptions(opts ...GetOptFunc) *GetOptions {
	o := &GetOptions{
		Sort: []Sort{},
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

type SortBy string

type Sort struct {
	By   SortBy
	Desc bool
}
