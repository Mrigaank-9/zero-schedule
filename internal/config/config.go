package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	Workers int64  `yaml:"workers" env-required:"true"`
	Address string `yaml:"address" env-required:"true"`
}
type Config struct {
	Version string `yaml:"version" env-required:"true"`
	Env     string `yaml:"env" env-required:"true"`
	App     `yaml:"app"`
}

func MustLoadConfig() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags := flag.String("config", "", "path to the config file")
		flag.Parse()

		configPath = *flags
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file doesn't exists : %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("config file not able to load: %s", err.Error())
	}

	return &cfg
}
