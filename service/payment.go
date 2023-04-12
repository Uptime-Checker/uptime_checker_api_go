package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type PaymentService struct {
	userDomain    *domain.UserDomain
	paymentDomain *domain.PaymentDomain
}

func NewPaymentService(userDomain *domain.UserDomain, paymentDomain *domain.PaymentDomain) *PaymentService {
	return &PaymentService{userDomain: userDomain, paymentDomain: paymentDomain}
}

func (p *PaymentService) CreateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	organizationID int64,
	planWithProduct pkg.PlanWithProduct,
) (*model.Subscription, error) {
	now := times.Now()

	isTrial := true
	status := stripe.SubscriptionStatusTrialing
	expiry := now.Add(time.Hour * 24 * constant.TrialSubscriptionDurationInDays)

	if planWithProduct.Tier == int32(resource.ProductTierFree) {
		isTrial = false
		status = stripe.SubscriptionStatusTrialing
		expiry = now.Add(time.Hour * 24 * constant.FreeSubscriptionDurationInDays)
	}

	return p.paymentDomain.CreateSubscription(ctx, tx, isTrial, string(status), expiry,
		planWithProduct.Plan.ID, planWithProduct.Product.ID, organizationID)
}

func (p *PaymentService) HandleStripeEvent(ctx context.Context, event stripe.Event) {
	tracingID := pkg.GetTracingID(ctx)
	if err := infra.Transaction(ctx, func(tx *sql.Tx) error {
		switch event.Type {
		case constant.StripeInvoiceCreated, constant.StripeInvoicePaid, constant.StripeInvoicePaymentFailed:
			var stripeInvoice stripe.Invoice
			if err := json.Unmarshal(event.Data.Raw, &stripeInvoice); err != nil {
				lgr.Error(tracingID, 1, "failed to unmarshal stripe invoice", err)
				return err
			}
			user, err := p.userDomain.GetUserFromPaymentCustomerID(ctx, stripeInvoice.Customer.ID)
			if err != nil {
				lgr.Error(tracingID, 2, "failed to get user, stripe customer:", stripeInvoice.Customer.ID, err)
				return errors.Newf("failed to get user, stripe customer: %s, err: %w", stripeInvoice.Customer.ID, err)
			}
			return p.createOrUpdateReceipt(ctx, tx, event, stripeInvoice, user)
		case constant.StripeCustomerSubscriptionCreated,
			constant.StripeCustomerSubscriptionUpdated,
			constant.StripeCustomerSubscriptionDeleted:
			var stripeSubscription stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &stripeSubscription); err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		sentry.CaptureException(errors.Newf("failed to commit stripe webhook transaction, err: %w", err))
	}
}

func (p *PaymentService) createOrUpdateReceipt(
	ctx context.Context,
	tx *sql.Tx,
	event stripe.Event,
	invoice stripe.Invoice,
	user *model.User,
) error {
	receipt := &model.Receipt{
		Price:              float64(invoice.Total),
		Currency:           string(invoice.Currency),
		ExternalID:         &invoice.ID,
		ExternalCustomerID: &invoice.Customer.ID,
		URL:                &invoice.HostedInvoiceURL,
		Status:             string(invoice.Status),
		Paid:               invoice.Paid,
		PaidAt:             p.getPaidAt(event),
		From:               lo.ToPtr(time.Unix(invoice.PeriodStart, 0)),
		To:                 lo.ToPtr(time.Unix(invoice.PeriodEnd, 0)),
		IsTrial:            false,
		PlanID:             new(int64),
		ProductID:          new(int64),
		SubscriptionID:     new(int64),
		OrganizationID:     *user.OrganizationID,
	}

	_, err := p.paymentDomain.CreateReceipt(ctx, tx, receipt)
	return err
}

func (p *PaymentService) getPaidAt(event stripe.Event) *time.Time {
	if event.Type == constant.StripeInvoicePaid {
		return lo.ToPtr(time.Unix(event.Created, 0))
	}
	return nil
}
