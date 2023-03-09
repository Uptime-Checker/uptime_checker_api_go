package main

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	if err := config.LoadConfig(".env"); err != nil {
		panic("config load failed")
	}

	// Setup Logger
	lgr.SetupLogger()

	// Setup Cache
	cache.SetupCache()

	// Setup Database
	if err := infra.ConnectDatabase(!config.IsProd); err != nil {
		panic("database connection failed")
	}
	lgr.Default.Print("database connected")
	web.Setup()
}
