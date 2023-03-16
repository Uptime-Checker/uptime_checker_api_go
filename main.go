package main

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())
	ctx = pkg.NewTracingID(ctx)
	tracingID := pkg.GetTracingID(ctx)

	// Setup Logger
	lgr.SetupLogger()

	if err := config.LoadConfig(".env"); err != nil {
		panic(err)
	}
	lgr.Default.Print(tracingID, "configuration loaded")

	// Setup Cache
	cache.SetupCache()
	lgr.Default.Print(tracingID, "cache started")

	// Setup Database
	if err := infra.ConnectDatabase(ctx, !config.IsProd); err != nil {
		panic(err)
	}
	lgr.Default.Print(tracingID, "database connected")
	web.Setup(ctx, shutdown)
}
