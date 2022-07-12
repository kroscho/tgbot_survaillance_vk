package govk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tgbot_surveillance/pkg/go-vk/models"
)

const (
	apiFriendsBuNameGet = "friends.search"
	apiFriendsGet       = "friends.get"
)

type Friend struct {
	vk *Client
}

type FriendsByNameGetParams struct {
	UserID int64
	Query  string
	Fields string
}

func (p *FriendsByNameGetParams) prepareParams() map[string]string {
	params := map[string]string{
		"user_id": fmt.Sprintf("%d", p.UserID),
		"q":       fmt.Sprintf("%s", p.Query),
		"fields":  fmt.Sprintf("%s", p.Fields),
	}

	return params
}

type FriendsByNameGetResult struct {
	Count int            `json:"count"`
	Items []*models.User `json:"items"`
}

func (u *User) FriendsByNameGet(friendsByNameGetParams FriendsByNameGetParams) (*FriendsByNameGetResult, error) {
	params := friendsByNameGetParams.prepareParams()

	result := &FriendsByNameGetResult{}

	request, err := u.vk.prepareRequest(http.MethodGet, apiFriendsBuNameGet, params, nil, "", true)
	if err != nil {
		return nil, err
	}

	response, err := u.vk.do(request)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(response, result); err != nil {
		return nil, err
	}

	if result.Count == 0 {
		return nil, nil
	} else {
		return result, nil
	}
}

type FriendsGetParams struct {
	UserID int64
	Fields string
}

func (p *FriendsGetParams) prepareParams() map[string]string {
	params := map[string]string{
		"user_id": fmt.Sprintf("%d", p.UserID),
		"fields":  fmt.Sprintf("%s", p.Fields),
	}

	return params
}

type FriendsGetResult struct {
	Count int            `json:"count"`
	Items []*models.User `json:"items"`
}

func (u *User) FriendsGet(friendsGetParams FriendsGetParams) (map[int64]models.User, error) {
	params := friendsGetParams.prepareParams()

	result := &FriendsGetResult{}

	request, err := u.vk.prepareRequest(http.MethodGet, apiFriendsGet, params, nil, "", true)
	if err != nil {
		return nil, err
	}

	response, err := u.vk.do(request)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(response, result); err != nil {
		return nil, err
	}

	if result.Count == 0 {
		return nil, nil
	} else {
		listFriendsMap := make(map[int64]models.User)
		for _, friend := range result.Items {
			listFriendsMap[friend.UID] = *friend
		}
		return listFriendsMap, nil
	}
}
