package channel

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

func SendAlarmWebhook(
	ctx context.Context,
	errorLog *model.ErrorLog,
	monitor *model.Monitor,
	alarm *model.Alarm,
	integration *model.MonitorIntegration,
	notification *model.MonitorNotification,
) {
	eventType := resource.MonitorNotificationType(notification.Type).String()
	data := map[string]any{
		"eventType": eventType,
		"name":      monitor.Name,
		"url":       monitor.URL,
	}
	if resource.MonitorNotificationType(notification.Type) == resource.MonitorNotificationTypeMonitorUp {
		data["downtime"] = alarm.ResolvedAt.Sub(alarm.InsertedAt).String()
	} else {
		data["error"] = errorLog.Text
	}

	webhookData := pkg.WebhookData{
		EventType: eventType,
		EventID:   *notification.ExternalID,
		EventAt:   notification.InsertedAt,
		Data:      data,
	}
	infra.SendWebhook(ctx, *integration.ExternalID, webhookData)
}
