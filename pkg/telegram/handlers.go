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
		if err := b.handleStartCommand(message.Chat.ID); err != nil {
			return err
		}
	}

	if b.downloader.IsValidURL(message.Text) {
		go func() {
			errChan := make(chan error, 1)
			go b.handleURL(message.Chat.ID, message.Text, errChan)
			err := <-errChan
			fmt.Println("error: ", err)
			if err == downloader.ErrorDurationTooLong {
				b.sendMessage(message.Chat.ID, fmt.Sprintf("Видео должно быть короче %.0f минут.", b.maxVideoDuration.Minutes()))
				return
			}
			if err != nil {
				fmt.Println("error: ", err)
				return
			}
			b.sendMessage(message.Chat.ID, "Не получается обработать эту ссылку.")
		}()
		return nil
	}

	if err := b.sendMessage(message.Chat.ID, "Не могу прочитать ссылку на видео."); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleStartCommand(chatID int64) error {
	if err := b.sendMessage(chatID, "Привет! Отправь ссылку на видео на YouTube, из которого хочешь достать музыку."); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleURL(chatID int64, url string, errChan chan error) {
	defer close(errChan)
	err := b.sendMessage(chatID, "Загрузка...")
	if err != nil {
		errChan <- err
		return
	}
	filename, err := b.downloader.Download(url)
	if err != nil {
		errChan <- fmt.Errorf("failed to download: %w\n", err)
		return
	}
	if err = b.sendAudio(chatID, filename); err != nil {
		errChan <- fmt.Errorf("failed to send audio: %w\n", err)
		return
	}
}
