package user

import (
	"time"

	"github.com/pkg/errors"
)

type ID int64

type TelegramID int64

type TypeMessage struct {
	TypeMessage int
}

type VkID *int64

/*
type User struct {
	ID          ID
	TgID        TelegramID
	VkID        VkID
	Username    string
	Token       *string
	Enabled     *bool
	LastPayment *time.Time
	Plan        *string
	Price       *float64
	CreatedAt   time.Time
}
*/

type User struct {
	ID        ID
	TgID      TelegramID
	VkID      VkID
	Username  string
	Token     *string
	CreatedAt time.Time
}

func New(tgID TelegramID, username string, createdAt time.Time) *User {
	return &User{TgID: tgID, Username: username, CreatedAt: createdAt}
}

/*
func (u *User) GetSubscriptionMsg() string {
	if u.Enabled == nil || !*u.Enabled {
		return "У Вас нет активных подписок"
	}

	duration := 30 * 24 * time.Hour

	return fmt.Sprintf("**Информация по подписке:** \n\n"+
		"	*План:* _%s_ \n"+
		"	*Дата продления:* _%v_ \n"+
		"	*Дата окончания:* _%v_ \n"+
		"	*Стоимость:* _%v руб._",
		*u.Plan, u.LastPayment.Local().Format("2006-01-02"), u.LastPayment.Local().Add(duration).Format("2006-01-02"), *u.Price)
}
*/

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrInvalidUser      = errors.New("invalid user")
)
