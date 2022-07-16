package telegram

import (
	"context"
	"fmt"
	"tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	govk "tgbot_surveillance/pkg/go-vk"

	"github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func (s Server) runNotifications() error {
	s.logger.Info("Start notifications...")
	trackeds, err := s.services.TrackedService.GetTrackeds(context.Background())
	if err != nil {
		return errors.Wrap(err, "callback, api vk")
	}

	for key, users := range trackeds {
		usr, err := s.services.UserService.GetUserByTgID(context.Background(), user.TelegramID(users[0].TgID))

		if err != nil {
			return errors.Wrapf(err, "tg_user_id: %v", users[0].TgID)
		}

		apiVk, _ := govk.NewApiClient(*usr.Token)
		params := govk.FriendsGetParams{
			UserID: int64(key),
			Fields: "id, first_name, last_name",
		}
		newListFriends, err := apiVk.FriendsGet(params)
		if err != nil {
			return errors.Wrap(err, "api vk")
		}
		prevListFriends, err := s.services.TrackedService.GetPrevFriends(
			context.Background(),
			&tracked.TrackedInfo{
				ID: users[0].ID,
			},
		)
		if err != nil {
			return errors.Wrap(err, "get prev")
		}

		addedFriends, deletedFriends, err := s.checkDeletedAndNewFriends(usr, newListFriends, prevListFriends, &tracked.TrackedInfo{ID: users[0].ID})
		if err != nil {
			return errors.Wrap(err, "message, check lists friends")
		}
		text := s.getTextAboutAddedAndDeletedFriends(addedFriends, deletedFriends)

		if text != "Пока изменений нет" {
			for _, u := range users {
				text1 := fmt.Sprintf("Отслеживаемый: %s \n%s", users[0].FirstName+" "+users[0].LastName, text)
				m := tgbotapi.NewMessage(int64(u.TgID), text1)
				if _, err = s.tg.Send(m); err != nil {
					return errors.Wrap(err, "send msg")
				}
			}
			err = s.services.TrackedService.AddInHistory(context.Background(), &tracked.TrackedInfo{ID: users[0].ID}, addedFriends, deletedFriends)
			if err != nil {
				fmt.Println("Error: ", err)
				return errors.Wrap(err, "add in history")
			}
		}
	}

	return nil
}
