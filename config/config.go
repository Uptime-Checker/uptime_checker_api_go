package config

import (
	"github.com/spf13/viper"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	Release string `mapstructure:"RELEASE"`
	APIKey  string `mapstructure:"X_API_KEY"`

	Version string `mapstructure:"VERSION"`

	SentryDSN string `mapstructure:"SENTRY_DSN"`

	NewRelicApp     string `mapstructure:"NEWRELIC_APP"`
	NewRelicLicense string `mapstructure:"NEWRELIC_LICENSE"`
}

var App Config
var IsProd bool

func LoadConfig(path string) error {
	viper.AutomaticEnv()
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&App); err != nil {
		return err
	}
	IsProd = App.Release == string(constant.EnvironmentProd)
	return nil
}
