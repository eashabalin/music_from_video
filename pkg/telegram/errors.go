package telegram

import (
	"errors"
	"fmt"
)

var (
	errInvalidURL       = errors.New("invalid url")
	errDurationTooLong  = errors.New("duration too long")
	errUnknownCommand   = errors.New("unknown command")
	errFailedToDownload = errors.New("failed to download")
	errFailedToSend     = errors.New("failed to send audio")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string
	fmt.Println(err)

	switch err {
	case errInvalidURL:
		messageText = "Не получается обработать эту ссылку."
	case errDurationTooLong:
		messageText = fmt.Sprintf("Видео должно быть короче %.0f минут.", b.maxVideoDuration.Minutes())
	case errUnknownCommand:
		messageText = "Не знаю такую команду :("
	case errFailedToDownload:
		messageText = "Не удалось скачать аудио."
	case errFailedToSend:
		messageText = "Не удалось отправить аудио."
	default:
		messageText = "Произошла неизвестная ошибка."
	}

	b.sendMessage(chatID, messageText)
}
