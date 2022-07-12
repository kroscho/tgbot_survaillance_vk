package user

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, opts ...GetOptFunc) ([]*User, error)
	Update(ctx context.Context, user *User) error
	GetUserByTgID(ctx context.Context, tgID TelegramID) (*User, error)
}

type service struct {
	store Store
}

func NewService(store Store) Service {
	return &service{store: store}
}

func (s service) Create(ctx context.Context, user *User) (*User, error) {
	return s.store.Create(ctx, user)
}

func (s service) Get(ctx context.Context, opts ...GetOptFunc) ([]*User, error) {
	return s.store.Get(ctx, opts...)
}

func (s service) Update(ctx context.Context, user *User) error {
	return s.store.Update(ctx, user)
}

func (s service) GetUserByTgID(ctx context.Context, tgID TelegramID) (*User, error) {
	return s.store.GetUserByTgID(ctx, tgID)
}
