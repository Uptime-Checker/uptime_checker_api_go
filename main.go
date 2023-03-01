package main

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	if err := config.LoadConfig(".env"); err != nil {
		panic("config load failed")
	}

	// Setup Logger
	log.SetupLogger()

	// Setup Database
	if err := infra.ConnectDatabase(); err != nil {
		panic("database connection failed")
	}
	log.Default.Print("database connected")
	web.Setup()
}
