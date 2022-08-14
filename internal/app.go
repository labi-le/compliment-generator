package internal

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type App struct {
	Logger      *logrus.Logger
	lp          *telebot.LongPoller
	accessToken string

	Bot *telebot.Bot
}

func NewApp(accessToken string, logLevel logrus.Level, lp *telebot.LongPoller) *App {
	l := logrus.New()
	l.SetLevel(logLevel)
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	b, err := telebot.NewBot(telebot.Settings{
		Token:  accessToken,
		Poller: lp,
	})
	if err != nil {
		l.Fatalf("Error creating bot: %s", err)
	}

	l.Infof("App started with token: %s Nickname: %s ID: %d", accessToken, b.Me.FirstName+" "+b.Me.LastName, b.Me.ID)

	return &App{Logger: l, lp: lp, accessToken: accessToken, Bot: b}

}

func (a *App) Run() error {
	AddHandlers(a.Bot, a.Logger)

	a.Bot.Start()

	return nil
}
