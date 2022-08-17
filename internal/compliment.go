package internal

import (
	"gopkg.in/telebot.v3"
	"test/internal/daemon"
	"test/pkg"
	"time"
)

var (
	menu = telebot.ReplyMarkup{}

	btnStart = menu.Text("/start")
	btnStop  = menu.Text("/stop")
)

type Compliments []*Compliment

type Compliment struct {
	UserName  telebot.Recipient
	IsStarted bool
	NeedStop  bool
}

func genGoal(c *telebot.Bot, r telebot.Recipient) daemon.Goal {
	return func(t *daemon.Task) error {
		menu.Reply(
			menu.Row(btnStart),
			menu.Row(btnStop),
		)

		compliment, _ := pkg.GetCompliment()
		_, err := c.Send(r, compliment.Text, &menu)
		if err != nil {
			return err
		}

		return nil
	}
}

func (cmd *Command) Compliment(c telebot.Context) error {
	daemon.Add(&daemon.Task{
		ID:          c.Sender().Recipient(),
		Goal:        genGoal(c.Bot(), c.Sender()),
		IsStarted:   false,
		NeedStop:    false,
		RepeatEvery: 1 * time.Second,
	})

	return nil
}

func (cmd *Command) StopCompliment(c telebot.Context) error {
	daemon.Delete(c.Sender().Recipient())

	menu.RemoveKeyboard = true
	return c.Reply("Хорошо", &menu)
}
