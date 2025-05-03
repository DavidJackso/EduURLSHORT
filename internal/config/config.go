package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"http://localhost:8081"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"10s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoadConfig() *Config { //Must приписываем когда не возвращаем ошибку, а кидаем панику. Использовать только в редких случаях здесь обусловлено тем что мы только запускаем программу.
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}
	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal(err)
	}

	return &config
}
