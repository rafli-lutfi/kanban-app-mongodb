package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rafli-lutfi/kanban-app-mongodb/config"
)

type UserFetch interface {
	Register(fullname, email, password string) (respCode int, err error)
	Login(email, password string) (userID string, respCode int, err error)
	Logout(userID string) error
}

type userFetch struct{}

func NewUserFetch() *userFetch {
	return &userFetch{}
}

func (u *userFetch) Register(fullname, email, password string) (respCode int, err error) {
	datajson := map[string]string{
		"fullname": fullname,
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/users/register"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (u *userFetch) Login(email, password string) (userID string, respCode int, err error) {
	datajson := map[string]string{
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return "", -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/users/login"), bytes.NewBuffer(data))
	if err != nil {
		return "", -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", -1, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", -1, err
	}

	var result map[string]interface{}

	json.Unmarshal(b, &result)

	if result["data"] == nil {
		return "", -1, nil
	}

	dataRespone := result["data"].(map[string]interface{})
	id := dataRespone["id"].(string)

	return id, http.StatusOK, nil
}

func (u *userFetch) Logout(userID string) error {
	client, err := getClientWithCookie(userID)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/users/logout"), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
