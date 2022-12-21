package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"telegram/config"
	"telegram/internal/entity"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/telebot.v3"
)

var (
	BackendClient = &http.Client{Timeout: 10 * time.Second}
	UserStates    = make(map[int64]map[string]int)    // map['chatID'] = map'btnAddOrigin' = 'message.ID'
	AddPermStates = make(map[int64]map[string]string) // map['chatID'] = map'userEmail' = 'userEmail', map'host' = 'host'
	NewUserStates = make(map[int64]map[string]string) // map['chatID'] = map'email' = 'userEmail', map'chatID' = 'ID'
)

func GetId(m *tb.Message) string {
	var userChat entity.Recipient

	if !m.Private() {
		log.Error("Error: chat is not private")
		return ""
	}

	userChat.ID = int(m.Chat.ID)

	message := "Сообщи этот ID админу для авторизации: " + userChat.Recipient()

	return message
}

func IsAdmin(chatId int) error {
	req, err := setRequest(entity.Payload{
		ChatID: strconv.Itoa(chatId),
	}, config.Args.AUTH_URL)
	if err != nil {
		log.Error("Error setting request for admin: ", err)

		return err
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error checking admin from Auth Backend: ", err)

		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("User with ID = %d is not an admin", chatId)
	}

	return nil
}

func setRequest(payload entity.Payload, url string) (http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating json body: %w", err)
	}

	responseBody := bytes.NewReader(body)

	req, err := http.NewRequest("POST", url, responseBody)
	if err != nil {
		return http.Request{}, fmt.Errorf("Error creating request to Backend: %w", err)
	}

	setHeaders(req)

	return *req, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("X-Green-Origin", "telegram-bot")
	req.Header.Set("Api-Key", config.Args.API_KEY)
}

func GetUsersString(l []entity.User) string {
	var str string

	for _, elem := range l {
		str += elem.User + "\n"
	}

	return str
}

func GetOriginString(l []entity.Origs) string {
	var str string

	for _, elem := range l {
		str += elem.Origin + "\n"
	}

	return str
}

func GetStruct(url string, data *entity.ResponseData) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error making GET request to %s: %w", url, err)
	}

	setHeaders(req)

	resp, err := BackendClient.Do(req)
	if err != nil {
		log.Error("\nError calling GET ", url, ": ", err)

		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("\nError in response from GET %s", url)
	}

	err = readJson(resp.Body, data)
	if err != nil {
		return fmt.Errorf("\nError in readJson from %s: %w", url, err)
	}

	return err
}

func readJson(resp io.ReadCloser, data *entity.ResponseData) error {
	body, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("\nError reading response from GET /origins")
	}

	err = json.Unmarshal([]byte(body), data)
	if err != nil {
		return fmt.Errorf("\nError unmarshalling response JSON %w", err)
	}

	return err
}

func AddUserState(chatID int64, state string, msgID int) {
	if _, userExist := UserStates[chatID]; !userExist {
		UserStates[chatID] = make(map[string]int)
	}

	UserStates[chatID][state] = msgID
}

func AddOriginToBackend(origin string) (string, error) {
	var data entity.ResponseData

	req, err := setRequest(entity.Payload{
		Origin: origin,
	}, config.Args.NEW_ORIGIN_URL)
	if err != nil {
		log.Error("Error setting request for .ValidateOrigin: ", err)

		return entity.TextInternalError, err
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error creating new Origin: ", err)

		return entity.TextInternalError, err
	}

	if resp.StatusCode == http.StatusInternalServerError ||
		resp.StatusCode == http.StatusBadRequest {
		return entity.TextInternalError,
			fmt.Errorf("Error: received error status code from backend: %d",
				resp.StatusCode)
	}

	err = readJson(resp.Body, &data)
	if err != nil {
		log.Error("Error in .ValidateOrigin: ", err)

		return entity.TextInternalError, err
	}

	return data.Response, err
}

func AddPermState(chatID int64, stage string, value string) {
	if _, userExist := AddPermStates[chatID]; !userExist {
		AddPermStates[chatID] = make(map[string]string)
	}

	AddPermStates[chatID][stage] = value
}

func AddNewUserState(chatID int64, param string, value string) {
	if _, userExist := NewUserStates[chatID]; !userExist {
		NewUserStates[chatID] = make(map[string]string)
	}

	NewUserStates[chatID][param] = value

	log.Debug("current NewUserStates: ", NewUserStates)
}

func AddPermission(email string, origin string) (int, error) {
	req, err := setRequest(entity.Payload{
		Email:  email,
		Origin: origin,
	}, config.Args.ADD_PERMISSION_URL)
	if err != nil {
		log.Error("Error setting request for .AddPermission: ", err)

		return 0, err
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error sending request from .AddPermission to Backend: ", err)

		return 0, err
	}

	return resp.StatusCode, nil
}

func NewUserToBackend(email string, chatID string) (int, error) {
	req, err := setRequest(entity.Payload{
		Email:  email,
		ChatID: chatID,
	}, config.Args.NEW_USER_URL)
	if err != nil {
		log.Error("Error setting request for .NewUserToBackend: ", err)

		return 0, err
	}

	resp, err := BackendClient.Do(&req)
	if err != nil {
		log.Error("Error sending request from .AddPermission to Backend: ", err)

		return 0, err
	}

	return resp.StatusCode, nil
}

func CheckInput(s string) bool {
	if s == "" {
		return false
	}

	for _, ch := range s {
		if ch == 32 || ch == 92 || ch == 10 {
			return false
		}
	}

	return true
}
