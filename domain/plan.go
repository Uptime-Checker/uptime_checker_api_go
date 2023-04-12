package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type PlanDomain struct{}

func NewPlanDomain() *PlanDomain {
	return &PlanDomain{}
}

func (p *PlanDomain) Create(
	ctx context.Context, tx *sql.Tx, plan *model.Plan,
	price float64, planType resource.PlanType,
) (*model.Plan, error) {
	if !planType.Valid() {
		return nil, constant.ErrInvalidPlanType
	}
	plan.Price = price
	insertStmt := Plan.INSERT(Plan.MutableColumns.Except(Plan.InsertedAt, Plan.UpdatedAt)).
		MODEL(plan).ON_CONFLICT(Plan.ExternalID).DO_UPDATE(SET(Plan.Price.SET(Float(price)))).RETURNING(Plan.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, plan)
	return plan, err
}
