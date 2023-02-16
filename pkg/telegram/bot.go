package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"musicFromVideo/pkg/downloader"
	"os"
	"time"
)

type Bot struct {
	bot              *tgbotapi.BotAPI
	username         string
	downloader       downloader.Downloader
	maxVideoDuration time.Duration
	maxDownloadTime  time.Duration
}

func NewBot(token string, username string, maxVideoDuration, maxDownloadTime time.Duration) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		err = fmt.Errorf("failed to create bot: %w\n", err)
		return nil, err
	}
	downloaderService, err := downloader.NewDownloader(maxVideoDuration, maxDownloadTime)
	if err != nil {
		err = fmt.Errorf("failed to create downloader: %w\n", err)
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("failed to create poller: %w\n", err)
		return nil, err
	}
	return &Bot{
		bot:              botAPI,
		username:         username,
		downloader:       *downloaderService,
		maxVideoDuration: maxVideoDuration,
		maxDownloadTime:  maxDownloadTime,
	}, nil
}

func (b *Bot) Run() error {
	b.bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := b.bot.GetUpdatesChan(updateConfig)

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) sendAudio(chatID int64, filename string) error {
	path := "downloads/" + filename
	defer os.Remove(path)

	file := tgbotapi.FilePath(path)
	audioCfg := tgbotapi.NewAudio(chatID, file)
	audioCfg.Caption = "Downloaded via @" + b.username
	if _, err := b.bot.Send(audioCfg); err != nil {
		return fmt.Errorf("error sending audio: %w\n", err)
	}
	return nil
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}
