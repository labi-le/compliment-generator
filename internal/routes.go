package internal

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func AddHandlers(b *telebot.Bot, l *logrus.Logger) {
	c := Command{Logger: l}

	b.Use(c.LoggerMiddleware)
	b.Handle("/start", c.Compliment)
	b.Handle("/stop", c.StopCompliment)
}
