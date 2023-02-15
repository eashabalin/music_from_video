package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	YamlFileCfg
	EnvVars
}

type YamlFileCfg struct {
	BotUsername        string `yaml:"bot_username"`
	MaxDurationMin     int    `yaml:"max_video_duration_min"`
	MaxDownloadTimeSec int    `yaml:"max_download_time_sec"`
}

type EnvVars struct {
	Token string
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
	yamlFile, err := os.ReadFile("pkg/config/config.yaml")
	if err != nil {
		return fmt.Errorf("error reading config.yaml: %v\n", err)
	}

	err = yaml.Unmarshal(yamlFile, &c.YamlFileCfg)
	if err != nil {
		return fmt.Errorf("error unmarshaling config.yaml: %w\n", err)
	}

	c.Token = os.Getenv("TG_BOT_TOKEN")
	return nil
}
