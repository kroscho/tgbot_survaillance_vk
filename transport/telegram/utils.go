package telegram

import (
	"context"
	"strings"
	"tgbot_surveillance/internal/domain/tracked"
	"tgbot_surveillance/internal/domain/user"
	govk "tgbot_surveillance/pkg/go-vk"
	vkmodels "tgbot_surveillance/pkg/go-vk/models"

	"github.com/pkg/errors"
)

// вырезать токен из строки, отправленной пользователем
func CutAccessTokenAndUserId(str string) (string, string) {
	runes := []rune(str)
	tokenStr := ""
	userIdStr := ""
	startToken := strings.Index(str, "access_token=")
	endToken := strings.Index(str, "&expires_in")
	startUserId := strings.Index(str, "user_id=")
	if startToken != -1 && endToken != -1 && startUserId != -1 {
		tokenStr = string(runes[startToken+len("access_token=") : endToken])
		userIdStr = string(runes[startUserId+len("user_id="):])
	}
	return tokenStr, userIdStr
}

// выявить новых друзей на основе прежнего и нового списка друзей
func (s Server) checkDeletedAndNewFriends(user *user.User, newListFriends map[int64]vkmodels.User, prevListFriends []int64, tracked *tracked.TrackedInfo) (map[int64]vkmodels.User, map[int64]vkmodels.User, error) {
	addedFriendsIds := make(map[int64]vkmodels.User)
	deletedFriendsIds := make(map[int64]vkmodels.User)
	i := 0
	for key, val := range newListFriends {
		isDeleted := false
		// проверяем удаленных друзей (если человек есть в прежнем списке, но нет в новом)
		if i < len(prevListFriends) {
			_, ok := newListFriends[prevListFriends[i]]
			if !ok {
				params := govk.UserGetParams{
					UserIDS: int64(prevListFriends[i]),
					Fields:  "id, first_name, last_name",
				}
				apiVk, _ := govk.NewApiClient(*user.Token)
				friend, err := apiVk.UserGet(params)
				if err != nil {
					return addedFriendsIds, deletedFriendsIds, errors.Wrap(err, "api vk")
				}
				if friend != nil {
					deletedFriendsIds[int64(friend.UID)] = *friend
					err = s.services.TrackedService.DeleteUserFromPrevFriends(context.Background(), friend, tracked)
					if err != nil {
						return addedFriendsIds, deletedFriendsIds, errors.Wrap(err, "delete friend")
					}
					isDeleted = true
				}
			}
		}
		if !isDeleted && !CheckExistInList(key, prevListFriends) {
			addedFriendsIds[key] = val
			err := s.services.TrackedService.AddUserInPrevFriends(context.Background(), &val, tracked)
			if err != nil {
				return addedFriendsIds, deletedFriendsIds, errors.Wrap(err, "added friend")
			}
		}
		i += 1
	}
	// если в прежнем списке еще остались ids, проверяем, есть ли они в новом (проверка удаленных друзей)
	if i < len(prevListFriends) {
		prevListPart := prevListFriends[i:]
		for _, idVK := range prevListPart {
			_, ok := newListFriends[idVK]
			if !ok {
				params := govk.UserGetParams{
					UserIDS: int64(prevListFriends[i]),
					Fields:  "id, first_name, last_name",
				}
				apiVk, _ := govk.NewApiClient(*user.Token)
				friend, err := apiVk.UserGet(params)
				if err != nil {
					return addedFriendsIds, deletedFriendsIds, errors.Wrap(err, "api vk")
				}
				if friend != nil {
					deletedFriendsIds[int64(friend.UID)] = *friend
					err = s.services.TrackedService.DeleteUserFromPrevFriends(context.Background(), friend, s.curTracked)
					if err != nil {
						return addedFriendsIds, deletedFriendsIds, errors.Wrap(err, "delete friend")
					}
				}
			}
		}
	}
	return addedFriendsIds, deletedFriendsIds, nil
}

func CheckExistInList(x int64, list []int64) bool {
	for _, v := range list {
		if v == x {
			return true
		}
	}
	return false
}

// получить текст о добавленных и удаленных друзьях
func (s Server) getTextAboutAddedAndDeletedFriends(addedFriendsIds map[int64]vkmodels.User, deletedFriendsIds map[int64]vkmodels.User) string {
	text := ""
	if len(addedFriendsIds) == 0 && len(deletedFriendsIds) == 0 {
		text += "Пока изменений нет"
	} else {
		if len(addedFriendsIds) != 0 {
			text += "Новые друзья:\n"
		}
		for _, addedFriend := range addedFriendsIds {
			text += "Новый друг: " + addedFriend.FirstName + " " + addedFriend.LastName + "\n"
		}
		if len(deletedFriendsIds) != 0 {
			text += "Удаленные друзья:\n"
		}
		for _, deletedFriend := range deletedFriendsIds {
			text += "Удаленный друг: " + deletedFriend.FirstName + " " + deletedFriend.LastName + "\n"
		}
	}
	return text
}

func (s Server) isMessageButton(updateText string) bool {
	switch updateText {
	case startButton:
		return true
	case mainButton:
		return true
	case mySubscriptionButton:
		return true
	case contactsButton:
		return true
	case getTokenVkButton:
		return true
	case updateTokenVkButton:
		return true
	case survaillanceButton:
		return true
	case addToTrackedButton:
		return true
	case addByVkIdButton:
		return true
	case addByLinkVkButton:
		return true
	case addByVkNameButton:
		return true
	case trackedsButton:
		return true
	case newFriendsButton:
		return true
	case addedAndDeletedFriendsButton:
		return true
	case friendsByTrackedButton:
		return true
	case likesByTrackedButtom:
		return true
	case getNewInfoAboutFriendsButton:
		return true
	case getHistoryAboutFriends:
		return true
	case deletedFromSurvaillanceButton:
		return true
	default:
		return false
	}
}
