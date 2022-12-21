package controller

import (
	"telegram/config"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func ShowMyId() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnMyId clicked by ", c.Sender().Username)

		message := utils.GetId(c.Message())

		if message != "" {
			return c.Send(message)
		}

		return nil
	}
}

func NewUser() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewUser clicked by ", c.Sender().Username)

		MenuIn.Inline(
			MenuIn.Row(BtnShowUsers, BtnAddUser),
		)

		return c.Send(entity.TextAddUser, MenuIn)
	}
}

func AddUser() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnAddUser clicked by ", c.Sender().Username)

		utils.AddUserState(c.Chat().ID, entity.StateAddUserEmail, c.Message().ID+2)

		c.Send(entity.TextSendEmailMsg)

		return c.Respond()
	}
}

func ShowUsers() tb.HandlerFunc {
	return func(c tb.Context) error {
		var data entity.ResponseData

		log.Info("BtnShowUsers clicked by ", c.Sender().Username)

		if err := utils.GetStruct(config.Args.USERS_URL, &data); err != nil {
			log.Error(err)
			c.Send(entity.TextInternalError)

			return c.Respond()
		}

		users := utils.GetUsersString(data.Users)

		MenuIn.Inline(
			MenuIn.Row(BtnShowUsers, BtnAddUser),
		)

		c.Send(users, MenuIn)

		return c.Respond()
	}
}
