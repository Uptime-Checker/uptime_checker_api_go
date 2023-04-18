package domain

import (
	"context"
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type OrganizationDomain struct{}

func NewOrganizationDomain() *OrganizationDomain {
	return &OrganizationDomain{}
}

func (o *OrganizationDomain) Create(
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
		RoleID:         roleID,
		UserID:         userID,
		OrganizationID: organizationID,
	}
	insertStmt := OrganizationUser.INSERT(
		OrganizationUser.RoleID,
		OrganizationUser.UserID,
		OrganizationUser.OrganizationID,
	).MODEL(orgUser).RETURNING(OrganizationUser.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, orgUser)
	return orgUser, err
}

func (o *OrganizationDomain) Get(ctx context.Context, id int64) (*model.Organization, error) {
	stmt := SELECT(Organization.AllColumns).FROM(Organization).WHERE(Organization.ID.EQ(Int(id))).LIMIT(1)

	organization := &model.Organization{}
	err := stmt.QueryContext(ctx, infra.DB, organization)
	return organization, err
}

func (o *OrganizationDomain) GetRoleByType(ctx context.Context, roleType resource.RoleType) (*model.Role, error) {
	stmt := SELECT(Role.AllColumns).FROM(Role).WHERE(Role.Type.EQ(Int(int64(roleType)))).LIMIT(1)

	role := &model.Role{}
	err := stmt.QueryContext(ctx, infra.DB, role)
	return role, err
}

func (o *OrganizationDomain) ListOrganizationsOfUser(
	ctx context.Context,
	userID int64,
) ([]pkg.OrganizationUserRole, error) {
	stmt := SELECT(OrganizationUser.AllColumns, Organization.AllColumns, Role.AllColumns).
		FROM(
			OrganizationUser.
				LEFT_JOIN(Role, OrganizationUser.RoleID.EQ(Role.ID)).
				LEFT_JOIN(Organization, OrganizationUser.OrganizationID.EQ(Organization.ID)),
		).
		WHERE(OrganizationUser.UserID.EQ(Int(userID)))

	var organizationUserRoles []pkg.OrganizationUserRole
	err := stmt.QueryContext(ctx, infra.DB, &organizationUserRoles)
	return organizationUserRoles, err
}
