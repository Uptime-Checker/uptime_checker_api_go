package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vgarvardt/gue/v5"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/worker"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

type RunCheckTaskPayload struct {
	MonitorID, MonitorRegionID, RegionID int64
}

func RunCheckAsync(ctx context.Context, monitorID, monitorRegionID, regionID int64, runAt time.Time) error {
	body := RunCheckTaskPayload{
		MonitorID:       monitorID,
		MonitorRegionID: monitorRegionID,
		RegionID:        regionID,
	}

	tid := pkg.GetTracingID(ctx)
	lgr.Print(tid, "scheduling check run for monitor", monitorID, "region", regionID, "monitor region", monitorRegionID)

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	job := &gue.Job{Type: worker.TaskRunCheck, RunAt: runAt, Args: payload}
	return worker.Wheel.Enqueue(ctx, job)
}
