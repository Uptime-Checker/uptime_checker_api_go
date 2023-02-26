package main

import "github.com/spf13/viper"

type Config struct {
	Host    string `json:"HOST"`
	Port    string `json:"PORT"`
	Release string `json:"RELEASE"`
}

func loadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
