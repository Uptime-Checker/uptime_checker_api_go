package service

import "github.com/Uptime-Checker/uptime_checker_api_go/domain"

type MonitorService struct {
	monitorDomain *domain.MonitorDomain
}

func NewMonitorService(monitorDomain *domain.MonitorDomain) *MonitorService {
	return &MonitorService{monitorDomain: monitorDomain}
}
