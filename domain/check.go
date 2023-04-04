package domain

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type CheckDomain struct{}

func NewCheckDomain() *CheckDomain {
	return &CheckDomain{}
}

// Create creates a check, it does not use transaction
func (c *CheckDomain) Create(
	ctx context.Context,
	check *model.Check,
) (*model.Check, error) {
	insertStmt := Check.INSERT(Check.MutableColumns.Except(Check.InsertedAt, Check.UpdatedAt)).
		MODEL(check).
		RETURNING(Check.AllColumns)
	err := insertStmt.QueryContext(ctx, infra.DB, check)
	return check, err
}
