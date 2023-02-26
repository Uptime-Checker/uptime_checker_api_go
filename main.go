package main

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	if err := config.LoadConfig(".env"); err != nil {
		panic("Config load failed")
	}

	web.Setup()
}
