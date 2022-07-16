package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"tgbot_surveillance/config"
	trackedsvc "tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	"tgbot_surveillance/pkg/clock"
	encrypt "tgbot_surveillance/pkg/encrypt"
	govk "tgbot_surveillance/pkg/go-vk"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"
	"time"

	"github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func (s *Server) callbackQueryHandler(update *tgbotapi.CallbackQuery) error {
	usr, err := s.services.UserService.GetUserByTgID(context.Background(), user.TelegramID(update.From.ID))
	if err != nil {
		return errors.Wrapf(err, "callbackQueryHandler, tg_user_id: %v", update.From.ID)
	}

	s.logger.Infof("User: %v;	Data: %v", usr.ID, update.Data)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch {
	case strings.Contains(update.Data, addByVkIdCallback+"_"):
		err = s.addInTrackedByVkId(usr, &msg, update)

	case strings.Contains(update.Data, trackedCallback+"_"):
		err = s.trackedCallback(usr, &msg, update)

	default:
		return nil
	}

	if err != nil {
		return err
	}

	_, err = s.tg.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send msg")
	}

	return nil
}

func (s *Server) messageHandler(update *tgbotapi.Message, cfg *config.Config) error {

	if update.IsCommand() {
		return s.commandHandler(update)
	}

	usr, err := s.services.UserService.GetUserByTgID(context.Background(), user.TelegramID(update.From.ID))
	if err != nil {
		return errors.Wrapf(err, "MessageHandler, tg_user_id: %v", update.From.ID)
	}

	if usr == nil {
		if usr, err = s.services.UserService.Create(context.Background(), user.New(
			user.TelegramID(update.From.ID),
			update.From.UserName,
			clock.Real{}.Now(),
		)); err != nil {
			return errors.Wrapf(err, "MessageHandler, tg_user_id: %v", update.From.ID)
		}
	}

	s.logger.Infof("User: %v;	Text: %v", usr, update.Text)

	msg := tgbotapi.NewMessage(update.Chat.ID, "")

	if !s.isMessageButton(update.Text) && s.curTypeMessage != nil && s.curTypeMessage.TypeMessage != NotMessage {
		switch s.curTypeMessage.TypeMessage {
		case TokenVK:
			time.Sleep(1 * time.Second)
			msg1 := tgbotapi.NewDeleteMessage(update.Chat.ID, update.MessageID)
			if _, err = s.tg.Send(msg1); err != nil {
				return errors.Wrap(err, "message, not valid token")
			}
			err = s.tokenVkMessage(usr, &msg, update, cfg)
		case SearchByIdVk:
			err = s.searchByVkID(usr, &msg, update)
		case SearchByNameVk:
			err = s.searchByVkName(usr, &msg, update)
		case SearchByLinkVk:
			err = s.searchByLinkVk(usr, &msg, update)
		}
	} else {
		switch update.Text {
		case startButton:
			err = s.startButton(usr, &msg, update)

		case infoAboutBot:
			err = s.infoAboutBot(&msg, update)

		case mySubscriptionButton:
			//msg.Text = usr.GetSubscriptionMsg()
			msg.Text = "Эта часть еще не реализована("

		case contactsButton:
			msg.Text = fmt.Sprintf("Введите /start для отображения меню\n\n" +
				"Связаться с админом: @mister_kros")

		case mainButton:
			err = s.mainButton(usr, &msg)

		case getTokenVkButton:
			err = s.getTokenVkButton(&msg)

		case updateTokenVkButton:
			err = s.getTokenVkButton(&msg)

		case survaillanceButton:
			err = s.survaillanceButton(usr, &msg)

		case addToTrackedButton:
			err = s.addToTrackedButton(&msg)

		case addByVkIdButton:
			err = s.addByVkIdButton(&msg)

		case addByVkNameButton:
			err = s.addByVkNameButton(&msg)

		case addByLinkVkButton:
			err = s.addByLinkVkButton(&msg)

		case trackedButton:
			err = s.trackedButton(usr, &msg)

		case friendsByTrackedButton:
			err = s.friendsByTrackedButton(&msg)

		case getNewInfoAboutFriendsButton:
			err = s.getNewInfoAboutFriendsButton(usr, &msg)

		case deletedFromSurvaillanceButton:
			err = s.deletedFromSurvaillanceButton(usr, &msg)

		case getHistoryAboutFriends:
			err = s.getHistoryAboutFriends(usr, &msg)

		default:
			msg1 := tgbotapi.NewDeleteMessage(update.Chat.ID, update.MessageID)
			if _, err = s.tg.Send(msg1); err != nil {
				return errors.Wrap(err, "message, not delete")
			}
			return nil
		}
	}
	if err != nil {
		return err
	}

	_, err = s.tg.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send msg")
	}

	return nil
}

