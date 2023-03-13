package task

import (
	"context"
	"database/sql"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type SyncProductsTask struct{}

func NewSyncProductsTask() *SyncProductsTask {
	return &SyncProductsTask{}
}

func (s SyncProductsTask) Do(tx *sql.Tx) {
	ctx := pkg.NewTracingID(context.Background())
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Default.Print(tid, 1, "Running SyncProductsTask")
}
