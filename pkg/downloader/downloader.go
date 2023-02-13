package downloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const expression = "^(http(s)?:\\/\\/)?((w){3}.)?(music\\.)?youtu(be|.be)?(\\.com)?\\/.+"

type Downloader struct {
	regexpr regexp.Regexp
}

func NewDownloader() (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}
	return &Downloader{
		regexpr: *r,
	}, nil
}

func (d *Downloader) Download(url string) (string, error) {
	if !d.IsValidURL(url) {
		return "", errors.New("invalid url")
	}

	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return "", err
	}

	filename := video.Title

	if video.Duration.Minutes() > 10 {
		return "", errors.New("duration too long")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*180)

	cmd := exec.CommandContext(ctx, "youtube-dl", "-x", "--audio-format", "mp3", url, "-o", "downloads/"+filename+".%(ext)s", "--no-playlist")
	defer cancel()

	data, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	if strings.Contains(string(data), "ERROR") {
		err = d.deleteFile("filename")
		if err != nil {
			log.Printf("file %s wasn't deleted\n", filename)
		}
		return "", errors.New(fmt.Sprintf("error downloading video with youtube-dl, output: %s", string(data)))
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
