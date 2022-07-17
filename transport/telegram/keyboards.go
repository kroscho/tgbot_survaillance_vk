package telegram

import (
	"strconv"
	"tgbot_surveillance/internal/domain/tracked"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	startCommand                  = "/start"
	startButton                   = "Начать"
	infoAboutBot                  = "Что делает бот?"
	surveillanceButton            = "Слежка"
	logoutOfTheBotButton          = "Покинуть бота("
	yesLogoutOfTheBotButton       = "Да, уверен"
	noLogoutOfTheBotButton        = "Нет, случайно нажал"
	mySubscriptionButton          = "Подписка"
	contactsButton                = "Помощь"
	mainButton                    = "Главная"
	getTokenVkButton              = "Получить токен ВК"
	updateTokenVkButton           = "Обновить токен ВК"
	survaillanceButton            = "Слежка" // слежка
	addByVkIdButton               = "Добавить по VK ID"
	addByLinkVkButton             = "Добавить по ссылке профиля"
	addByVkNameButton             = "Добавить по имени VK"
	trackedsButton                = "Отслеживаемые"            // отслеживаемые люди
	addToTrackedButton            = "Добавить нового"          // добавить в отслеживаемые
	addInTrackedButton            = "Добавить в отслеживаемые" // добавить в отслеживаемые конкретного
	trackedButton                 = "Отслеживаемый"            // отслеживаемый
	newFriendsButton              = "new friends"
	addedAndDeletedFriendsButton  = "added and deleted friends"
	friendsByTrackedButton        = "Друзья"
	likesByTrackedButtom          = "Лайки"
	getNewInfoAboutFriendsButton  = "Проверить друзей"
	getHistoryAboutFriends        = "История изменений"
	DELETED_FRIENDS               = "deleted friends"
	deletedFromSurvaillanceButton = "Убрать из отслеживаемых"
	HELP_TRACKED                  = "help tracked"
	likesButton                   = "likes"
	CHANGE_TYPE_SURVEILLANCE      = "change type surveillance"
	ADD_BY_NAME1                  = "add by name"
	ADD_BY_ID1                    = "add by id"

	addByVkIdCallback = "AddByVkID"
	trackedCallback   = "Tracked"
)

var (
	startKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(startButton),
			tgbotapi.NewKeyboardButton(mySubscriptionButton),
			tgbotapi.NewKeyboardButton(contactsButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(infoAboutBot),
		),
	)

	getTokenVkKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	startCommandKeybpard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(startCommand),
		),
	)

	addToTrackedKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(addByVkNameButton),
			tgbotapi.NewKeyboardButton(addByVkIdButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(addByLinkVkButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(survaillanceButton),
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	addByVkIdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(addToTrackedButton),
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	trackedKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			//tgbotapi.NewKeyboardButton(getNewInfoAboutFriendsButton),
			tgbotapi.NewKeyboardButton(getHistoryAboutFriends),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(trackedButton),
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	deletedFromTrackedKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(survaillanceButton),
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	logoutOfTheBotKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(yesLogoutOfTheBotButton),
			tgbotapi.NewKeyboardButton(noLogoutOfTheBotButton),
		),
	)
)

// клавиатура для "Главная"
func (s *Server) getMainKeyboard(isToken bool) (*tgbotapi.ReplyKeyboardMarkup, error) {
	var keyboard tgbotapi.ReplyKeyboardMarkup
	if isToken {
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(updateTokenVkButton),
				tgbotapi.NewKeyboardButton(survaillanceButton),
				tgbotapi.NewKeyboardButton(contactsButton),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(infoAboutBot),
				tgbotapi.NewKeyboardButton(logoutOfTheBotButton),
			),
		)
	} else {
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(getTokenVkButton),
				tgbotapi.NewKeyboardButton(contactsButton),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(infoAboutBot),
			),
		)
	}
	return &keyboard, nil
}

