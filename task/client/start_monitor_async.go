package client

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

func StartMonitorAsync(ctx context.Context, monitorID int64) error {
	body := task.StartMonitorTaskPayload{
		MonitorID: monitorID,
	}

	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, 1, "scheduling start monitor for monitor", monitorID)

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	t := asynq.NewTask(worker.TaskStartMonitor, payload)
	info, err := worker.AsynqEnqueue(ctx, t)
	if err != nil {
		return err
	}
	lgr.Print(tid, 2, "start monitor task sent", info.ID, info.State.String())
	return nil
}
