package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type RegionDomain struct{}

func NewRegionDomain() *RegionDomain {
	return &RegionDomain{}
}

func (m *RegionDomain) Get(ctx context.Context, key string) (*model.Region, error) {
	stmt := SELECT(Region.AllColumns).FROM(Region).WHERE(Region.Key.EQ(String(key))).LIMIT(1)

	region := &model.Region{}
	err := stmt.QueryContext(ctx, infra.DB, region)
	return region, err
}

func (p *RegionDomain) List(ctx context.Context) ([]model.Region, error) {
	stmt := SELECT(Region.AllColumns).FROM(Region)

	var regions []model.Region
	err := stmt.QueryContext(ctx, infra.DB, &regions)
	return regions, err
}
