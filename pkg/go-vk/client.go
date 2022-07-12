package govk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"tgbot_surveillance/pkg/go-vk/models"
	"time"

	"github.com/pkg/errors"
)

const (
	version  = "5.131"
	apiURL   = "https://api.vk.com/method/"
	tokenURL = "https://oauth.vk.com/access_token"
)

type VkApiClient interface {
	Request(method string, params map[string]string) ([]byte, error)
	MethodsUsers
}

type MethodsUsers interface {
	UserGet(usersGetParams UserGetParams) (*models.User, error)
	FriendsByNameGet(friendsByNameGetParams FriendsByNameGetParams) (*FriendsByNameGetResult, error)
	UserSearch(userSearchParams UserSearchParams) (*UserSearchResult, error)
	FriendsGet(friendsGetParams FriendsGetParams) (map[int64]models.User, error)
}

type Client struct {
	AccessToken  string
	ClientSecret string
	Version      string
	httpClient   *http.Client
	httpTimeout  time.Duration
	User
}

func NewApiClient(token string) (VkApiClient, error) {
	vk := &Client{
		AccessToken: token,
		Version:     version,
		httpClient:  &http.Client{},
		httpTimeout: 5 * time.Second,
	}

	vk.User = User{vk: vk}

	return vk, nil
}

func (c *Client) Request(method string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(apiURL + method)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}

	query.Set("access_token", c.AccessToken)
	query.Set("v", c.Version)
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var handler struct {
		Error    *models.Error
		Response json.RawMessage
	}
	err = json.Unmarshal(body, &handler)

	if handler.Error != nil {
		return nil, handler.Error
	}

	return handler.Response, nil
}

func (c *Client) prepareRequest(method, path string, params, headers map[string]string, body string, auth bool) (*http.Request, error) {
	u := apiURL + path

	req, err := http.NewRequest(method, u, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}

	if auth {
		query.Set("access_token", c.AccessToken)
	}

	query.Set("v", c.Version)
	req.URL.RawQuery = query.Encode()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var handler struct {
		Error    *models.Error
		Response json.RawMessage
	}

	err = json.Unmarshal(response, &handler)
	if err != nil {
		return nil, err
	}

	if handler.Error != nil {
		return nil, handler.Error
	}

	return handler.Response, nil
}

func (c *Client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()

	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from VK API")
	}
}
