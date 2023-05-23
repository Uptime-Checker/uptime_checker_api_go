package controller

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/web/controller/resp"

	"github.com/gofiber/fiber/v2"
)

type RegionController struct {
	regionDomain *domain.RegionDomain
}

func NewRegionController(
	regionDomain *domain.RegionDomain,
) *RegionController {
	return &RegionController{regionDomain: regionDomain}
}

func (r *RegionController) List(c *fiber.Ctx) error {
	ctx := c.Context()
	regions, err := r.regionDomain.List(ctx)
	if err != nil {
		return resp.SendError(c, fiber.StatusInternalServerError, err)
	}
	return resp.ServeData(c, fiber.StatusOK, regions)
}
