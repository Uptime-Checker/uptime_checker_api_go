package infra

import (
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
)

func SetupNewRelic() (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(config.App.NewRelicApp),
		newrelic.ConfigLicense(config.App.NewRelicLicense),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
}
