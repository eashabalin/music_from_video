package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"musicFromVideo/pkg/downloader"
)

const (
	commandStart = "start"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}
		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message.Chat.ID)
	default:
		return b.handleUnknownCommand(message.Chat.ID)

	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if b.downloader.IsValidURL(message.Text) {
		go func() {
			errChan := make(chan error, 1)
			go b.handleURL(message.Chat.ID, message.Text, errChan)
			err := <-errChan
			if err != nil {
				b.handleError(message.Chat.ID, err)
			}
		}()
		return nil
	}
	return errInvalidURL
}

func (b *Bot) handleStartCommand(chatID int64) error {
	if err := b.sendMessage(chatID, b.messages.Start); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUnknownCommand(chatID int64) error {
	return b.sendMessage(chatID, b.messages.UnknownCommand)
}

func (b *Bot) handleURL(chatID int64, url string, errChan chan error) {
	defer close(errChan)
	err := b.sendMessage(chatID, b.messages.Loading)
	if err != nil {
		errChan <- err
		return
	}
	filename, err := b.downloader.Download(url)
	if err != nil {
		if err == downloader.ErrorDurationTooLong {
			errChan <- errDurationTooLong
		} else {
			errChan <- errFailedToDownload
		}
	}
	if err = b.sendAudio(chatID, filename); err != nil {
		errChan <- errFailedToSend
		return
	}
}