func (s *Server) commandHandler(update *tgbotapi.Message) error {
	switch update.Command() {
	case "start":
		var err error

		msg := tgbotapi.NewMessage(update.Chat.ID, "")
		msg.ReplyMarkup = startKeyboard
		msg.Text = "Выберите команду"

		m := tgbotapi.NewMessage(401948312, fmt.Sprintf("New user press /start, %v, %v", update.From.ID, update.From.String()))
		if _, err = s.tg.Send(m); err != nil {
			return errors.Wrap(err, "send msg")
		}
		if _, err = s.tg.Send(msg); err != nil {
			return errors.Wrap(err, "send msg")
		}

	case "invoice":
		invoice := tgbotapi.NewInvoice(update.Chat.ID, "Test invoice", "description here", "custom_payload",
			"381764678:TEST:14079", "start_param", "RUB",
			&[]tgbotapi.LabeledPrice{{Label: "RUB", Amount: 100000}})
		_, err := s.tg.Send(invoice)
		if err != nil {
			return errors.Wrap(err, "new invoice")
		}
	default:
		return nil
	}

	return nil
}

func (s *Server) startButton(user *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.Message) error {
	isToken := false
	if user.Token != nil && *user.Token != "" {
		isToken = true
	}
	text := ""

	if isToken {
		text = fmt.Sprintf(`Привет <b>%s</b>, меня зовут mr. Kros, хочу помочь тебе в слежке`, update.From.FirstName)
	} else {
		text = fmt.Sprintf("Привет <b>%s</b>, меня зовут mr. Kros, хочу помочь тебе в слежке.\n%s", update.From.FirstName, MAIN_NO_TOKEN_TEXT)
	}
	msg.Text = text
	keyboard, err := s.getMainKeyboard(isToken)
	if err != nil {
		return errors.Wrap(err, "message handler, start button")
	}
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "html"

	return nil
}

func (s *Server) infoAboutBot(msg *tgbotapi.MessageConfig, update *tgbotapi.Message) error {

	text := fmt.Sprintf(`Привет <b>%s</b>`, update.From.FirstName)
	text += "\n Бот будет отправлять сообщения, если у тех, кого вы отслеживаете, появились изменения в друзьях (появился новый друг или друг удален).\n Для этого вам необходимо добавить людей в отслеживаемые. Успехов)"

	msg.Text = text
	msg.ParseMode = "html"
	msg.ReplyMarkup = startKeyboard

	return nil
}

func (s *Server) mainButton(user *user.User, msg *tgbotapi.MessageConfig) error {
	s.curTypeMessage.ChangeTypeMessage(NotMessage)
	isToken := false
	if user.Token != nil && *user.Token != "" {
		isToken = true
	}

	msg.Text = "Выберите команду"
	keyboard, err := s.getMainKeyboard(isToken)
	if err != nil {
		return errors.Wrap(err, "message handler, posts button")
	}
	msg.ReplyMarkup = keyboard

	return nil
}

func (s *Server) getTokenVkButton(msg *tgbotapi.MessageConfig) error {
	msg.Text = GET_TOKEN_VK_TEXT
	msg.ReplyMarkup = getTokenVkKeyboard
	s.curTypeMessage.ChangeTypeMessage(TokenVK)

	return nil
}

