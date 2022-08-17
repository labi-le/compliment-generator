package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"test/internal"
	"test/internal/daemon"
	"time"
)

var (
	accessToken string
)

func init() {
	flag.StringVar(&accessToken, "access-token", "", "Telegram access token")
}

func main() {
	flag.Parse()

	if accessToken == "" {
		flag.PrintDefaults()
		return
	}

	app := internal.NewApp(
		accessToken,
		logrus.DebugLevel,
		&tele.LongPoller{Timeout: 10 * time.Second},
	)

	go func() {
		daemon.NewDaemon(app.Logger).RunDaemon()
	}()

	app.Run()
}
