package tracked

import (
	"context"
	"tgbot_surveillance/internal/domain/user"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"
)

type Service interface {
	Get(ctx context.Context, user *user.User) ([]*TrackedInfo, error)
	Create(ctx context.Context, user *user.User, userAdd *vkmodels.User) error
	GetPrevFriends(ctx context.Context, tracked *TrackedInfo) ([]int64, error)
	UpdatePrevFriends(ctx context.Context, tracked *TrackedInfo, newFriends map[int64]vkmodels.User) error
	GetTrackedByVkID(ctx context.Context, user *user.User, vkId int64) (*TrackedInfo, error)
	DeleteUserFromPrevFriends(ctx context.Context, deleteUser *vkmodels.User, tracked *TrackedInfo) error
	AddUserInPrevFriends(ctx context.Context, addUser *vkmodels.User, tracked *TrackedInfo) error
	DeleteUserFromTracked(ctx context.Context, user *user.User, tracked *TrackedInfo) error
	GetTrackeds(ctx context.Context) (map[VkID][]*UserTrackedInfo, error)
}

type service struct {
	userService user.Service
	store       Store
}

func NewService(userService user.Service, store Store) Service {
	return &service{userService: userService, store: store}
}

func (s service) Get(ctx context.Context, user *user.User) ([]*TrackedInfo, error) {
	return s.store.Get(ctx, user)
}

func (s service) Create(ctx context.Context, user *user.User, userAdd *vkmodels.User) error {
	return s.store.Create(ctx, user, userAdd)
}

func (s service) GetPrevFriends(ctx context.Context, tracked *TrackedInfo) ([]int64, error) {
	return s.store.GetPrevFriends(ctx, tracked)
}

func (s service) UpdatePrevFriends(ctx context.Context, tracked *TrackedInfo, newFriends map[int64]vkmodels.User) error {
	return s.store.UpdatePrevFriends(ctx, tracked, newFriends)
}

func (s service) GetTrackedByVkID(ctx context.Context, user *user.User, vkId int64) (*TrackedInfo, error) {
	return s.store.GetTrackedByVkID(ctx, user, vkId)
}

func (s service) DeleteUserFromPrevFriends(ctx context.Context, deleteUser *vkmodels.User, tracked *TrackedInfo) error {
	return s.store.DeleteUserFromPrevFriends(ctx, deleteUser, tracked)
}

func (s service) AddUserInPrevFriends(ctx context.Context, addUser *vkmodels.User, tracked *TrackedInfo) error {
	return s.store.AddUserInPrevFriends(ctx, addUser, tracked)
}

func (s service) DeleteUserFromTracked(ctx context.Context, user *user.User, tracked *TrackedInfo) error {
	return s.store.DeleteUserFromTracked(ctx, user, tracked)
}

func (s service) GetTrackeds(ctx context.Context) (map[VkID][]*UserTrackedInfo, error) {
	return s.store.GetTrackeds(ctx)
}
