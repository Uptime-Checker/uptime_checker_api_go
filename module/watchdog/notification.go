package watchdog

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/sourcegraph/conc/iter"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func (w *WatchDog) notify(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	errorLog *model.ErrorLog,
	monitor *model.Monitor,
	alarm *model.Alarm,
) error {
	alarmChannels, err := w.getAlarmChannels(ctx, monitor)
	if err != nil {
		return err
	}
	iterator := iter.Iterator[model.AlarmChannel]{
		MaxGoroutines: len(alarmChannels),
	}
	iterator.ForEach(alarmChannels, func(alarmChannel *model.AlarmChannel) {
		w.notifyAlarmChannel(ctx, tx, check, errorLog, monitor, alarm, alarmChannel)
	})
	return nil
}

func (w *WatchDog) getAlarmChannels(ctx context.Context, monitor *model.Monitor) ([]model.AlarmChannel, error) {
	alarmChannels, err := w.alarmChannelDomain.ListByMonitor(ctx, monitor.ID)
	if err != nil {
		return nil, errors.Newf("failed to list channels by monitor: %d, err: %w", monitor.ID, err)
	}
	if len(alarmChannels) == 0 {
		alarmChannels, err = w.alarmChannelDomain.ListByOrganization(ctx, monitor.OrganizationID)
		if err != nil {
			return nil, errors.Newf("failed to list channels by org: %d, err: %w", monitor.OrganizationID, err)
		}
	}
	if len(alarmChannels) == 0 {
		return nil, errors.Newf("no alarm channels found for monitor: %d", monitor.ID)
	}
	return alarmChannels, nil
}

func (w *WatchDog) notifyAlarmChannel(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	errorLog *model.ErrorLog,
	monitor *model.Monitor,
	alarm *model.Alarm,
	alarmChannel *model.AlarmChannel,
) {
}
