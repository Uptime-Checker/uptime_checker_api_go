package domain

import (
	"context"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type PaymentDomain struct{}

func NewPaymentDomain() *PaymentDomain {
	return &PaymentDomain{}
}

type PlanWithProduct struct {
	*model.Plan
	*model.Product
}

func (p *PaymentDomain) GetPlanWithProduct(ctx context.Context, id int64) (*PlanWithProduct, error) {
	stmt := SELECT(Plan.AllColumns, Product.AllColumns).
		FROM(Plan.LEFT_JOIN(Product, Plan.ProductID.EQ(Product.ID))).WHERE(Plan.ID.EQ(Int(id)))

	planWithProduct := &PlanWithProduct{}
	err := stmt.QueryContext(ctx, infra.DB, planWithProduct)
	return planWithProduct, err
}
