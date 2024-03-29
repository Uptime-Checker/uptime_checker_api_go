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

type AssertionDomain struct{}

func NewAssertionDomain() *AssertionDomain {
	return &AssertionDomain{}
}

func (a *AssertionDomain) Create(
	ctx context.Context,
	tx *sql.Tx,
	assertion *model.Assertion,
	assertionSource resource.AssertionSource,
	assertionComparison resource.AssertionComparison,
) (*model.Assertion, error) {
	if !assertionSource.Valid() {
		return nil, constant.ErrInvalidAssertionSource
	}
	if !assertionComparison.Valid() {
		return nil, constant.ErrInvalidAssertionComparison
	}

	assertion.Source = int32(assertionSource)
	assertion.Comparison = int32(assertionComparison)

	insertStmt := Assertion.INSERT(Assertion.MutableColumns.Except(Assertion.InsertedAt, Assertion.UpdatedAt)).
		MODEL(assertion).RETURNING(Assertion.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, assertion)
	return assertion, err
}

func (a *AssertionDomain) ListAssertions(ctx context.Context, monitorID int64) ([]model.Assertion, error) {
	stmt := SELECT(Assertion.AllColumns).FROM(Assertion).WHERE(Assertion.MonitorID.EQ(Int(monitorID)))

	var assertions []model.Assertion
	err := stmt.QueryContext(ctx, infra.DB, &assertions)
	return assertions, err
}
