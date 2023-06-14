package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type PropertyDomain struct{}

func NewPropertyDomain() *PropertyDomain {
	return &PropertyDomain{}
}

func (p *PropertyDomain) List(ctx context.Context) ([]model.Property, error) {
	stmt := SELECT(Property.AllColumns).FROM(Property)
	var properties []model.Property
	err := stmt.QueryContext(ctx, infra.DB, &properties)
	return properties, err
}
