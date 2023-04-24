package infra

import (
	"context"

	svix "github.com/svix/svix-webhooks/go"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
)

var svixClient *svix.Svix

func SetupSvix() {
	svixClient = svix.New(config.App.SvixKey, nil)
}

func CreateOrganizationApplication(ctx context.Context, organizationSlug string) (*svix.ApplicationOut, error) {
	return svixClient.Application.Create(ctx, &svix.ApplicationIn{
		Name: organizationSlug,
	})
}

func SendWebhook(ctx context.Context, appID string, webhookData pkg.WebhookData) {
	tracingID := pkg.GetTracingID(ctx)
	outMessage, err := svixClient.Message.Create(ctx, appID, &svix.MessageIn{
		EventType: webhookData.EventType,
		EventId:   *svix.NullableString(webhookData.EventID),
		Payload:   webhookData.Data,
	})
	if err != nil {
		lgr.Error(tracingID, "failed to send webhook", err)
	} else {
		lgr.Print(tracingID, "webhook sent", outMessage.Id, "at", outMessage.Timestamp.String())
	}
}
