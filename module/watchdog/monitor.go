package watchdog

import (
	"context"
	"database/sql"
	"time"

	"github.com/cockroachdb/errors"
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
	dailyReport *model.DailyReport,
) error {
	tracingID := pkg.GetTracingID(ctx)
	alarmPolicy, err := w.alarmPolicyService.Get(ctx, monitor.ID, monitor.OrganizationID)
	if err != nil {
		return errors.Newf("could not find alarm policy %d, err: %w", monitor.ID, err)
	}
	monitorRegion, err = w.monitorRegionDomain.UpdateDown(ctx, tx, monitorRegion.ID, !check.Success)
	if err != nil {
		return errors.Newf("could not update monitor region %d, err: %w", monitorRegion.ID, err)
	}

	monitorStatus, err := w.monitorStatusDomain.GetLatest(ctx, monitor.ID)
	if err != nil {
		return errors.Newf("failed to get monitor status, err: %w", err)
	}
	status, err := w.handleAlarmPolicy(ctx, monitor, monitorStatus, alarmPolicy, check.Success)
	if err != nil {
		return err
	}
	lgr.Print(tracingID, 1, "monitor status", status.String())

	if resource.MonitorStatus(monitorStatus.Status) != status {
		_, err := w.monitorStatusDomain.Create(ctx, tx, monitorStatus, status)
		if err != nil {
			return errors.Newf("failed to create monitor status, err: %w", err)
		}
	}

	// Update monitor params - status | consecutive
	var lastFailedAt *time.Time
	if !check.Success {
		lastFailedAt = &check.InsertedAt
	}
	consecutiveCount := w.getMonitorConsecutiveCount(monitor, check.Success)
	monitor, err = w.monitorDomain.UpdateConsecutive(ctx, tx, monitor.ID, status, consecutiveCount, lastFailedAt)
	if err != nil {
		return errors.Newf("failed to update monitor consecutive count, err: %w", err)
	}

	return w.alarmCheck(ctx, tx, monitor, check, status, dailyReport)
}

func (w *WatchDog) handleAlarmPolicy(
	ctx context.Context,
	monitor *model.Monitor,
	monitorStatus *model.MonitorStatusChange,
	alarmPolicy *model.AlarmPolicy,
	success bool,
) (resource.MonitorStatus, error) {
	now := times.Now()
	status := resource.MonitorStatus(monitor.Status)
	reason := resource.AlarmPolicyName(alarmPolicy.Reason)
	consecutiveCount := pkg.Abs(w.getMonitorConsecutiveCount(monitor, success))

	if !success {
		status = resource.MonitorStatusDegraded

		switch reason {
		case resource.AlarmPolicyErrorThreshold:
			if consecutiveCount >= alarmPolicy.Threshold {
				status = resource.MonitorStatusFailing
			}
		case resource.AlarmPolicyDurationThreshold:
			currentStatus := resource.MonitorStatus(monitorStatus.Status)
			if currentStatus == resource.MonitorStatusFailing || currentStatus == resource.MonitorStatusDegraded {
				if now.Sub(monitorStatus.InsertedAt).Seconds() > float64(alarmPolicy.Threshold) {
					status = resource.MonitorStatusFailing
				}
			}
		case resource.AlarmPolicyRegionThreshold:
			monitorRegions, err := w.monitorRegionDomain.GetAll(ctx, monitor.ID)
			if err != nil {
				return status, errors.Newf("failed to get monitor regions, err: %w", err)
			}
			downRegions := lo.Filter(monitorRegions, func(monitorRegion model.MonitorRegion, index int) bool {
				return monitorRegion.Down
			})
			if int32(len(downRegions)) >= alarmPolicy.Threshold {
				status = resource.MonitorStatusFailing
			}
		}
	} else {
		if status == resource.MonitorStatusFailing {
			status = resource.MonitorStatusDegraded
		}
		switch reason {
		case resource.AlarmPolicyErrorThreshold:
			if consecutiveCount >= alarmPolicy.Threshold {
				status = resource.MonitorStatusPassing
			}
		case resource.AlarmPolicyDurationThreshold:
			if monitor.LastFailedAt != nil &&
				now.Sub(*monitor.LastFailedAt).Seconds() > float64(alarmPolicy.Threshold) {
				status = resource.MonitorStatusPassing
			}
		case resource.AlarmPolicyRegionThreshold:
			monitorRegions, err := w.monitorRegionDomain.GetAll(ctx, monitor.ID)
			if err != nil {
				return status, errors.Newf("failed to get monitor regions, err: %w", err)
			}
			upRegions := lo.Filter(monitorRegions, func(monitorRegion model.MonitorRegion, index int) bool {
				return !monitorRegion.Down
			})
			if int32(len(upRegions)) >= alarmPolicy.Threshold {
				status = resource.MonitorStatusPassing
			}
		}
	}
	return status, nil
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
