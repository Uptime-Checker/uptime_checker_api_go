package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/stripe/stripe-go/v74"

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

func (p *PaymentDomain) GetSubscriptionFromExternalID(
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

func (p *PaymentDomain) GetLocalActiveSubscription(
	ctx context.Context,
	organizationID int64,
) (*model.Subscription, error) {
	stmt := SELECT(Subscription.AllColumns).FROM(Subscription).WHERE(
		Subscription.OrganizationID.EQ(Int(organizationID)).
			AND(Subscription.ExternalID.IS_NULL()).
			AND(Subscription.Status.EQ(String(string(stripe.SubscriptionStatusActive)))),
	).LIMIT(1)

	subscription := &model.Subscription{}
	err := stmt.QueryContext(ctx, infra.DB, subscription)
	return subscription, err
}

func (p *PaymentDomain) ListActiveSubscriptions(
	ctx context.Context,
	organizationID int64,
) ([]model.Subscription, error) {
	stmt := SELECT(Subscription.AllColumns).FROM(Subscription).
		WHERE(
			Subscription.OrganizationID.EQ(Int(organizationID)).
				AND(Subscription.Status.EQ(String(string(stripe.SubscriptionStatusActive)))),
		)
	var subscriptions []model.Subscription
	err := stmt.QueryContext(ctx, infra.DB, &subscriptions)
	return subscriptions, err
}

func (p *PaymentDomain) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	subscription *model.Subscription,
) (*model.Subscription, error) {
	insertStmt := Subscription.INSERT(
		Subscription.MutableColumns.Except(Subscription.InsertedAt, Subscription.UpdatedAt),
	).MODEL(subscription).RETURNING(Subscription.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, subscription)
	return subscription, err
}

func (p *PaymentDomain) UpdateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	subscription *model.Subscription,
) (*model.Subscription, error) {
	now := times.Now()
	subscription.UpdatedAt = now

	updateStmt := Subscription.UPDATE(
		Subscription.Status,
		Subscription.IsTrial,
		Subscription.ExpiresAt,
		Subscription.CanceledAt,
		Subscription.CancellationReason,
		Subscription.PlanID,
		Subscription.ProductID,
		Subscription.UpdatedAt,
	).MODEL(subscription).WHERE(Subscription.ID.EQ(Int(id))).RETURNING(Subscription.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, subscription)
	return subscription, err
}

func (p *PaymentDomain) GetReceiptFromExternalID(
	ctx context.Context,
	externalID string,
) (*model.Receipt, error) {
	stmt := SELECT(Receipt.AllColumns).FROM(Receipt).WHERE(
		Receipt.ExternalID.EQ(String(externalID)),
	).LIMIT(1)

	receipt := &model.Receipt{}
	err := stmt.QueryContext(ctx, infra.DB, receipt)
	return receipt, err
}

func (p *PaymentDomain) CreateReceipt(
	ctx context.Context,
	tx *sql.Tx,
	receipt *model.Receipt,
) (*model.Receipt, error) {
	insertStmt := Receipt.INSERT(Receipt.MutableColumns.Except(Receipt.InsertedAt, Receipt.UpdatedAt)).
		MODEL(receipt).RETURNING(Receipt.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, receipt)
	return receipt, err
}

func (p *PaymentDomain) UpdateReceipt(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	receipt *model.Receipt,
) (*model.Receipt, error) {
	now := times.Now()
	receipt.UpdatedAt = now

	updateStmt := Receipt.UPDATE(
		Receipt.Status,
		Receipt.PaidAt,
		Receipt.Price,
		Receipt.Paid,
		Receipt.URL,
		Receipt.IsTrial,
		Receipt.SubscriptionID,
		Receipt.ExternalSubscriptionID,
		Receipt.UpdatedAt,
	).MODEL(receipt).WHERE(Receipt.ID.EQ(Int(id))).RETURNING(Receipt.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, receipt)
	return receipt, err
}

func (p *PaymentDomain) ExpireSubscription(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
) (*model.Subscription, error) {
	now := times.Now()
	subscription := &model.Subscription{
		Status:    string(stripe.SubscriptionStatusPaused),
		ExpiresAt: &now,
		UpdatedAt: now,
	}

	updateStmt := Subscription.UPDATE(Subscription.Status, Subscription.ExpiresAt, Subscription.UpdatedAt).
		MODEL(subscription).WHERE(Subscription.ID.EQ(Int(id))).RETURNING(Subscription.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, subscription)
	return subscription, err
}
