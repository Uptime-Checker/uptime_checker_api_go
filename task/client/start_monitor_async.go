package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

func StartMonitor(ctx context.Context, monitorID int64, runAt time.Time) error {
	body := task.StartMonitorTaskPayload{
		MonitorID: monitorID,
	}

	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, "scheduling start monitor for monitor", monitorID)

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	job := &gue.Job{Type: worker.TaskStartMonitor, RunAt: runAt, Args: payload}
	return worker.Wheel.Enqueue(ctx, job)
}
