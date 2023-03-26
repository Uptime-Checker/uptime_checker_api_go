package domain

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type CheckDomain struct{}

func NewCheckDomain() *CheckDomain {
	return &CheckDomain{}
}

func (c *CheckDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
) (*model.Check, error) {
	insertStmt := Check.INSERT(Check.MutableColumns.Except(Check.InsertedAt, Check.UpdatedAt)).
		MODEL(check).
		RETURNING(Check.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, check)
	return check, err
}
