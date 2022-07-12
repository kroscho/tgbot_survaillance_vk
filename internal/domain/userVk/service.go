package userVk

import (
	"context"
	"tgbot_surveillance/internal/domain/user"
)

type Service interface {
	Create(ctx context.Context, from *UserVk) error
}

type service struct {
	userService user.Service
	store       Store
}

func NewService(userService user.Service, store Store) Service {
	return &service{userService: userService, store: store}
}

func (s service) Create(ctx context.Context, from *UserVk) error {
	return s.store.Create(ctx, from)
}
