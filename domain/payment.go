package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"

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

func (p *PaymentDomain) GetPlanWithProductFromExternalPlanID(
	ctx context.Context,
	id string,
) (*pkg.PlanWithProduct, error) {
	stmt := SELECT(Plan.AllColumns, Product.AllColumns).
		FROM(Plan.LEFT_JOIN(Product, Plan.ProductID.EQ(Product.ID))).WHERE(Plan.ExternalID.EQ(String(id)))

	planWithProduct := &pkg.PlanWithProduct{}
	err := stmt.QueryContext(ctx, infra.DB, planWithProduct)
	return planWithProduct, err
}

func (u *PaymentDomain) GetSubscriptionFromExternalID(
	ctx context.Context,
	externalID string,
) (*model.Subscription, error) {
	stmt := SELECT(Subscription.AllColumns).FROM(Subscription).WHERE(
		Subscription.ExternalID.EQ(String(externalID)),
	).LIMIT(1)

	subscription := &model.Subscription{}
	err := stmt.QueryContext(ctx, infra.DB, subscription)
	return subscription, err
}

func (p *PaymentDomain) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	subscription *model.Subscription,
) (*model.Subscription, error) {
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
	insertStmt := Receipt.INSERT(Receipt.MutableColumns.Except(Receipt.InsertedAt, Receipt.UpdatedAt)).MODEL(receipt).
		ON_CONFLICT(Receipt.ExternalID).DO_UPDATE(SET(
		Receipt.Price.SET(Float(receipt.Price)),
		Receipt.Paid.SET(Bool(receipt.Paid)),
		Receipt.Status.SET(String(receipt.Status)),
		Receipt.URL.SET(infra.GetStringExpression(receipt.URL)),
		Receipt.PaidAt.SET(infra.GetTimestampExpression(receipt.PaidAt)),
		Receipt.IsTrial.SET(Bool(receipt.IsTrial)),
		Receipt.SubscriptionID.SET(infra.GetIntegerExpression(receipt.SubscriptionID)),
	)).RETURNING(Receipt.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, receipt)
	return receipt, err
}
