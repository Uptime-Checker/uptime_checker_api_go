package watchdog

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func (w *WatchDog) alarmCheck(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
	check *model.Check,
	status resource.MonitorStatus,
) error {
	var ongoingAlarm *model.Alarm
	alarm, err := w.alarmDomain.GetOngoing(ctx, monitor.ID)
	if err == nil {
		ongoingAlarm = alarm
	}
	if status == resource.MonitorStatusPassing {
		return w.resolveAlarm(ctx, tx, monitor, check, ongoingAlarm)
	} else if status == resource.MonitorStatusFailing {
		return w.raiseAlarm(ctx, tx, monitor, check, ongoingAlarm)
	}
	return nil
}

func (w *WatchDog) resolveAlarm(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
	check *model.Check,
	alarm *model.Alarm,
) error {
	tracingID := pkg.GetTracingID(ctx)
	if alarm == nil {
		lgr.Print(tracingID, 1, "no ongoing alarm to resolve", monitor.ID)
		return nil
	}
	_, err := w.alarmDomain.Resolve(ctx, tx, alarm.ID, check.ID)
	if err != nil {
		return errors.Newf("failed to resolve alarm, err: %w", err)
	}
	lgr.Print(tracingID, 2, "resolved alarm", alarm.ID)
	// update daily report duration
	// send notification
	return nil
}

func (w *WatchDog) raiseAlarm(
	ctx context.Context,
	tx *sql.Tx,
	monitor *model.Monitor,
	check *model.Check,
	alarm *model.Alarm,
) error {
	tracingID := pkg.GetTracingID(ctx)
	if alarm != nil {
		lgr.Print(tracingID, "no new alarm to raise", monitor.ID, "ongoing alarm", alarm.ID)
		return nil
	}
	alarm = &model.Alarm{
		Ongoing:            true,
		TriggeredByCheckID: &check.ID,
		MonitorID:          monitor.ID,
		OrganizationID:     monitor.OrganizationID,
	}
	_, err := w.alarmDomain.Create(ctx, tx, alarm)
	if err != nil {
		return errors.Newf("failed to create alarm, err: %w", err)
	}
	lgr.Print(tracingID, 2, "raised alarm", alarm.ID)
	// send notification
	return nil
}
