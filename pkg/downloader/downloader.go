package downloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const expression = "^(http(s)?:\\/\\/)?((w){3}.)?(music\\.)?youtu(be|.be)?(\\.com)?\\/.+"

var (
	ErrorDurationTooLong = errors.New("duration too long")
)

type Downloader struct {
	regexpr         regexp.Regexp
	maxDuration     time.Duration
	maxDownloadTime time.Duration
}

func NewDownloader(maxVideoDuration time.Duration, maxDownloadTime time.Duration) (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		err = fmt.Errorf("error compiling regexp: %w\n", err)
		return nil, err
	}
	return &Downloader{
		regexpr:         *r,
		maxDuration:     maxVideoDuration,
		maxDownloadTime: maxDownloadTime,
	}, nil
}

func (d *Downloader) DownloadAudio(url string) (string, error) {
	filename, err := d.checkUrlAndDownload(url, d.callYtDlpForAudio)
	if err != nil {
		return "", err
	}
	return filename + ".mp3", err
}

func (d *Downloader) DownloadVideo(url string) (string, error) {
	filename, err := d.checkUrlAndDownload(url, d.callYtDlpForVideo)
	if err != nil {
		return "", err
	}
	return filename + ".mp4", err
}

func (d *Downloader) checkUrlAndDownload(url string, ytDlpCall YtDlpCall) (string, error) {
	if !d.IsValidURL(url) {
		err := fmt.Errorf("invalid url\n")
		return "", err
	}

	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		err = fmt.Errorf("unable to get video: %w\n", err)
		return "", err
	}

	filename := video.Title

	if video.Duration > d.maxDuration {
		return "", ErrorDurationTooLong
	}

	err = ytDlpCall(url, filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

type YtDlpCall func(string, string) error

func (d *Downloader) callYtDlpForAudio(url string, filename string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.maxDownloadTime)

	cmd := exec.CommandContext(ctx, "yt-dlp", "-x", "--audio-format", "mp3", url, "-o", "downloads/"+filename+".mp3", "--no-playlist")
	defer cancel()

	data, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("yt-dlp output error: %w\n. output: %s\n", err, string(data))
		return err
	}

	return nil
}

func (d *Downloader) callYtDlpForVideo(url string, filename string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.maxDownloadTime)

	cmd := exec.CommandContext(ctx, "yt-dlp", url, "-o", "downloads/"+filename+".mp4", "--format", "mp4")
	defer cancel()

	data, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("yt-dlp output error: %w\n. output: %s\n", err, string(data))
		return err
	}
	return nil
}

func (d *Downloader) IsValidURL(url string) bool {
	return d.regexpr.MatchString(url)
}

func (d *Downloader) deleteFile(name string) error {
	files, err := os.ReadDir("downloads")
	if err != nil {
		return err
	}
	for _, f := range files {
		filename := f.Name()
		if !f.IsDir() && strings.HasPrefix(filename, name) {
			if err = os.Remove("downloads/" + filename); err != nil {
				return err
			}
		}
	}
	return nil
}
