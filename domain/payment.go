package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type PaymentDomain struct{}

func NewPaymentDomain() *PaymentDomain {
	return &PaymentDomain{}
}

func (p *PaymentDomain) GetPlanWithProduct(ctx context.Context, id int64) (*pkg.PlanWithProduct, error) {
	stmt := SELECT(Plan.AllColumns, Product.AllColumns).
		FROM(Plan.LEFT_JOIN(Product, Plan.ProductID.EQ(Product.ID))).WHERE(Plan.ID.EQ(Int(id)))

	planWithProduct := &pkg.PlanWithProduct{}
	err := stmt.QueryContext(ctx, infra.DB, planWithProduct)
	return planWithProduct, err
}

func (p *PaymentDomain) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	isTrial bool,
	status string,
	expiresAt time.Time,
	planID, productID, organizationID int64,
) (*model.Subscription, error) {
	now := times.Now()
	if times.CompareDate(now, expiresAt) == constant.Date1AfterDate2 {
		return nil, constant.ErrExpiresAtInThePast
	}

	subscription := &model.Subscription{
		Status:         status,
		StartsAt:       &now,
		ExpiresAt:      &expiresAt,
		IsTrial:        isTrial,
		PlanID:         planID,
		ProductID:      productID,
		OrganizationID: organizationID,
	}
	insertStmt := Subscription.INSERT(Subscription.MutableColumns.
		Except(Subscription.InsertedAt, Subscription.UpdatedAt)).
		MODEL(subscription).
		RETURNING(Subscription.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, subscription)
	return subscription, err
}

func (p *PaymentDomain) CreateReceipt(
	ctx context.Context,
	tx *sql.Tx,
	receipt *model.Receipt,
) (*model.Receipt, error) {
	return nil, nil
}
