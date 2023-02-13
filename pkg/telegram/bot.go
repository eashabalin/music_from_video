package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	username string
}

func NewBot(token string, username string) *Bot {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return &Bot{bot: botAPI, username: username}
}

func (b *Bot) Run() error {
	b.bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			return err
		}

	}
	return nil
}

func (b *Bot) sendAudio(chatID int64, filename string) error {
	path := "downloads/" + filename
	defer os.Remove(path)

	file := tgbotapi.FilePath(path)
	audioCfg := tgbotapi.NewAudio(chatID, file)
	audioCfg.Caption = "Downloaded via @" + b.username
	if _, err := b.bot.Send(audioCfg); err != nil {
		return errors.New("error sending audio")
	}
	return nil
}
