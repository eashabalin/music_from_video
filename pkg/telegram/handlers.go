package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"musicFromVideo/pkg/downloader"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	d, err := downloader.NewDownloader()
	if err != nil {
		return err
	}
	filename := message.Text
	if d.IsValidURL(message.Text) {
		filename, err = d.Download(message.Text)
		if err != nil {
			return err
		}
		if err := b.sendAudio(message.Chat.ID, filename); err != nil {
			panic(err)
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.username)

	b.bot.Send(msg)

	return nil
}