func (s *Server) survaillanceButton(user *user.User, msg *tgbotapi.MessageConfig) error {
	msg.Text = "Выберите команду либо отслеживаемого человека"
	trackeds, err := s.services.TrackedService.Get(context.Background(), user)
	if err != nil {
		return errors.Wrap(err, "callback, api vk")
	}

	keyboardInline, keyboard, err := s.getSurvaillanceKeyboard(trackeds)
	if err != nil {
		return errors.Wrap(err, "callback, main button")
	}
	msg.ReplyMarkup = keyboard
	_, err = s.tg.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send msg")
	}
	msg.Text = "Отслеживаемые люди:"
	msg.ReplyMarkup = keyboardInline

	return nil
}

func (s *Server) addToTrackedButton(msg *tgbotapi.MessageConfig) error {
	s.curTypeMessage.ChangeTypeMessage(NotMessage)
	msg.Text = "Выберите способ добавления человека в отслеживаемые:"
	msg.ReplyMarkup = addToTrackedKeyboard

	return nil
}

func (s *Server) addByVkIdButton(msg *tgbotapi.MessageConfig) error {
	msg.Text = ADD_BY_ID_TEXT
	msg.ReplyMarkup = addByVkIdKeyboard
	s.curTypeMessage.ChangeTypeMessage(SearchByIdVk)

	return nil
}

func (s *Server) searchByVkID(user *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.Message) error {
	listUsers := make(map[int]vkmodels.User)
	idInt, err := strconv.Atoi(update.Text)
	if err != nil {
		msg.Text = ID_NOT_INT_ERROR
	} else {
		params := govk.UserGetParams{
			UserIDS: int64(idInt),
			Fields:  "id, first_name, last_name",
		}
		apiVk, _ := govk.NewApiClient(*user.Token)
		res, err := apiVk.UserGet(params)
		if err != nil {
			return errors.Wrap(err, "api vk")
		}
		if res == nil {
			msg.Text = FRIENDS_BY_NAME_EMPTY
		} else {
			listUsers[int(res.UID)] = *res
			msg.Text = "Пользователь найден! Теперь вы можете его добавить в отслеживаемые."

			keyboard, err := s.getListUsersByVkIDKeyboard(listUsers)
			if err != nil {
				return errors.Wrap(err, "message, main button")
			}
			msg.ReplyMarkup = keyboard
		}
	}

	return nil
}

func (s *Server) addByVkNameButton(msg *tgbotapi.MessageConfig) error {
	msg.Text = ADD_BY_NAME_TEXT
	msg.ReplyMarkup = addByVkIdKeyboard
	s.curTypeMessage.ChangeTypeMessage(SearchByNameVk)

	return nil
}

func (s *Server) searchByVkName(user *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.Message) error {
	listUsers := make(map[int]vkmodels.User)

	params := govk.FriendsByNameGetParams{
		UserID: int64(*user.VkID),
		Query:  update.Text,
		Fields: "id, first_name, last_name",
	}
	apiVk, _ := govk.NewApiClient(*user.Token)
	res, err := apiVk.FriendsByNameGet(params)
	if err != nil {
		return errors.Wrap(err, "api vk")
	}
	if res == nil {
		msg.Text = FRIENDS_BY_NAME_EMPTY
	} else {
		for _, friend := range res.Items {
			listUsers[int(friend.UID)] = *friend
		}
		msg.Text = "Друг найден! Теперь вы можете его добавить в отслеживаемые. \nЕсли в этом списке нет, введите полное имя фамилию пользователя."

		keyboard, err := s.getListUsersByVkNameKeyboard(listUsers)
		if err != nil {
			return errors.Wrap(err, "message, main button")
		}
		msg.ReplyMarkup = keyboard
	}

	return nil
}

func (s *Server) addByLinkVkButton(msg *tgbotapi.MessageConfig) error {
	msg.Text = ADD_BY_LINK_VK_TEXT
	msg.ReplyMarkup = addByVkIdKeyboard
	s.curTypeMessage.ChangeTypeMessage(SearchByLinkVk)

	return nil
}