// клавиатура для "Слежка"
func (s *Server) getSurvaillanceKeyboard(trackeds []*tracked.TrackedInfo) (*tgbotapi.InlineKeyboardMarkup, *tgbotapi.ReplyKeyboardMarkup, error) {
	var keyboardInline tgbotapi.InlineKeyboardMarkup
	var keyboard tgbotapi.ReplyKeyboardMarkup

	listButtons1 := []tgbotapi.InlineKeyboardButton{}
	listButtons2 := []tgbotapi.InlineKeyboardButton{}
	listButtons3 := []tgbotapi.InlineKeyboardButton{}
	listButtons4 := []tgbotapi.InlineKeyboardButton{}
	listButtons5 := []tgbotapi.InlineKeyboardButton{}
	i := 1
	for _, user := range trackeds {
		nameFriend := user.UserVK.FirstName + " " + user.UserVK.LastName
		data := trackedCallback + "_" + strconv.Itoa(int(user.UserVK.UID))
		if i <= 2 {
			listButtons1 = append(listButtons1, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 2 && i <= 4 {
			listButtons2 = append(listButtons2, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 4 && i <= 6 {
			listButtons3 = append(listButtons3, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 6 && i <= 8 {
			listButtons4 = append(listButtons4, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		} else if i > 8 && i <= 10 {
			listButtons5 = append(listButtons5, tgbotapi.NewInlineKeyboardButtonData(nameFriend, data))
		}
		i += 1
	}

	if len(trackeds) != 0 {
		keyboardInline = tgbotapi.NewInlineKeyboardMarkup(
			listButtons1[:],
			listButtons2[:],
			listButtons3[:],
			listButtons4[:],
			listButtons5[:],
		)
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(addToTrackedButton),
				tgbotapi.NewKeyboardButton(mainButton),
			),
		)
	} else {
		keyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(addToTrackedButton),
				tgbotapi.NewKeyboardButton(mainButton),
			),
		)
	}
	return &keyboardInline, &keyboard, nil
}

// клавиатура для "Добавить в отслеживаемые"
func (s *Server) getListUsersByVkIDKeyboard(listUsers map[int]vkmodels.User) (*tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup
	listButtons := []tgbotapi.InlineKeyboardButton{}

	for key, user := range listUsers {
		nameUser := user.FirstName + " " + user.LastName
		data := addByVkIdCallback + "_" + strconv.Itoa(key)
		listButtons = append(listButtons, tgbotapi.NewInlineKeyboardButtonData(nameUser, nameUser))
		listButtons = append(listButtons, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
	}

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(listButtons...),
	)
	return &keyboard, nil
}

// клавиатура для "Добавить в отслеживаемые"
func (s *Server) getListUsersByVkNameKeyboard(listUsers map[int]vkmodels.User) (*tgbotapi.InlineKeyboardMarkup, error) {
	var keyboard tgbotapi.InlineKeyboardMarkup
	listButtons1 := []tgbotapi.InlineKeyboardButton{}
	listButtons2 := []tgbotapi.InlineKeyboardButton{}
	listButtons3 := []tgbotapi.InlineKeyboardButton{}
	listButtons4 := []tgbotapi.InlineKeyboardButton{}
	listButtons5 := []tgbotapi.InlineKeyboardButton{}
	listButtons6 := []tgbotapi.InlineKeyboardButton{}
	listButtons7 := []tgbotapi.InlineKeyboardButton{}

	i := 1
	for key, user := range listUsers {
		nameFriend := user.FirstName + " " + user.LastName
		data := addByVkIdCallback + "_" + strconv.Itoa(key)
		if i == 1 {
			listButtons1 = append(listButtons1, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons1 = append(listButtons1, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 2 {
			listButtons2 = append(listButtons2, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons2 = append(listButtons2, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 3 {
			listButtons3 = append(listButtons3, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons3 = append(listButtons3, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 4 {
			listButtons4 = append(listButtons4, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons4 = append(listButtons4, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 5 {
			listButtons5 = append(listButtons5, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons5 = append(listButtons5, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 6 {
			listButtons6 = append(listButtons6, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons6 = append(listButtons6, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		} else if i == 7 {
			listButtons7 = append(listButtons7, tgbotapi.NewInlineKeyboardButtonData(nameFriend, nameFriend))
			listButtons7 = append(listButtons7, tgbotapi.NewInlineKeyboardButtonData("Добавить", data))
		}
		i += 1
	}

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		listButtons1[:],
		listButtons2[:],
		listButtons3[:],
		listButtons4[:],
		listButtons5[:],
		listButtons6[:],
		listButtons7[:],
	)
	return &keyboard, nil
}

// клавиатура для "Добавить в отслеживаемые"
func (s *Server) getTrackedKeyboard() (*tgbotapi.ReplyKeyboardMarkup, error) {
	var keyboard tgbotapi.ReplyKeyboardMarkup

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(friendsByTrackedButton),
			tgbotapi.NewKeyboardButton(likesByTrackedButtom),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(deletedFromSurvaillanceButton),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(survaillanceButton),
			tgbotapi.NewKeyboardButton(mainButton),
		),
	)

	return &keyboard, nil
}
