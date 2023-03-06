package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

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

func (p *PaymentDomain) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	isTrial bool,
	status resource.SubscriptionStatus,
	expiresAt time.Time,
	planID, productID, organizationID int64,
) (*model.Subscription, error) {

	if !status.Valid() {
		return nil, constant.ErrInvalidSubscriptionStatus
	}
	statusValue := int32(status)

	now := times.Now()
	if times.CompareDate(now, expiresAt) == constant.Date1AfterDate2 {
		return nil, constant.ErrExpiresAtInThePast
	}

	subscription := &model.Subscription{
		Status:         &statusValue,
		StartsAt:       &now,
		ExpiresAt:      &expiresAt,
		IsTrial:        &isTrial,
		PlanID:         &planID,
		ProductID:      &productID,
		OrganizationID: &organizationID,
	}
	insertStmt := Subscription.INSERT(Subscription.MutableColumns).MODEL(subscription).
		RETURNING(Subscription.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, subscription)
	return subscription, err
}
