package watchdog

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samber/lo"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
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

	status, err := w.handleAlarmPolicy(ctx, monitor, alarmPolicy, check.Success)
	if err != nil {
		return err
	}
	lgr.Print(tracingID, 1, "monitor status", status.String())
	return nil
}

func (w *WatchDog) handleAlarmPolicy(
	ctx context.Context,
	monitor *model.Monitor,
	alarmPolicy *model.AlarmPolicy,
	success bool,
) (*resource.MonitorStatus, error) {
	now := times.Now()
	status := resource.MonitorStatusPassing
	if !success {
		status = resource.MonitorStatusDegraded
		consecutiveCount := pkg.Abs(w.getMonitorConsecutiveCount(monitor, success))
		reason := resource.AlarmPolicyName(alarmPolicy.Reason)

		switch reason {
		case resource.AlarmPolicyErrorThreshold:
			if consecutiveCount >= alarmPolicy.Threshold {
				status = resource.MonitorStatusFailing
			}
		case resource.AlarmPolicyDurationThreshold:
			monitorStatus, err := w.monitorStatusDomain.GetLatest(ctx, monitor.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get monitor status, err: %w", err)
			}
			currentStatus := resource.MonitorStatus(monitorStatus.Status)
			if currentStatus == resource.MonitorStatusFailing || currentStatus == resource.MonitorStatusDegraded {
				if now.Sub(monitorStatus.InsertedAt).Seconds() > float64(alarmPolicy.Threshold) {
					status = resource.MonitorStatusFailing
				}
			}
		case resource.AlarmPolicyRegionThreshold:
			monitorRegions, err := w.monitorRegionDomain.GetAll(ctx, monitor.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get monitor regions, err: %w", err)
			}
			downRegions := lo.Filter(monitorRegions, func(monitorRegion model.MonitorRegion, index int) bool {
				return monitorRegion.Down
			})
			if int32(len(downRegions)) >= alarmPolicy.Threshold {
				status = resource.MonitorStatusFailing
			}
		}
	}
	return &status, nil
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
