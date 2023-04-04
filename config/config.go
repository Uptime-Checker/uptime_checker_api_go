package config

import (
	"github.com/spf13/viper"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type config struct {
	Port       string `mapstructure:"PORT"`
	Release    string `mapstructure:"RELEASE"`
	APIKey     string `mapstructure:"X_API_KEY"`
	JWTKey     string `mapstructure:"JWT_KEY"`
	WorkerPool int    `mapstructure:"WORKER_POOL"`
	FlyRegion  string `mapstructure:"FLY_REGION"`

	Version string `mapstructure:"VERSION"`

	DatabaseHost     string `mapstructure:"DB_HOST"`
	DatabasePort     string `mapstructure:"DB_PORT"`
	DatabaseUser     string `mapstructure:"DB_USER"`
	DatabasePassword string `mapstructure:"DB_PASSWORD"`
	DatabaseSchema   string `mapstructure:"DB_SCHEMA"`

	SentryDSN string `mapstructure:"SENTRY_DSN"`

	NewRelicApp     string `mapstructure:"NEWRELIC_APP"`
	NewRelicLicense string `mapstructure:"NEWRELIC_LICENSE"`

	AxiomOrganizationID string `mapstructure:"AXIOM_ORG_ID"`
	AxiomToken          string `mapstructure:"AXIOM_TOKEN"`
	AxiomDataset        string `mapstructure:"AXIOM_DATASET"`
}

var (
	App    config
	IsProd bool
	JWTKey = []byte(App.JWTKey)
	Region *model.Region
)

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
