package telegram

import (
	"context"
	"fmt"
	"sync"
	"tgbot_surveillance/config"
	"tgbot_surveillance/internal/domain/tracked"
	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	"tgbot_surveillance/internal/domain/userVk"
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
	UserVkService  userVk.Service
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

	go func() {
		for {
			fmt.Println("GOGOGO")
			err := s.runNotifications()
			if err != nil {
				s.logger.Errorf("%+v", err)
				return
			}
			time.Sleep(time.Duration(cfg.NotificationDuration) * time.Minute)
		}
	}()

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
				s.wg.Add(1)
				go s.proccessUpdate(update, cfg)
			}
		}
	}()

	<-doneC

	return nil
}

func (s *Server) proccessUpdate(update tgbotapi.Update, cfg *config.Config) {
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

func (s *Server) Stop() {
	s.tg.StopReceivingUpdates()
}
