package main

import (
	"musicFromVideo/pkg/config"
	tg "musicFromVideo/pkg/telegram"
)

func main() {
	cfg := config.NewConfig()

	bot := tg.NewBot(cfg.Token, cfg.BotUsername)

	if err := bot.Run(); err != nil {
		panic(err)
	}

}
