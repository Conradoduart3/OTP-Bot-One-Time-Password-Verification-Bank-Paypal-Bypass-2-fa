package controller

import (
	"telegram/config"
	"telegram/internal/entity"
	"telegram/internal/utils"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

func NewOrigin() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnNewOrigin clicked by ", c.Sender().Username)

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		return c.Send(entity.TextAddOrigin, MenuIn)
	}
}

func ShowOrigins() tb.HandlerFunc {
	return func(c tb.Context) error {
		var data entity.ResponseData

		log.Info("BtnShowOrigins clicked by ", c.Sender().Username)

		if err := utils.GetStruct(config.Args.ORIGIN_URL, &data); err != nil {
			log.Error(err)
			c.Send(entity.TextInternalError)

			return c.Respond()
		}

		hosts := utils.GetOriginString(data.Origins)

		MenuIn.Inline(
			MenuIn.Row(BtnShowOrigins, BtnAddOrigin),
		)

		c.Send(hosts, MenuIn)

		return c.Respond()
	}
}

func AddOrigin() tb.HandlerFunc {
	return func(c tb.Context) error {
		log.Info("BtnAddOrigin clicked by ", c.Sender().Username)

		utils.AddUserState(c.Chat().ID, entity.StateAddOrigin, c.Message().ID+2)

		c.Send(entity.TextSendHostMsg)

		return c.Respond()
	}
}
