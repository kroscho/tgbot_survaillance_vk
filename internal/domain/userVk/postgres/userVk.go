package postgres

import (
	userVksvc "tgbot_surveillance/internal/domain/userVk"
)

type userVk struct {
	ID        userVksvc.ID   `db:"id_user_vk"`
	VkID      userVksvc.VkID `db:"vk_id"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
}

func (t userVk) marshal() (*userVksvc.UserVk, error) {
	return &userVksvc.UserVk{
		ID:        userVksvc.ID(t.ID),
		VkID:      userVksvc.VkID(t.VkID),
		FirstName: string(t.FirstName),
		LastName:  string(t.LastName),
	}, nil
}

func (t *userVk) unmarshal(from *userVksvc.UserVk) {
	*t = userVk{
		ID:        userVksvc.ID(from.ID),
		VkID:      userVksvc.VkID(from.VkID),
		FirstName: string(from.FirstName),
		LastName:  string(from.LastName),
	}
}
