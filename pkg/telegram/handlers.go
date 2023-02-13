package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if message.IsCommand() && message.Command() == "start" {
		if err := b.sendMessage(message.Chat.ID, "Привет! Отправь ссылку на видео на YouTube, из которого хочешь достать музыку."); err != nil {
			return err
		}
		return nil
	}

	if b.downloader.IsValidURL(message.Text) {
		b.sendMessage(message.Chat.ID, "Загрузка...")
		filename, err := b.downloader.Download(message.Text)
		if err != nil {
			return fmt.Errorf("failed to download: %w\n", err)
		}
		if err = b.sendAudio(message.Chat.ID, filename); err != nil {
			return fmt.Errorf("failed to send audio: %w\n", err)
		}

		return nil
	}

	if err := b.sendMessage(message.Chat.ID, "Не могу прочитать ссылку на видео."); err != nil {
		return err
	}

	return nil
}
