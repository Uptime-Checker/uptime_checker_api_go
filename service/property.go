package service

import (
	"context"

	"github.com/getsentry/sentry-go"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
)

type PropertyService struct {
	propertyDomain *domain.PropertyDomain
}

func NewPropertyService(propertyDomain *domain.PropertyDomain) *PropertyService {
	return &PropertyService{propertyDomain: propertyDomain}
}

func (p *PropertyService) Get(ctx context.Context, propertyKey resource.PropertyKey) *string {
	properties, err := p.propertyDomain.List(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	for _, property := range properties {
		if property.Key == string(propertyKey) {
			return &property.Value
		}
	}
	return nil
}
