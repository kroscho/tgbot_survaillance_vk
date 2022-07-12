package postgres

import (
	usersvc "tgbot_surveillance/internal/domain/user"
	"time"
)

/*
type user struct {
	ID          usersvc.ID         `db:"id_user"`
	TgID        usersvc.TelegramID `db:"tg_id"`
	VkID        usersvc.VkID       `db:"vk_id"`
	Username    string             `db:"user_name"`
	Token       *string            `db:"user_token"`
	Enabled     *bool              `db:"enabled"`
	LastPayment *time.Time         `db:"last_payment"`
	Plan        *string            `db:"plan"`
	Price       *float64           `db:"price"`
	CreatedAt   time.Time          `db:"created_at"`
}
*/

type user struct {
	ID        usersvc.ID         `db:"id_user"`
	TgID      usersvc.TelegramID `db:"tg_id"`
	VkID      usersvc.VkID       `db:"vk_id"`
	Username  string             `db:"user_name"`
	Token     *string            `db:"user_token"`
	CreatedAt time.Time          `db:"created_at"`
}

func (u user) marshal() (*usersvc.User, error) {
	return &usersvc.User{
		ID:        usersvc.ID(u.ID),
		TgID:      usersvc.TelegramID(u.TgID),
		VkID:      usersvc.VkID(u.VkID),
		Username:  u.Username,
		Token:     u.Token,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (u *user) unmarshal(from *usersvc.User) {
	*u = user{
		ID:        usersvc.ID(from.ID),
		TgID:      usersvc.TelegramID(from.TgID),
		VkID:      usersvc.VkID(from.VkID),
		Username:  from.Username,
		Token:     from.Token,
		CreatedAt: from.CreatedAt.UTC(),
	}
}
