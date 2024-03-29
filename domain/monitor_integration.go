package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type MonitorIntegrationDomain struct{}

func NewMonitorIntegrationDomain() *MonitorIntegrationDomain {
	return &MonitorIntegrationDomain{}
}

func (m *MonitorIntegrationDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	monitorIntegration *model.MonitorIntegration,
	monitorIntegrationType resource.MonitorIntegrationType,
) (*model.MonitorIntegration, error) {
	if !monitorIntegrationType.Valid() {
		return nil, constant.ErrInvalidMonitorIntegrationType
	}
	monitorIntegration.Type = int32(monitorIntegrationType)
	insertStmt := MonitorIntegration.INSERT(MonitorIntegration.MutableColumns.
		Except(MonitorIntegration.InsertedAt, MonitorIntegration.UpdatedAt)).
		MODEL(monitorIntegration).
		RETURNING(MonitorIntegration.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, monitorIntegration)
	return monitorIntegration, err
}

func (m *MonitorIntegrationDomain) Get(ctx context.Context, id int64) (*model.MonitorIntegration, error) {
	stmt := SELECT(MonitorIntegration.AllColumns).FROM(MonitorIntegration).WHERE(MonitorIntegration.ID.EQ(Int(id))).
		LIMIT(1)

	integration := &model.MonitorIntegration{}
	err := stmt.QueryContext(ctx, infra.DB, integration)
	return integration, err
}

func (m *MonitorIntegrationDomain) List(
	ctx context.Context,
	organizationID int64,
) ([]model.MonitorIntegration, error) {
	stmt := SELECT(MonitorIntegration.AllColumns).FROM(MonitorIntegration).
		WHERE(MonitorIntegration.OrganizationID.EQ(Int(organizationID)))

	var monitorIntegrations []model.MonitorIntegration
	err := stmt.QueryContext(ctx, infra.DB, &monitorIntegrations)
	return monitorIntegrations, err
}

func (c *MonitorIntegrationDomain) Update(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	monitorIntegration *model.MonitorIntegration,
) (*model.MonitorIntegration, error) {
	updateStmt := MonitorIntegration.UPDATE(
		MonitorIntegration.Config, MonitorIntegration.UpdatedAt,
	).MODEL(monitorIntegration).WHERE(MonitorIntegration.ID.EQ(Int(id))).RETURNING(MonitorIntegration.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, monitorIntegration)
	return monitorIntegration, err
}
