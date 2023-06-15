package service

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
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
	monitorIntegrationType resource.MonitorIntegrationType, config map[string]any,
) (*model.MonitorIntegration, error) {
	jsonConfig, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	prevIntegration, err := m.findIntegrationFromType(ctx, organization.ID, monitorIntegrationType)
	if err != nil {
		return nil, err
	}
	if prevIntegration != nil {
		now := times.Now()
		prevIntegration.Config = string(jsonConfig)
		prevIntegration.UpdatedAt = now
		_, err := m.monitorIntegrationDomain.Update(ctx, tx, prevIntegration.ID, prevIntegration)
		if err != nil {
			return nil, err
		}

		return prevIntegration, nil
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
	}
	monitorIntegration, err = m.monitorIntegrationDomain.Create(ctx, tx, monitorIntegration, monitorIntegrationType)
	if err != nil {
		return nil, err
	}
	return monitorIntegration, err
}

func (m *MonitorIntegrationService) findIntegrationFromType(
	ctx context.Context,
	orgID int64,
	monitorIntegrationType resource.MonitorIntegrationType,
) (*model.MonitorIntegration, error) {
	integrations, err := m.monitorIntegrationDomain.List(ctx, orgID)
	if err != nil {
		return nil, err
	}
	for i, integration := range integrations {
		if resource.MonitorIntegrationType(integration.Type) == monitorIntegrationType {
			return &integrations[i], nil
		}
	}
	return nil, constant.ErrIntegrationNotFound
}
