package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type AssertionService struct {
	assertionDomain *domain.AssertionDomain
}

func NewAssertionService(assertionDomain *domain.AssertionDomain) *AssertionService {
	return &AssertionService{assertionDomain: assertionDomain}
}

func (a *AssertionService) Create(
	ctx context.Context,
	tx *sql.Tx, monitorID int64, source int32, property *string, comparison int32, value string,
) (*model.Assertion, error) {
	assertion := model.Assertion{
		Property:  property,
		Value:     &value,
		MonitorID: monitorID,
	}

	return a.assertionDomain.Create(
		ctx,
		tx,
		&assertion,
		resource.AssertionSource(source),
		resource.AssertionComparison(comparison),
	)
}
