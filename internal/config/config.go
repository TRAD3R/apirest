package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug bool `yaml:"is_debug" env-default:"false"`
	DB      struct {
		URL string `yaml:"url" required:"true"`
	} `yaml:"db"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			log.Fatalf("Error reading config: %v", err)
		}
	})

	return instance
}
