package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/service"
)

type MonitorController struct {
	monitorDomain  *domain.MonitorDomain
	monitorService *service.MonitorService
}

func NewMonitorController(
	monitorDomain *domain.MonitorDomain,
	monitorService *service.MonitorService,
) *MonitorController {
	return &MonitorController{monitorDomain: monitorDomain, monitorService: monitorService}
}

func (m *MonitorController) Create(c *fiber.Ctx) error {
	return nil
}

func (m *MonitorController) ListMonitors(c *fiber.Ctx) error {
	return nil
}
