package internal

import (
	"gopkg.in/telebot.v3"
)

type Compliments []*Compliment

type Compliment struct {
	UserName  telebot.Recipient
	IsStarted bool
	NeedStop  bool
}

func (cmd *Command) Compliment(c telebot.Context) error {

	Add(&Compliment{
		UserName:  c.Sender(),
		IsStarted: false,
		NeedStop:  false,
	})

	return nil
}

func (cmd *Command) StopCompliment(c telebot.Context) error {
	Delete(c.Sender())

	menu.RemoveKeyboard = true
	return c.Reply("Хорошо", &menu)
}
