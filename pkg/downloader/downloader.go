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

func (d *Downloader) Download(url string) (string, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), d.maxDownloadTime)

	cmd := exec.CommandContext(ctx, "youtube-dl", "-x", "--audio-format", "mp3", url, "-o", "downloads/"+filename+".%(ext)s", "--no-playlist")
	defer cancel()

	data, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("combined output error: %w\n", err)
		return "", err
	}

	if strings.Contains(string(data), "ERROR") {
		err = d.deleteFile("filename")
		if err != nil {
			err = fmt.Errorf("file %s wasn't deleted: %w\n", filename, err)
			return "", err
		}
		err = fmt.Errorf("error downloading video with youtube-dl, output: %s", string(data))
		return "", err
	}
	return filename + ".mp3", nil
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
