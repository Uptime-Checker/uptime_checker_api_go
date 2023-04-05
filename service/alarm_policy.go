package service

import "github.com/Uptime-Checker/uptime_checker_api_go/domain"

type AlarmPolicyService struct {
	alarmPolicyDomain *domain.AlarmPolicyDomain
}

func NewAlarmPolicyService(alarmPolicyDomain *domain.AlarmPolicyDomain) *AlarmPolicyService {
	return &AlarmPolicyService{alarmPolicyDomain: alarmPolicyDomain}
}
