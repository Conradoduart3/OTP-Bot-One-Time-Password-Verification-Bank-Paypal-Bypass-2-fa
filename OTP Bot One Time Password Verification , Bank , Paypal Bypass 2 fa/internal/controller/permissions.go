package controller

import (
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewPermission() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewPermission clicked")

		if err := UsersInlineKeyboard(MenuIn); err != nil {
			c.Send(entity.TextInternalError)
		}

		utils.AddUserState(c.Chat().ID, entity.StateChooseUser, c.Message().ID+1)

		return c.Send(entity.TextChooseUserMsg, MenuIn)
	}
}
