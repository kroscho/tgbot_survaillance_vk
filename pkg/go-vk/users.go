package govk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tgbot_surveillance/pkg/go-vk/models"
)

const (
	apiUserGet    = "users.get"
	apiUserSearch = "users.search"
)

type User struct {
	vk *Client
}

type UserGetParams struct {
	UserIDS int64
	Fields  string
}

func (p *UserGetParams) prepareParams() map[string]string {
	params := map[string]string{
		"user_ids": fmt.Sprintf("%d", p.UserIDS),
		"fields":   fmt.Sprintf("%s", p.Fields),
	}

	return params
}

func (u *User) UserGet(userGetParams UserGetParams) (*models.User, error) {
	params := userGetParams.prepareParams()

	result := &[]models.User{}

	request, err := u.vk.prepareRequest(http.MethodGet, apiUserGet, params, nil, "", true)
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

	if len(*result) == 0 {
		return nil, nil
	} else {
		return &(*result)[0], nil
	}
}

type UserSearchParams struct {
	Query  string
	Fields string
}

func (p *UserSearchParams) prepareParams() map[string]string {
	params := map[string]string{
		"q":      fmt.Sprintf("%s", p.Query),
		"fields": fmt.Sprintf("%s", p.Fields),
	}

	return params
}

type UserSearchResult struct {
	Count int            `json:"count"`
	Items []*models.User `json:"items"`
}

func (u *User) UserSearch(userSearchParams UserSearchParams) (*UserSearchResult, error) {
	params := userSearchParams.prepareParams()

	result := &UserSearchResult{}

	request, err := u.vk.prepareRequest(http.MethodGet, apiUserSearch, params, nil, "", true)
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

	fmt.Println("REsult: ", result.Count, result.Items)

	if result.Count == 0 {
		return nil, nil
	} else {
		return result, nil
	}
}
