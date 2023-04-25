package watchdog

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/sourcegraph/conc/iter"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/module/watchdog/channel"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
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
	if err := w.handleNotifyAlarmChannel(ctx, tx, check, errorLog, monitor, alarm, alarmChannel); err != nil {
		sentry.CaptureException(err)
	}
}

func (w *WatchDog) handleNotifyAlarmChannel(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
	errorLog *model.ErrorLog,
	monitor *model.Monitor,
	alarm *model.Alarm,
	alarmChannel *model.AlarmChannel,
) error {
	tracingID := pkg.GetTracingID(ctx)
	if !alarmChannel.On {
		lgr.Print(tracingID, "alarm channel is off", alarmChannel.ID)
		return nil
	}
	if alarmChannel.IntegrationID == nil && alarmChannel.UserContactID == nil {
		return errors.Newf("no integration or user contact found, alarm channel: %d", alarmChannel.ID)
	}
	notificationType := resource.MonitorNotificationTypeMonitorUp
	if check.Success {
		notificationType = resource.MonitorNotificationTypeMonitorDown
	}
	notification := &model.MonitorNotification{
		Successful:     false,
		AlarmID:        &alarm.ID,
		MonitorID:      monitor.ID,
		OrganizationID: monitor.OrganizationID,
	}
	if alarmChannel.IntegrationID != nil {
		integration, err := w.monitorIntegrationDomain.Get(ctx, *alarmChannel.IntegrationID)
		if err != nil {
			return errors.Newf("failed to get integration: %d, err: %w", *alarmChannel.IntegrationID, err)
		}
		integrationType := resource.MonitorIntegrationType(integration.Type)
		notification.IntegrationID = &integration.ID
		if integrationType == resource.MonitorIntegrationTypeWebhook {
			eventID := pkg.GetUniqueString()
			notification.ExternalID = &eventID
		}
		notification, err := w.monitorNotificationDomain.Create(ctx, tx, notification, notificationType)
		if err != nil {
			return errors.Newf("failed to create notification: %d, err: %w", alarmChannel.ID, err)
		}
		channel.SendAlarmWebhook(ctx, errorLog, monitor, alarm, integration, notification)
	}
	return nil
}
