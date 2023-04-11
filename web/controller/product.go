package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"
)

type ProductController struct {
}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (u *ProductController) ListExternal(c *fiber.Ctx) error {
	return resp.ServeData(c, fiber.StatusOK, infra.ListProductsWithPrices())
}

func (u *ProductController) ListInternal(c *fiber.Ctx) error {
	infra.ListProductsWithPrices()
	return resp.ServeNoContent(c, fiber.StatusNoContent)
}
