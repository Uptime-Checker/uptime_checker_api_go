package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
)

type RunCheckTaskPayload struct {
	MonitorID int64
}

func RunCheckAsync(ctx context.Context, monitorID int64, runAt time.Time) error {
	body := RunCheckTaskPayload{
		MonitorID: monitorID,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	job := &gue.Job{Type: worker.TaskRunCheck, RunAt: runAt, Args: payload}
	return worker.Client.Enqueue(ctx, job)
}
