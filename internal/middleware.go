package internal

import (
	"gopkg.in/telebot.v3"
)

func (cmd *Command) LoggerMiddleware(handlerFunc telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		cmd.Infof(
			"ID:%d (%s) Nickname:%s Message:%s",
			c.Sender().ID,
			c.Sender().Username,
			c.Sender().FirstName+" "+c.Sender().LastName,
			c.Text(),
		)

		return handlerFunc(c)
	}
}
