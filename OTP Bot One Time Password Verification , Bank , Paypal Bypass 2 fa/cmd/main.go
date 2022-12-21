package main

import (
	"context"
	"os"
	"os/signal"

	"telegram/internal/controller"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go controller.StartTelegramBot(ctx)

	<-c
	cancel()
}
