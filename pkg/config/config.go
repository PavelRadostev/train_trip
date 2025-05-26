package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	DB struct {
		DSN             string        `yaml:"dsn"`
		Schema          string        `yaml:"schema"`
		MaxConns        int32         `yaml:"max_conns"`
		ConnectTimeout  time.Duration `yaml:"connect_timeout"`
		ConnectAttempts int           `yaml:"connect_attempts"`
	} `yaml:"db"`

	Migration struct {
		Dir string `yaml:"migration_dir"`
	} `yaml:"migration"`
}

func Load() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read congig: %s", err)
	}

	return &cfg
}
