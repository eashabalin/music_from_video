package main

import (
	"musicFromVideo/pkg/config"
	tg "musicFromVideo/pkg/telegram"
	"time"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	bot, err := tg.NewBot(
		cfg.Token, cfg.BotUsername,
		time.Duration(cfg.MaxDurationMin)*time.Minute,
		time.Duration(cfg.MaxDownloadTimeSec)*time.Second,
		cfg.Messages,
	)
	if err != nil {
		panic(err)
	}

	if err := bot.Run(); err != nil {
		panic(err)
	}

}
