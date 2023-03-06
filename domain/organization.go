package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type OrganizationDomain struct{}

func NewOrganizationDomain() *OrganizationDomain {
	return &OrganizationDomain{}
}

func (o *OrganizationDomain) CreateOrganization(
	ctx context.Context,
	tx *sql.Tx,
	name, slug string,
) (*model.Organization, error) {

	org := &model.Organization{
		Name: name,
		Slug: slug,
	}
	insertStmt := Organization.INSERT(Organization.Name, Organization.Slug).MODEL(org).
		RETURNING(Organization.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, org)
	return org, err
}

func (o *OrganizationDomain) CreateOrganizationUser(
	ctx context.Context,
	tx *sql.Tx,
	organizationID, userID, roleID int64,
) (*model.OrganizationUser, error) {

	orgUser := &model.OrganizationUser{
		RoleID:         &roleID,
		UserID:         &userID,
		OrganizationID: &organizationID,
	}
	insertStmt := OrganizationUser.INSERT(
		OrganizationUser.RoleID,
		OrganizationUser.UserID,
		OrganizationUser.OrganizationID,
	).MODEL(orgUser).RETURNING(OrganizationUser.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, orgUser)
	return orgUser, err
}

func (o *OrganizationDomain) GetRoleByType(ctx context.Context, roleType resource.RoleType) (*model.Role, error) {
	stmt := SELECT(Role.AllColumns).FROM(Role).WHERE(Role.Type.EQ(Int(int64(roleType)))).LIMIT(1)

	role := &model.Role{}
	err := stmt.QueryContext(ctx, infra.DB, role)
	return role, err
}
