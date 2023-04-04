package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

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
	tx *sql.Tx,
	check *model.Check,
) (*model.Check, error) {
	insertStmt := Check.INSERT(Check.MutableColumns.Except(Check.InsertedAt, Check.UpdatedAt)).
		MODEL(check).
		RETURNING(Check.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, check)
	return check, err
}

func (c *CheckDomain) Update(
	ctx context.Context,
	tx *sql.Tx,
	check *model.Check,
) (*model.Check, error) {
	updateStmt := Check.UPDATE(
		Check.Success, Check.Body, Check.Traces, Check.Headers, Check.StatusCode, Check.ContentSize, Check.ContentType,
		Check.Duration, Check.UpdatedAt,
	).MODEL(check).
		WHERE(Check.ID.EQ(Int(check.ID))).
		RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, check)
	return check, err
}
