package main

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	StorageOption string `mapstructure:"storage_option"`
	BadgerDBPath  string `mapstructure:"badger_db_path"`
	RedisAddr     string `mapstructure:"redis_addr"`
	UserAddress   string `mapstructure:"user_address"`
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %s", err)
	}
}