func (s *Server) searchByLinkVk(user *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.Message) error {
	listUsers := make(map[int]vkmodels.User)

	params := govk.UserSearchParams{
		Query:  update.Text,
		Fields: "id, first_name, last_name",
	}
	apiVk, _ := govk.NewApiClient(*user.Token)
	res, err := apiVk.UserSearch(params)
	if err != nil {
		return errors.Wrap(err, "api vk")
	}
	if res == nil {
		msg.Text = USER_BY_LINK_VK_NOT_FOUND
	} else {
		for _, friend := range res.Items {
			listUsers[int(friend.UID)] = *friend
		}
		msg.Text = "Пользователь найден! Теперь вы можете его добавить в отслеживаемые."

		keyboard, err := s.getListUsersByVkNameKeyboard(listUsers)
		if err != nil {
			return errors.Wrap(err, "message, main button")
		}
		msg.ReplyMarkup = keyboard
	}

	return nil
}

func (s *Server) tokenVkMessage(usr *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.Message, cfg *config.Config) error {
	token, userId := CutAccessTokenAndUserId(update.Text)

	if token == "" || userId == "" {
		msg.Text = GET_UNSUCCESS_VK_TOKEN
		msg.ReplyMarkup = getTokenVkKeyboard
	} else {
		usr.Token = &token
		vkID, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			return errors.Wrap(err, "userID not int")
		}
		usr.VkID = user.VkID(&vkID)

		userCopy := usr
		encToken, err := encrypt.Encrypt(token, cfg.Secret)
		if err != nil {
			return errors.Wrap(err, "error encrypting your classified text")
		}
		userCopy.Token = &encToken

		if err = s.services.UserService.Update(context.Background(), userCopy); err != nil {
			return errors.Wrapf(err, "MessageHandler, tg_user_id: %v", update.From.ID)
		}

		msg.Text = GET_SUCCESS_VK_TOKEN
		msg.ReplyMarkup = getTokenVkKeyboard
		s.curTypeMessage.TypeMessage = NotMessage
	}

	return nil
}

func (s *Server) addInTrackedByVkId(usr *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.CallbackQuery) error {
	s.logger.Info("Callback: ", addByVkIdCallback)
	runes := []rune(update.Data)
	vkID := string(runes[len(addByVkIdCallback)+1:])
	idInt, _ := strconv.Atoi(vkID)

	apiVk, _ := govk.NewApiClient(*usr.Token)
	params := govk.UserGetParams{
		UserIDS: int64(idInt),
		Fields:  "id, first_name, last_name",
	}
	userAdd, err := apiVk.UserGet(params)
	if err != nil {
		return errors.Wrap(err, "api vk")
	}

	err = s.services.TrackedService.Create(context.Background(), usr, userAdd)
	if err != nil {
		if err == trackedsvc.ErrTrackedAlreadyExist {
			msg.Text = "Данный пользователь уже есть в ваших отслеживаемых!"
			msg.ReplyMarkup = addToTrackedKeyboard
		} else {
			return errors.Wrap(err, "callback, api vk")
		}
	} else {
		msg.Text = "Пользователь успешно добавлен в отслеживаемые."
		msg.ReplyMarkup = addToTrackedKeyboard
	}
	return nil
}

func (s *Server) trackedCallback(usr *user.User, msg *tgbotapi.MessageConfig, update *tgbotapi.CallbackQuery) error {
	s.logger.Info("Callback: ", trackedCallback)
	runes := []rune(update.Data)
	vkID := string(runes[len(trackedCallback)+1:])
	idInt, _ := strconv.Atoi(vkID)

	tracked, err := s.services.TrackedService.GetTrackedByVkID(context.Background(), usr, int64(idInt))
	if err != nil {
		return errors.Wrap(err, "message, trackedButton")
	}
	s.curTracked = tracked
	msg.Text = fmt.Sprintf("Текущий отслеживаемый: %s", tracked.UserVK.FirstName+" "+tracked.UserVK.LastName)
	keyboard, err := s.getTrackedKeyboard()
	if err != nil {
		return errors.Wrap(err, "message, main button")
	}
	msg.ReplyMarkup = keyboard

	return nil
}

