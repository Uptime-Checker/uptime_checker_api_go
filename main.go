package main

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/web"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())
	ctx = pkg.NewTracingID(ctx)
	tracingID := pkg.GetTracingID(ctx)

	if err := config.LoadConfig(".env"); err != nil {
		panic(err)
	}

	// Setup logger
	lgr.SetupLogger()

	// Setup Cache
	cache.SetupCache()
	lgr.Print(tracingID, "cache started")

	// Setup Database
	if err := infra.ConnectDatabase(ctx, !config.IsProd); err != nil {
		panic(err)
	}
	lgr.Print(tracingID, "database connected")
	web.Setup(ctx, shutdown)
}
