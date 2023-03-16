package main

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())
	if err := config.LoadConfig(".env"); err != nil {
		panic(err)
	}

	// Setup Logger
	lgr.SetupLogger()

	// Setup Cache
	cache.SetupCache()

	// Setup Database
	if err := infra.ConnectDatabase(ctx, !config.IsProd); err != nil {
		panic(err)
	}
	lgr.Default.Print("database connected")
	web.Setup(ctx, shutdown)
}
