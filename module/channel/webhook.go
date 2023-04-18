package channel

import "time"

type WebhookData struct {
	eventType, eventID string
	eventAt            time.Time
	data               map[string]interface{}
}

func SendWebhook(webhookData WebhookData) {
}
