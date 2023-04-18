package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type MonitorIntegrationService struct {
	monitorIntegrationDomain *domain.MonitorIntegrationDomain
}

func NewMonitorIntegrationService(
	monitorIntegrationDomain *domain.MonitorIntegrationDomain,
) *MonitorIntegrationService {
	return &MonitorIntegrationService{monitorIntegrationDomain: monitorIntegrationDomain}
}

func (m *MonitorIntegrationService) Create(
	ctx context.Context, tx *sql.Tx, organization *model.Organization,
	monitorIntegrationType resource.MonitorIntegrationType, config map[string]string,
) (*model.MonitorIntegration, error) {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	monitorIntegration := &model.MonitorIntegration{
		Config:         string(jsonConfig),
		OrganizationID: organization.ID,
	}
	if monitorIntegrationType == resource.MonitorIntegrationTypeWebhook {
		svixApp, err := infra.CreateOrganizationApplication(ctx, organization.Slug)
		if err != nil {
			sentry.CaptureException(err)
			return nil, err
		}
		monitorIntegration.ExternalID = &svixApp.Id
		monitorIntegration, err = m.monitorIntegrationDomain.Create(ctx, tx, monitorIntegration, monitorIntegrationType)
		if err != nil {
			return nil, err
		}
	}
	return monitorIntegration, err
}
