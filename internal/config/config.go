package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-required:"true"`
	ServerConfig `yaml:"server"`
	DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Url string `yaml:"url" env:"SERVER_URL"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, proceeding with environment variables.")
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file does not exist : %s", configPath))
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(fmt.Sprintf("cannot read config : %s", err))
	}
	return &cfg
}
