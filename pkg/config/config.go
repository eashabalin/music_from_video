package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Token              string
	BotUsername        string   `yaml:"bot_username"`
	MaxDurationMin     int      `yaml:"max_video_duration_min"`
	MaxDownloadTimeSec int      `yaml:"max_download_time_sec"`
	Messages           Messages `yaml:"messages"`
}

type Messages struct {
	Responses `yaml:"response"`
	Errors    `yaml:"error"`
}

type Responses struct {
	Start          string `yaml:"start"`
	UnknownCommand string `yaml:"unknown_command"`
	Loading        string `yaml:"loading"`
	Video          string `yaml:"video"`
}

type Errors struct {
	InvalidURL       string `yaml:"invalid_url"`
	DurationTooLong  string `yaml:"duration_too_long"`
	FailedToDownload string `yaml:"failed_to_download"`
	FailedToSend     string `yaml:"failed_to_send"`
	Default          string `yaml:"default"`
}

func init() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatal(err)
	//}
}

func NewConfig() (*Config, error) {
	var config Config
	if err := config.getConfigs(); err != nil {
		err = fmt.Errorf("error reading configs: %w\n", err)
		return nil, err
	}
	return &config, nil
}

func (c *Config) getConfigs() error {
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error reading config.yaml: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return fmt.Errorf("error unmarshaling config.yaml: %w\n", err)
	}

	token, ok := os.LookupEnv("TG_BOT_TOKEN")
	if !ok {
		return fmt.Errorf("not found env variable TG_BOT_TOKEN")
	}
	c.Token = token

	return nil
}
