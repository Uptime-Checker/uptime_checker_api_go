package config

import "github.com/spf13/viper"

type Config struct {
	Host    string `json:"HOST"`
	Port    string `json:"PORT"`
	Release string `json:"RELEASE"`
}

var App Config

func LoadConfig(path string) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&App); err != nil {
		return err
	}
	return nil
}
