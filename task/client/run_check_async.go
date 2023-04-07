package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/task"
)

func RunCheckAsync(ctx context.Context, monitorRegionID int64, runAt time.Time) error {
	body := task.RunCheckTaskPayload{
		MonitorRegionID: monitorRegionID,
	}

	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, 1, "scheduling check run for monitor region", monitorRegionID, "at", runAt.String())

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	job := &gue.Job{Type: worker.TaskRunCheck, RunAt: runAt, Args: payload}
	return worker.SlowWheel.Enqueue(ctx, job)
}

func RunCheckAsyncFast(ctx context.Context, monitorRegionID int64, runAt time.Time) error {
	body := task.RunCheckTaskPayload{
		MonitorRegionID: monitorRegionID,
	}

	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, 1, "scheduling check run for monitor region", monitorRegionID, "at", runAt.String())

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	t := asynq.NewTask(worker.TaskRunCheck, payload)
	info, err := worker.FastWheel.EnqueueContext(ctx, t, asynq.ProcessAt(runAt))
	if err != nil {
		return err
	}
	lgr.Print(tid, 2, "check run task sent", info.ID, info.State.String())
	return nil
}
