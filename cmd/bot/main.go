package main

import (
	"musicFromVideo/pkg/config"
	tg "musicFromVideo/pkg/telegram"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	bot, err := tg.NewBot(cfg.Token, cfg.BotUsername)
	if err != nil {
		panic(err)
	}

	if err := bot.Run(); err != nil {
		panic(err)
	}

}
