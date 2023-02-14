package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"musicFromVideo/pkg/downloader"
)

func (b *Bot) handleUpdate(update tgbotapi.Update) error {
	if update.Message != nil {
		err := b.handleMessage(update.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if message.IsCommand() && message.Command() == "start" {
		if err := b.sendMessage(message.Chat.ID, "Привет! Отправь ссылку на видео на YouTube, из которого хочешь достать музыку."); err != nil {
			return err
		}
		return nil
	}

	errChan := make(chan error, 1)

	if b.downloader.IsValidURL(message.Text) {
		go func() {
			err := b.sendMessage(message.Chat.ID, "Загрузка...")
			if err != nil {
				errChan <- err
				return
			}
			filename, err := b.downloader.Download(message.Text)
			if err == downloader.ErrorDurationTooLong {
				err = b.sendMessage(message.Chat.ID, "Видео должно быть короче 10 минут.")
				if err != nil {
					errChan <- err
					return
				}
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to download: %w\n", err)
				return
			}
			if err = b.sendAudio(message.Chat.ID, filename); err != nil {
				errChan <- fmt.Errorf("failed to send audio: %w\n", err)
				return
			}
			errChan <- nil
			close(errChan)
		}()
		return <-errChan
	}

	if err := b.sendMessage(message.Chat.ID, "Не могу прочитать ссылку на видео."); err != nil {
		return err
	}

	return nil
}
