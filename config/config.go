package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Env *viper.Viper
}

func NewConfig() *Config {
	config := viper.New()
	config.SetConfigFile(".env")
	config.AddConfigPath(".")
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	return &Config{Env: config}
}
