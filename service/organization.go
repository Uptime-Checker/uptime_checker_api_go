package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type OrganizationService struct {
	organizationDomain *domain.OrganizationDomain
}

func NewOrganizationService(organizationDomain *domain.OrganizationDomain) *OrganizationService {
	return &OrganizationService{organizationDomain: organizationDomain}
}

func (o *OrganizationService) Create(
	ctx context.Context,
	tx *sql.Tx,
	name, slug string,
	userID int64,
) (*model.Organization, error) {
	organization, err := o.organizationDomain.CreateOrganization(ctx, tx, name, slug)
	if err != nil {
		return nil, err
	}

	role, err := o.organizationDomain.GetRoleByType(ctx, resource.RoleTypeSuperAdmin)
	if err != nil {
		return nil, err
	}
	_, err = o.organizationDomain.CreateOrganizationUser(ctx, tx, organization.ID, userID, role.ID)
	if err != nil {
		return nil, err
	}
	return organization, nil
}
