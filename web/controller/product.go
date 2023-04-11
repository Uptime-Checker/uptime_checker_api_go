package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/middlelayer"
)

type ProductController struct {
	productDomain *domain.ProductDomain
	userDomain    *domain.UserDomain
}

func NewProductController(productDomain *domain.ProductDomain, userDomain *domain.UserDomain) *ProductController {
	return &ProductController{productDomain: productDomain, userDomain: userDomain}
}

func (p *ProductController) ListExternal(c *fiber.Ctx) error {
	return resp.ServeData(c, fiber.StatusOK, infra.ListProductsWithPrices())
}

func (p *ProductController) CreateBillingCustomer(c *fiber.Ctx) error {
	ctx := c.Context()
	user := middlelayer.GetUser(c)
	if user.PaymentCustomerID != nil {
		return resp.ServeNoContent(c, fiber.StatusNoContent)
	}
	billingCustomer, err := infra.CreateCustomer(user.Name, user.Email)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrBillingCustomerCreateFailed, err)
	}
	updatedUser, err := p.userDomain.UpdatePaymentID(ctx, user.ID, billingCustomer.ID)
	if err != nil {
		return resp.ServeError(c, fiber.StatusBadRequest, resp.ErrBillingCustomerUpdateFailed, err)
	}
	return resp.ServeData(c, fiber.StatusOK, updatedUser)
}

func (p *ProductController) ListInternal(c *fiber.Ctx) error {
	products, err := p.productDomain.ListProducts(c.Context())
	if err != nil {
		return resp.SendError(c, fiber.StatusInternalServerError, err)
	}
	return resp.ServeData(c, fiber.StatusOK, products)
}
