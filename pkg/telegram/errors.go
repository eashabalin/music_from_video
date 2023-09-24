package telegram

import (
	"errors"
	"fmt"
)

var (
	errInvalidURL       = errors.New("invalid url")
	errDurationTooLong  = errors.New("duration too long")
	errFailedToDownload = errors.New("failed to download")
	errFailedToSend     = errors.New("failed to send")
)

func (b *Bot) handleError(chatID int64, err error) {
	var messageText string
	fmt.Println(err)

	switch err {
	case errInvalidURL:
		messageText = b.messages.InvalidURL
	case errDurationTooLong:
		messageText = fmt.Sprintf(b.messages.DurationTooLong+"\n", int(b.maxVideoDuration.Minutes()))
	case errFailedToDownload:
		messageText = b.messages.FailedToDownload
	case errFailedToSend:
		messageText = b.messages.FailedToSend
	default:
		messageText = b.messages.Errors.Default
	}

	b.sendMessage(chatID, messageText)
}
