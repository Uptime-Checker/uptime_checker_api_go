package infra

import (
	"context"

	svix "github.com/svix/svix-webhooks/go"
)

var svixClient *svix.Svix

func SetupSvix() {
	svixClient = svix.New("AUTH_TOKEN", nil)
}

func CreateOrganizationApplication(ctx context.Context, organizationSlug string) (*svix.ApplicationOut, error) {
	return svixClient.Application.Create(ctx, &svix.ApplicationIn{
		Name: organizationSlug,
	})
}
