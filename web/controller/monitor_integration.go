package controller

import "github.com/Uptime-Checker/uptime_checker_api_go/service"

type MonitorIntegrationController struct {
	paymentService *service.PaymentService
}

func NewMonitorIntegrationController(
	paymentService *service.PaymentService,
) *MonitorIntegrationController {
	return &MonitorIntegrationController{paymentService: paymentService}
}