func (s *Server) trackedButton(usr *user.User, msg *tgbotapi.MessageConfig) error {
	s.logger.Info("Message: ", trackedButton)

	idInt := s.curTracked.UserVK.UID

	tracked, err := s.services.TrackedService.GetTrackedByVkID(context.Background(), usr, int64(idInt))
	if err != nil {
		return errors.Wrap(err, "message, trackedButton")
	}
	s.curTracked = tracked
	msg.Text = fmt.Sprintf("Текущий отслеживаемый: %s", tracked.UserVK.FirstName+" "+tracked.UserVK.LastName)
	keyboard, err := s.getTrackedKeyboard()
	if err != nil {
		return errors.Wrap(err, "message, main button")
	}
	msg.ReplyMarkup = keyboard

	return nil
}

func (s *Server) friendsByTrackedButton(msg *tgbotapi.MessageConfig) error {
	msg.Text = "Как только изменения в друзьях появятся, бот вам об этом сообщит)\nНо вы можете посмотреть историю изменений за все время отслеживания данного пользователя."
	msg.ReplyMarkup = trackedKeyboard

	return nil
}

func (s *Server) getNewInfoAboutFriendsButton(usr *user.User, msg *tgbotapi.MessageConfig) error {
	s.logger.Info("Message: ", getNewInfoAboutFriendsButton)

	vkID := s.curTracked.UserVK.UID

	apiVk, _ := govk.NewApiClient(*usr.Token)
	params := govk.FriendsGetParams{
		UserID: int64(vkID),
		Fields: "id, first_name, last_name",
	}
	newListFriends, err := apiVk.FriendsGet(params)
	if err != nil {
		return errors.Wrap(err, "api vk")
	}
	prevListFriends, err := s.services.TrackedService.GetPrevFriends(context.Background(), s.curTracked)
	if err != nil {
		return errors.Wrap(err, "get prev")
	}

	addedFriends, deletedFriends, err := s.checkDeletedAndNewFriends(usr, newListFriends, prevListFriends, s.curTracked)
	if err != nil {
		return errors.Wrap(err, "message, check lists friends")
	}
	text := fmt.Sprintf("Текущий отслеживаемый: %s \n", s.curTracked.UserVK.FirstName+" "+s.curTracked.UserVK.LastName)
	text += s.getTextAboutAddedAndDeletedFriends(addedFriends, deletedFriends)
	msg.Text = text
	msg.ReplyMarkup = trackedKeyboard

	return nil
}

func (s *Server) deletedFromSurvaillanceButton(usr *user.User, msg *tgbotapi.MessageConfig) error {
	s.logger.Info("Message: ", deletedFromSurvaillanceButton)

	err := s.services.TrackedService.DeleteUserFromTracked(context.Background(), usr, s.curTracked)
	if err != nil {
		return errors.Wrap(err, "message, deletedFromSurvaillanceButton")
	}
	msg.Text = fmt.Sprintf("%s убран из отслеживаемых успешно", s.curTracked.UserVK.FirstName+" "+s.curTracked.UserVK.LastName)
	s.curTracked = nil
	msg.ReplyMarkup = deletedFromTrackedKeyboard

	return nil
}

func (s *Server) getHistoryAboutFriends(usr *user.User, msg *tgbotapi.MessageConfig) error {
	s.logger.Info("Message: ", getHistoryAboutFriends)

	addedFriends, deletedFriends, err := s.services.TrackedService.GetHistoryAboutFriends(context.Background(), usr, s.curTracked)
	if err != nil {
		return errors.Wrap(err, "get history")
	}

	text := s.getTextHistoryFriends(addedFriends, deletedFriends)

	msg.Text = text
	msg.ParseMode = "html"
	msg.ReplyMarkup = trackedKeyboard

	return nil
}
