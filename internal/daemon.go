package internal

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
	"test/pkg"
	"time"
)

var (
	cDaemon *Daemon
	menu    = telebot.ReplyMarkup{}

	btnStart = menu.Text("/start")
	btnStop  = menu.Text("/stop")
)

type Daemon struct {
	Compliments *Compliments
	*telebot.Bot
	*logrus.Logger
}

func NewComplimentsDaemon(bot *telebot.Bot, logger *logrus.Logger) *Daemon {
	if cDaemon == nil {
		cDaemon = &Daemon{
			Compliments: &Compliments{},
			Bot:         bot,
			Logger:      logger,
		}
	}

	return cDaemon
}

func Add(compliment *Compliment) {
	cDaemon.Add(compliment)
}

func Delete(u telebot.Recipient) {
	cDaemon.Delete(u)
}

func (d *Daemon) Add(compliment *Compliment) {
	if d.Find(compliment.UserName) != nil {
		d.Infof("Compliments daemon user %s already exist", compliment.UserName.Recipient())
		return
	}

	*d.Compliments = append(*d.Compliments, compliment)
}

func (d *Daemon) Find(u telebot.Recipient) *Compliment {
	for _, compliment := range *d.Compliments {
		if compliment.UserName.Recipient() == u.Recipient() {
			return compliment
		}
	}
	return nil
}

func (d *Daemon) Delete(u telebot.Recipient) {
	for i, compliment := range *d.Compliments {
		if compliment.UserName.Recipient() == u.Recipient() {
			*d.Compliments = append((*d.Compliments)[:i], (*d.Compliments)[i+1:]...)
			compliment.NeedStop = true

			d.Infof("Compliments daemon delete user %s", u.Recipient())
		}
	}
}

func (d *Daemon) RunDaemon() {
	d.Info("Compliments daemon started")

	go func() {
		for {
			d.Infof("Compliments daemon users count: %d", len(*d.Compliments))
			time.Sleep(10 * time.Second)
		}

	}()

	for {
		for _, compliment := range *d.Compliments {
			if compliment == nil {
				time.Sleep(1 * time.Second)
				continue
			}

			if compliment.IsStarted == false {
				d.Infof("Compliments daemon start user %s", compliment.UserName.Recipient())

				go d.sendCompliment(compliment)
				compliment.IsStarted = true
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (d *Daemon) sendCompliment(u *Compliment) {
	menu.Reply(
		menu.Row(btnStart),
		menu.Row(btnStop),
	)

	for {
		if u.NeedStop {
			break
		}

		compliment, _ := pkg.GetCompliment()
		if _, err := d.Send(u.UserName, compliment.Text, &menu); err != nil {
			d.Errorf("Compliments daemon error: %s", err)
		}
		time.Sleep(1 * time.Second)
	}
}
