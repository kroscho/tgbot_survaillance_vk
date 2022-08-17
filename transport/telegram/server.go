package telegram

import (
	"context"
	"fmt"
	"sync"
	"tgbot_surveillance/config"
	"tgbot_surveillance/internal/domain/tracked"
	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	govk "tgbot_surveillance/pkg/go-vk"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	NotMessage     = "Not message"
	TokenVK        = "TokenVK"
	SearchByNameVk = "SearchByNameVk"
	SearchByIdVk   = "SearchByIdVK"
	SearchByLinkVk = "SearchBuLinkVk"
)

type TypeMessage struct {
	TypeMessage string
}

func NewTypeMessage() *TypeMessage {
	typeMes := TypeMessage{}
	typeMes.TypeMessage = NotMessage

	return &typeMes
}

func (t *TypeMessage) ChangeTypeMessage(typeMes string) {
	t.TypeMessage = typeMes
}

type Services struct {
	UserService    user.Service
	TrackedService tracked.Service
}

type Server struct {
	tg             *tgbotapi.BotAPI
	logger         logrus.FieldLogger
	services       Services
	usersUpdates   *usersUpdates
	curTypeMessage *TypeMessage
	curTracked     *trackedsvc.TrackedInfo
	wg             *sync.WaitGroup
}

func NewServer(tg *tgbotapi.BotAPI, logger logrus.FieldLogger, services Services) *Server {
	return &Server{
		tg:             tg,
		logger:         logger,
		services:       services,
		usersUpdates:   newUsersUpdates(),
		curTypeMessage: NewTypeMessage(),
		curTracked:     nil,
		wg:             &sync.WaitGroup{},
	}
}

func (s *Server) Run(ctx context.Context, cfg *config.Config) error {
	s.logger.Info("Start application...")
	defer s.logger.Warnln("Stop application")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := s.tg.GetUpdatesChan(u)
	if err != nil {
		return errors.Wrap(err, "get updates chan")
	}

	go s.runNotifications()

	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		for {
			select {
			case <-ctx.Done():
				s.logger.Warnln(ctx.Err())
				s.wg.Wait()
				return
			case update := <-updates:
				go s.proccessUpdate(update, cfg)
			case <-time.After(10 * time.Minute):
				go s.runNotifications()
			}
		}
	}()

	<-doneC

	return nil
}

func (s *Server) proccessUpdate(update tgbotapi.Update, cfg *config.Config) {
	s.wg.Add(1)
	defer s.wg.Done()

	var err error

	switch {
	case update.CallbackQuery != nil:
		fields := logrus.Fields{
			"tg_user_id": update.CallbackQuery.From.ID,
			"user_data":  update.CallbackQuery.From,
			"data":       update.CallbackQuery.Data,
		}
		err = s.callbackQueryHandler(update.CallbackQuery)
		if err != nil {
			s.logger.WithFields(fields).Errorf("%+v", err)
			return
		}
		s.logger.WithFields(fields).Infof("New callback query %v", update.CallbackQuery.From.String())
	case update.Message != nil:
		fields := logrus.Fields{
			"tg_user_id": update.Message.From.ID,
			"user_data":  update.Message.From,
			"text":       update.Message.Text,
		}
		err = s.messageHandler(update.Message, cfg)
		if err != nil {
			s.logger.WithFields(fields).Errorf("%+v", err)
			return
		}
		s.logger.WithFields(fields).Infof("New message %v", update.Message.From.String())
	}
}

func (s Server) runNotifications() {
	s.wg.Add(1)
	defer s.wg.Done()

	s.logger.Info("Start notifications...")
	trackeds, err := s.services.TrackedService.GetTrackeds(context.Background())
	if err != nil {
		s.logger.Errorf("%+v", err)
		return
	}

	for key, users := range trackeds {
		usr, err := s.services.UserService.GetUserByTgID(context.Background(), user.TelegramID(users[0].TgID))

		if err != nil {
			s.logger.Errorf("%+v", err)
			return
		}

		apiVk, _ := govk.NewApiClient(*usr.Token)
		params := govk.FriendsGetParams{
			UserID: int64(key),
			Fields: "id, first_name, last_name",
		}
		newListFriends, err := apiVk.FriendsGet(params)
		if err != nil {
			/*
				text := fmt.Sprintf("%s, %s", usr.Username, "У вас проблема с токеном, обновите его\nЛибо покиньте бота, чтоб не получать уведомления.")
				m := tgbotapi.NewMessage(int64(usr.TgID), text)
				m.ReplyMarkup, err = s.getMainKeyboard(true)
				if err != nil {
					s.logger.Errorf("%+v", err)
					return
				}
				if _, err = s.tg.Send(m); err != nil {
					s.logger.Errorf("%+v", err)
					return
				}
				s.logger.Errorf("%+v", err)
			*/
			text := fmt.Sprintf("%s, %s", usr.Username, "У вас проблема с токеном, обновите его\nЛибо покиньте бота, чтоб не получать уведомления.")
			s.logger.Info(text)
			return
		}
		prevListFriends, err := s.services.TrackedService.GetPrevFriends(
			context.Background(),
			&tracked.TrackedInfo{
				ID: users[0].ID,
			},
		)
		if err != nil {
			s.logger.Errorf("%+v", err)
			return
		}

		addedFriends, deletedFriends, err := s.checkDeletedAndNewFriends(usr, newListFriends, prevListFriends, &tracked.TrackedInfo{ID: users[0].ID})
		if err != nil {
			s.logger.Errorf("%+v", err)
			return
		}
		text := s.getTextAboutAddedAndDeletedFriends(addedFriends, deletedFriends, *users[0])

		if text != "Пока изменений нет" {
			for _, u := range users {
				text1 := fmt.Sprintf("Отслеживаемый: %s \n%s", users[0].FirstName+" "+users[0].LastName, text)
				m := tgbotapi.NewMessage(int64(u.TgID), text1)
				if _, err = s.tg.Send(m); err != nil {
					s.logger.Errorf("%+v", err)
					return
				}
			}
			err = s.services.TrackedService.AddInHistory(context.Background(), &tracked.TrackedInfo{ID: users[0].ID}, addedFriends, deletedFriends)
			if err != nil {
				s.logger.Errorf("%+v", err)
				return
			}
		}
	}
	s.logger.Info("Start notifications...")
}

func (s *Server) Stop() {
	s.tg.StopReceivingUpdates()
}
