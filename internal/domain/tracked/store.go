package tracked

import (
	"context"
	"tgbot_surveillance/internal/domain/user"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"
)

type Store interface {
	Get(ctx context.Context, user *user.User) ([]*TrackedInfo, error)
	Create(ctx context.Context, user *user.User, userAdd *vkmodels.User) error
	GetPrevFriends(ctx context.Context, tracked *TrackedInfo) ([]int64, error)
	UpdatePrevFriends(ctx context.Context, tracked *TrackedInfo, newFriends map[int64]vkmodels.User) error
	GetTrackedByVkID(ctx context.Context, user *user.User, vkId int64) (*TrackedInfo, error)
	DeleteUserFromPrevFriends(ctx context.Context, deleteUser *vkmodels.User, tracked *TrackedInfo) error
	AddUserInPrevFriends(ctx context.Context, addUser *vkmodels.User, tracked *TrackedInfo) error
	DeleteUserFromTracked(ctx context.Context, user *user.User, tracked *TrackedInfo) error
	GetTrackeds(ctx context.Context) (map[VkID][]*UserTrackedInfo, error)
	AddInHistory(ctx context.Context, from *TrackedInfo, addedFriends map[int64]vkmodels.User, deletedFriends map[int64]vkmodels.User) error
	GetHistoryAboutFriends(ctx context.Context, user *user.User, tracked *TrackedInfo) (map[string][]*HistoryVk, map[string][]*HistoryVk, error)
}
