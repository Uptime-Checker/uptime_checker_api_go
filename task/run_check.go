package task

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type RunCheckTask struct{}

func NewRunCheckTask() *RunCheckTask {
	return &RunCheckTask{}
}

func (s RunCheckTask) Do(ctx context.Context, job *gue.Job) error {
	tid := pkg.GetTracingID(ctx)
	defer sentry.RecoverWithContext(ctx)

	lgr.Default.Print(tid, 1, "Running RunCheckTask")

	return nil
}
