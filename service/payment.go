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

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
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
			return p.createOrUpdateReceipt(ctx, tx, event, stripeInvoice)
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
) error {
	user, err := p.userDomain.GetUserFromPaymentCustomerID(ctx, invoice.Customer.ID)
	if err != nil {
		return errors.Newf("failed to get stripe customer: %s, err: %w", invoice.Customer.ID, err)
	}
	line := invoice.Lines.Data[0]
	plan, err := p.paymentDomain.GetPlanWithProductFromExternalPlanID(ctx, line.Price.ID)
	if err != nil {
		return errors.Newf("failed to get plan, external plan ID: %s, err: %w", line.Price.ID, err)
	}
	subscription, err := p.getSubscriptionID(ctx, invoice)
	if err != nil {
		return errors.Newf("failed to get subscription, external ID: %s, err: %w", invoice.Subscription.ID, err)
	}

	var subscriptionID int64
	if subscription != nil {
		subscriptionID = subscription.ID
	}
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
		PlanID:             &plan.Plan.ID,
		ProductID:          &plan.Product.ID,
		SubscriptionID:     &subscriptionID,
		OrganizationID:     *user.OrganizationID,
	}

	eventAt := time.Unix(invoice.PeriodStart, 0)
	lastEventAt := cache.GetReceiptEventForCustomer(ctx, invoice.Customer.ID)
	if lastEventAt == nil || times.CompareDate(eventAt, *lastEventAt) == constant.Date1AfterDate2 {
		_, err = p.paymentDomain.CreateReceipt(ctx, tx, receipt)
	}

	return err
}

func (p *PaymentService) getPaidAt(event stripe.Event) *time.Time {
	if event.Type == constant.StripeInvoicePaid {
		return lo.ToPtr(time.Unix(event.Created, 0))
	}
	return nil
}

func (p *PaymentService) getSubscriptionID(ctx context.Context, invoice stripe.Invoice) (*model.Subscription, error) {
	if invoice.Subscription == nil {
		return nil, nil
	}
	return p.paymentDomain.GetSubscriptionFromExternalID(ctx, invoice.Subscription.ID)
}
