package task

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type SyncProductsTask struct{}

func NewSyncProductsTask() *SyncProductsTask {
	return &SyncProductsTask{}
}

func (s *SyncProductsTask) Run() {
	ctx := pkg.NewTracingID(context.Background())
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	log.Default.Print(tid, 1, "Running SyncProductsTask")
}
