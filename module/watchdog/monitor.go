package watchdog

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func (w *WatchDog) verify(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	monitor *model.Monitor,
	monitorRegion *model.MonitorRegion,
) error {
	tracingID := pkg.GetTracingID(ctx)
	alarmPolicy, err := w.alarmPolicyService.Get(ctx, monitor.ID, monitor.OrganizationID)
	if err != nil {
		return fmt.Errorf("could not find alarm policy %d, err: %w", monitor.ID, err)
	}
	monitorRegion, err = w.monitorRegionDomain.UpdateDown(ctx, tx, monitorRegion.ID, !check.Success)
	if err != nil {
		return fmt.Errorf("could not update monitor region %d, err: %w", monitorRegion.ID, err)
	}

	status := w.handleAlarmPolicy(monitor, alarmPolicy, check.Success)
	lgr.Print(tracingID, 1, "monitor status", status.String())
	return nil
}

func (w *WatchDog) handleAlarmPolicy(
	monitor *model.Monitor,
	alarmPolicy *model.AlarmPolicy,
	success bool,
) resource.MonitorStatus {
	status := resource.MonitorStatusPassing
	if !success {
		status = resource.MonitorStatusDegraded
		consecutiveCount := pkg.Abs(w.getMonitorConsecutiveCount(monitor, success))
		if alarmPolicy.Reason == string(resource.AlarmPolicyErrorThreshold) &&
			alarmPolicy.Threshold == consecutiveCount {
			status = resource.MonitorStatusFailing
		}
	}
	return status
}

func (w *WatchDog) getMonitorConsecutiveCount(monitor *model.Monitor, success bool) int32 {
	current := monitor.ConsecutiveCount
	if success {
		if current < 0 {
			current = 1
		} else {
			current++
		}
	}
	if current > 0 {
		current = -1
	} else {
		current--
	}
	return current
}
