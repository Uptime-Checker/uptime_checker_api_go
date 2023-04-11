package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type ProductDomain struct{}

func NewProductDomain() *ProductDomain {
	return &ProductDomain{}
}

func (a *ProductDomain) ListProducts(ctx context.Context) ([]model.Product, error) {
	stmt := SELECT(Product.AllColumns).FROM(Product)

	var products []model.Product
	err := stmt.QueryContext(ctx, infra.DB, &products)
	return products, err
}
