package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v74/webhook"

	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
)

type WebhookController struct {
	paymentService *service.PaymentService
}

func NewWebhookController(
	paymentService *service.PaymentService,
) *WebhookController {
	return &WebhookController{paymentService: paymentService}
}

func (w *WebhookController) StripePayment(c *fiber.Ctx) error {
	event, err := webhook.ConstructEvent(c.Body(), c.Get(constant.StripeSignatureHeader), config.App.StripeWebhookKey)
	if err != nil {
		return resp.ServeInternalServerError(c, err)
	}
	switch event.Type {
	case constant.StripeInvoiceCreated,
		constant.StripeInvoicePaid,
		constant.StripeInvoicePaymentFailed,
		constant.StripeCustomerSubscriptionCreated,
		constant.StripeCustomerSubscriptionUpdated,
		constant.StripeCustomerSubscriptionDeleted:
		w.paymentService.HandleStripeEvent(c.Context(), event)
	}

	return resp.ServeNoContent(c, fiber.StatusNoContent)
}
