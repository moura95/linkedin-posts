package jsonplaceholder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	httpClient = &http.Client{}
)

var (
	baseURL  = "https://jsonplaceholder.typicode.com"
	apiUsers = "/users"
)

type User struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func makeHttpRequest(req *http.Request, responseObject interface{}) (*http.Response, error) {
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal object
	if responseObject != nil {
		err = json.Unmarshal(body, responseObject)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal json, body: %s, unmarshal err: %s", string(body), err)
		}
	}

	return res, nil
}

func newHttpRequest(method string, path string, data io.Reader) (*http.Request, error) {
	urlRequest := baseURL + path
	req, err := http.NewRequest(method, urlRequest, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return req, nil
}

func ListUsers() (users []User, err error) {
	req, err := newHttpRequest(http.MethodGet, apiUsers, nil)
	if err != nil {
		return nil, err
	}

	_, err = makeHttpRequest(req, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id int32) (user *User, err error) {
	urlApi := fmt.Sprintf("%s/%d", apiUsers, id)
	req, err := newHttpRequest(http.MethodGet, urlApi, nil)
	if err != nil {
		return nil, err
	}

	res, err := makeHttpRequest(req, &user)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, nil
	}

	return user, nil
}

func CreateUser(name, username, email string) (user *User, err error) {
	userReq := User{
		Name:     name,
		Username: username,
		Email:    email,
	}
	data, err := json.Marshal(userReq)

	if err != nil {
		return nil, err

	}

	req, err := newHttpRequest(http.MethodPost, apiUsers, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	_, err = makeHttpRequest(req, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
