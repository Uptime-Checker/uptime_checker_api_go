package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
		status = stripe.SubscriptionStatusActive
		expiry = now.Add(time.Hour * 24 * constant.FreeSubscriptionDurationInDays)
	}

	subscription := &model.Subscription{
		Status:         string(status),
		StartsAt:       &now,
		ExpiresAt:      &expiry,
		IsTrial:        isTrial,
		PlanID:         planWithProduct.Plan.ID,
		ProductID:      planWithProduct.Product.ID,
		OrganizationID: organizationID,
	}

	return p.paymentDomain.CreateSubscription(ctx, tx, subscription)
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
			return p.createOrUpdateSubscription(ctx, tx, event, stripeSubscription)
		}
		return nil
	}); err != nil {
		sentry.CaptureException(errors.Newf("failed to commit stripe webhook transaction, err: %w", err))
	}
}

func (p *PaymentService) createOrUpdateSubscription(
	ctx context.Context,
	tx *sql.Tx,
	event stripe.Event,
	remoteSubscription stripe.Subscription,
) error {
	user, err := p.userDomain.GetUserFromPaymentCustomerID(ctx, remoteSubscription.Customer.ID)
	if err != nil {
		return errors.Newf("failed to get stripe customer: %s, err: %w", remoteSubscription.Customer.ID, err)
	}
	item := remoteSubscription.Items.Data[0]
	planWithProduct, err := p.paymentDomain.GetPlanWithProductFromExternalPlanID(ctx, item.Price.ID)
	if err != nil {
		return errors.Newf("failed to get plan, external plan ID: %s, err: %w", item.Price.ID, err)
	}

	localSubscription, err := p.paymentDomain.GetLocalActiveSubscription(ctx, *user.OrganizationID)
	if err == nil { // local active subscription found
		_, err := p.paymentDomain.ExpireSubscription(ctx, tx, localSubscription.ID)
		if err != nil {
			return errors.Newf("failed to expire free subscription, org: %d, err: %w", *user.OrganizationID, err)
		}
	}

	subscription := &model.Subscription{
		Status:             string(remoteSubscription.Status),
		StartsAt:           lo.ToPtr(time.Unix(remoteSubscription.StartDate, 0)),
		ExpiresAt:          lo.ToPtr(time.Unix(remoteSubscription.CurrentPeriodEnd, 0)),
		CanceledAt:         p.getCanceledAt(remoteSubscription),
		CancellationReason: p.getCancellationReason(remoteSubscription),
		IsTrial:            remoteSubscription.Status == stripe.SubscriptionStatusTrialing,
		ExternalID:         &remoteSubscription.ID,
		ExternalCustomerID: user.PaymentCustomerID,
		PlanID:             planWithProduct.Plan.ID,
		ProductID:          planWithProduct.Product.ID,
		OrganizationID:     *user.OrganizationID,
	}

	eventAt := time.Unix(event.Created, 0)
	lastEventAt := cache.GetPaymentEventForCustomer(ctx, cache.GetSubscriptionEventKey(remoteSubscription.Customer.ID))
	if lastEventAt == nil || times.CompareDate(eventAt, *lastEventAt) == constant.Date1AfterDate2 {
		_, err = p.paymentDomain.CreateSubscription(ctx, tx, subscription)
		if err != nil {
			return errors.Newf("failed to create subscription, err: %w", err)
		}
		cache.SetPaymentEventForCustomer(ctx, cache.GetSubscriptionEventKey(remoteSubscription.Customer.ID), eventAt)
	}

	activeSubscriptions, err := p.paymentDomain.ListActiveSubscriptions(ctx, *user.OrganizationID)
	if err != nil {
		return errors.Newf("failed to list active subscriptions, err: %w", err)
	}

	if len(activeSubscriptions) > 1 {
		sentry.CaptureMessage(fmt.Sprintf("more than one active subscription, org: %d", *user.OrganizationID))
	}
	return nil
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
	planWithProduct, err := p.paymentDomain.GetPlanWithProductFromExternalPlanID(ctx, line.Price.ID)
	if err != nil {
		return errors.Newf("failed to get plan, external plan ID: %s, err: %w", line.Price.ID, err)
	}

	var subscriptionID int64
	if invoice.Subscription != nil {
		subscription, err := p.paymentDomain.GetSubscriptionFromExternalID(ctx, invoice.Subscription.ID)
		if err != nil {
			return errors.Newf("failed to get subscription, external ID: %s, err: %w", invoice.Subscription.ID, err)
		}
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
		PlanID:             &planWithProduct.Plan.ID,
		ProductID:          &planWithProduct.Product.ID,
		SubscriptionID:     &subscriptionID,
		OrganizationID:     *user.OrganizationID,
	}

	eventAt := time.Unix(event.Created, 0)
	lastEventAt := cache.GetPaymentEventForCustomer(ctx, cache.GetReceiptEventKey(invoice.Customer.ID))
	if lastEventAt == nil || times.CompareDate(eventAt, *lastEventAt) == constant.Date1AfterDate2 {
		_, err = p.paymentDomain.CreateReceipt(ctx, tx, receipt)
		if err != nil {
			return errors.Newf("failed to create receipt, err: %w", err)
		}
		cache.SetPaymentEventForCustomer(ctx, cache.GetReceiptEventKey(invoice.Customer.ID), eventAt)
	}

	return nil
}

func (p *PaymentService) getPaidAt(event stripe.Event) *time.Time {
	if event.Type == constant.StripeInvoicePaid {
		return lo.ToPtr(time.Unix(event.Created, 0))
	}
	return nil
}

func (p *PaymentService) getCanceledAt(subscription stripe.Subscription) *time.Time {
	if subscription.CancelAt != 0 {
		return lo.ToPtr(time.Unix(subscription.CancelAt, 0))
	}
	return nil
}

func (p *PaymentService) getCancellationReason(subscription stripe.Subscription) *string {
	var reason string
	cancellationDetails := subscription.CancellationDetails
	if cancellationDetails != nil {
		reason = fmt.Sprintf("%s.%s", string(cancellationDetails.Reason), string(cancellationDetails.Feedback))
		if cancellationDetails.Comment != "" {
			reason = fmt.Sprintf("%s.%s", reason, cancellationDetails.Comment)
		}
	}
	return &reason
}
