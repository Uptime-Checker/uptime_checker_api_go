package task

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type SyncProductsTask struct{}

func NewSyncProductsTask() *SyncProductsTask {
	return &SyncProductsTask{}
}

func (s SyncProductsTask) Do(ctx context.Context, tx *sql.Tx) {
	tid := pkg.GetTracingID(ctx)

	lgr.Print(tid, 1, "running SyncProductsTask")
}
