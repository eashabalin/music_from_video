package config

import (
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
	BotUsername string `yaml:"bot_username"`
}

type EnvVars struct {
	Token string
}

func NewConfig() *Config {
	var config Config
	config.getConfigs()
	return &config
}

func (c *Config) getConfigs() {
	yamlFile, err := os.ReadFile("pkg/config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &c.YamlFileCfg)

	c.Token = os.Getenv("TG_BOT_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
}
