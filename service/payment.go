package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/stripe/stripe-go/v74"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type PaymentService struct {
	paymentDomain *domain.PaymentDomain
}

func NewPaymentService(paymentDomain *domain.PaymentDomain) *PaymentService {
	return &PaymentService{paymentDomain: paymentDomain}
}

func (p *PaymentService) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	organizationID int64,
	planWithProduct pkg.PlanWithProduct,
) (*model.Subscription, error) {
	now := times.Now()

	isTrial := true
	status := resource.SubscriptionStatusTrialing
	expiry := now.Add(time.Hour * 24 * constant.TrialSubscriptionDurationInDays)

	if planWithProduct.Tier == int32(resource.ProductTierFree) {
		isTrial = false
		status = resource.SubscriptionStatusActive
		expiry = now.Add(time.Hour * 24 * constant.FreeSubscriptionDurationInDays)
	}

	return p.paymentDomain.CreateSubscription(ctx, tx, isTrial, status, expiry,
		planWithProduct.Plan.ID, planWithProduct.Product.ID, organizationID)
}

func (p *PaymentService) HandleStripeEvent(ctx context.Context, event stripe.Event) {
	switch event.Type {
	case constant.StripeInvoiceCreated, constant.StripeInvoicePaid, constant.StripeInvoicePaymentFailed:
		var stripeInvoice stripe.Invoice
		if err := json.Unmarshal(event.Data.Raw, &stripeInvoice); err != nil {
			panic(err)
		}
	case constant.StripeCustomerSubscriptionCreated,
		constant.StripeCustomerSubscriptionUpdated,
		constant.StripeCustomerSubscriptionDeleted:
		var stripeSubscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &stripeSubscription); err != nil {
			panic(err)
		}
	}
}
