package controller

import (
	"context"
	"time"

	"telegram/config"

	log "github.com/sirupsen/logrus"

	tb "gopkg.in/telebot.v3"
)

var Bot = &tb.Bot{}

func StartTelegramBot(ctx context.Context) {
	config.Validate()

	level, err := log.ParseLevel(config.Args.LOG_LEVEL)
	if err != nil {
		log.Fatal(err)
	}

	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(level)

	settings := tb.Settings{
		Token: config.Args.TG_BOT_KEY,
		Poller: &tb.LongPoller{
			Timeout: 1 * time.Second,
		},
	}

	Bot, _ = tb.NewBot(settings)

	Bot.Handle("/start", OnStart())
	Bot.Handle(tb.OnText, OnText())
	Bot.Handle(tb.OnCallback, OnCallback())

	// Buttons
	Bot.Handle(&BtnMyId, ShowMyId())
	Bot.Handle(&BtnNewUser, NewUser())
	Bot.Handle(&BtnNewOrigin, NewOrigin())
	Bot.Handle(&BtnNewPermission, NewPermission())

	// Inline Buttons
	Bot.Handle(&BtnShowOrigins, ShowOrigins())
	Bot.Handle(&BtnAddOrigin, AddOrigin())
	Bot.Handle(&BtnShowUsers, ShowUsers())
	Bot.Handle(&BtnAddUser, AddUser())

	go func() {
		Bot.Start()
	}()

	log.Info("Telegram Bot started")

	<-ctx.Done()

	log.Info("Telegram Bot stopped")
	Bot.Stop()
}
