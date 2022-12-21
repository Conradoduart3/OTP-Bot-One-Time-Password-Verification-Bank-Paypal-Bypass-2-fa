package controller

import (
	"net/http"
	"strings"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

var (
	Menu             = &tb.ReplyMarkup{}
	MenuIn           = &tb.ReplyMarkup{}
	BtnNewUser       = Menu.Text(entity.TextNewUserBtn)
	BtnNewOrigin     = Menu.Text(entity.TextNewOriginBtn)
	BtnNewPermission = Menu.Text(entity.TextNewPermissionBtn)
	BtnMyId          = Menu.Text(entity.TextMyID)
	BtnShowOrigins   = MenuIn.Data(entity.TextCurrentOriginsBtn, "origins")
	BtnShowUsers     = MenuIn.Data(entity.TextShowUsersBtn, "users")
	BtnAddOrigin     = MenuIn.Data(entity.TextAddOriginBtn, "newOrigin")
	BtnAddUser       = MenuIn.Data(entity.TextAddUserBtn, "newUser")
)

func OnStart() tb.HandlerFunc {
	return func(c tb.Context) error {
		message := utils.GetId(c.Message())

		if message != "" {
			if err := utils.IsAdmin(int(c.Sender().ID)); err != nil {
				log.Info(err)
				Menu.Reply(
					Menu.Row(BtnMyId),
				)
			} else {
				log.Info("Admin user signed in: ", c.Sender().Username)

				Menu.Reply(
					Menu.Row(BtnMyId, BtnNewUser),
					Menu.Row(BtnNewOrigin, BtnNewPermission),
				)
			}
		}

		return c.Send(message, Menu)
	}
}

func OnText() tb.HandlerFunc {
	return func(c tb.Context) error {
		var (
			state string
			msgID int
		)

		if _, userExist := utils.UserStates[c.Chat().ID]; !userExist {
			utils.UserStates[c.Chat().ID] = make(map[string]int)
		}

		for state, msgID = range utils.UserStates[c.Chat().ID] {
			if msgID == c.Message().ID {
				switch state {
				case entity.StateAddOrigin:
					return AddingOrigin(c)
				case entity.StateAddUserEmail:
					return InsertingEmail(c)
				case entity.StateAddUserChatID:
					return InsertingChatID(c)
				}
			}
		}

		return c.Send(entity.TextUnknownMsg)
	}
}

func OnCallback() tb.HandlerFunc {
	return func(c tb.Context) error {
		var (
			data  string
			state string
			msgID int
		)

		data = strings.TrimPrefix(c.Callback().Data, "\f")

		if _, userExist := utils.UserStates[c.Chat().ID]; !userExist {
			utils.UserStates[c.Chat().ID] = make(map[string]int)
		}

		for state, msgID = range utils.UserStates[c.Chat().ID] {
			if msgID == c.Message().ID {
				switch state {
				case entity.StateChooseUser:
					return ChoosingUser(c, data)
				case entity.StateChooseHost:
					return ChoosingHost(c, data)
				}
			}
		}

		c.Send(entity.TextUnknownMsg)

		return c.Respond()
	}
}

func ChoosingUser(c tb.Context, data string) error {
	utils.AddPermState(c.Chat().ID, "email", data)
	msg := entity.TextChooseOriginMsg + data

	if err := OriginsInlineKeyboard(MenuIn); err != nil {
		c.Send(entity.TextInternalError)

		return c.Respond()
	}

	utils.AddUserState(c.Chat().ID, entity.StateChooseHost, c.Message().ID+1)
	c.Send(msg, MenuIn)

	return c.Respond()
}

func ChoosingHost(c tb.Context, data string) error {
	email := utils.AddPermStates[c.Chat().ID]["email"]

	status, err := utils.AddPermission(email, data)
	if err != nil {
		log.Error("Error adding permission: ", err)
		c.Send(entity.TextInternalError)

		return c.Respond()
	}

	msg := entity.TextInternalError
	if status == http.StatusConflict {
		msg = "У пользователя " + email + " уже есть доступ к " + data
	}

	if status == http.StatusCreated {
		msg = "Выдали пользователю " + email + " доступ к сервису " + data
		log.Info("User ", c.Sender().Username, " added permission for ", email, " to ", data)

		delete(utils.AddPermStates[c.Chat().ID], "email")
	}

	c.Send(msg)

	return c.Respond()
}

func AddingOrigin(c tb.Context) error {
	MenuIn.Inline(
		MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
	)

	if valid := utils.CheckInput(c.Message().Text); !valid {
		return c.Send(entity.TextSpacesNotAllowed, MenuIn)
	}

	data := strings.ToLower(c.Message().Text)

	msg, err := utils.AddOriginToBackend(data)
	if err == nil {
		log.Info("User ", c.Sender().Username, " added origin ", data)
	}

	return c.Send(msg, MenuIn)
}

func InsertingEmail(c tb.Context) error {
	if valid := utils.CheckInput(c.Message().Text); !valid {
		MenuIn.Inline(
			MenuIn.Row(BtnShowUsers, BtnAddUser),
		)

		return c.Send(entity.TextSpacesNotAllowed, MenuIn)
	}

	data := strings.ToLower(c.Message().Text)

	utils.AddNewUserState(c.Chat().ID, "email", data)
	utils.AddUserState(c.Chat().ID, entity.StateAddUserChatID, c.Message().ID+2)

	msg := entity.TextSendChatIDMsg

	return c.Send(msg)
}

func InsertingChatID(c tb.Context) error {
	MenuIn.Inline(
		MenuIn.Row(BtnShowUsers, BtnAddUser),
	)

	if valid := utils.CheckInput(c.Message().Text); !valid {
		return c.Send(entity.TextSpacesNotAllowed, MenuIn)
	}

	email := utils.NewUserStates[c.Chat().ID]["email"]

	status, err := utils.NewUserToBackend(email, c.Message().Text)
	if err != nil {
		return c.Send(entity.TextInternalError, MenuIn)
	}

	delete(utils.NewUserStates[c.Chat().ID], "email")

	msg := entity.TextInternalError

	if status == http.StatusCreated {
		log.Info("New user ", email, " was added by ", c.Sender().Username)
		msg = "Пользователь " + email + " успешно добавлен. Можно проверять :)"
	}

	if status == http.StatusConflict {
		msg = "Пользователь " + email + " уже существует"
	}

	return c.Send(msg, MenuIn)
}
