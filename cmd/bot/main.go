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

	bot := tg.NewBot(cfg.Token, cfg.BotUsername)

	if err := bot.Run(); err != nil {
		panic(err)
	}

}
