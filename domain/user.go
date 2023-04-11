package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type UserDomain struct{}

func NewUserDomain() *UserDomain {
	return &UserDomain{}
}

// Guest User

// CreateGuest does not need transaction
func (u *UserDomain) CreateGuest(ctx context.Context, email, code string) (*model.GuestUser, error) {
	now := times.Now()
	user := &model.GuestUser{
		Email:     email,
		Code:      code,
		ExpiresAt: now.Add(time.Minute * constant.GuestUserCodeExpiryInMinutes),
	}
	insertStmt := GuestUser.INSERT(GuestUser.Email, GuestUser.Code, GuestUser.ExpiresAt).MODEL(user).
		RETURNING(GuestUser.AllColumns)
	err := insertStmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) GetLatestGuestUser(ctx context.Context, email string) (*model.GuestUser, error) {
	stmt := SELECT(GuestUser.AllColumns).FROM(GuestUser).WHERE(GuestUser.Email.EQ(String(email))).
		ORDER_BY(GuestUser.ExpiresAt.DESC()).LIMIT(1)

	user := &model.GuestUser{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) GetGuestUser(ctx context.Context, email, code string) (*model.GuestUser, error) {
	stmt := SELECT(GuestUser.AllColumns).FROM(GuestUser).WHERE(
		GuestUser.Email.EQ(String(email)).AND(GuestUser.Code.EQ(String(code))),
	).ORDER_BY(GuestUser.ExpiresAt.DESC()).LIMIT(1)

	user := &model.GuestUser{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) DeleteGuestUser(ctx context.Context, tx *sql.Tx, id int64) error {
	deleteStmt := GuestUser.DELETE().WHERE(GuestUser.ID.EQ(Int(id)))
	_, err := deleteStmt.ExecContext(ctx, tx)
	return err
}

// User

func (u *UserDomain) GetUser(ctx context.Context, email string) (*model.User, error) {
	stmt := SELECT(User.AllColumns).FROM(User).WHERE(User.Email.EQ(String(email))).LIMIT(1)

	user := &model.User{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) GetUserWithRoleAndSubscription(
	ctx context.Context,
	id int64,
) (*pkg.UserWithRoleAndSubscription, error) {
	stmt := SELECT(
		User.AllColumns,
		Role.AllColumns,
		RoleClaim.AllColumns,
		Organization.AllColumns,
		Subscription.AllColumns,
		Plan.AllColumns,
		Product.AllColumns,
		Feature.AllColumns,
		ProductFeature.AllColumns,
	).
		FROM(
			User.
				LEFT_JOIN(Role, User.RoleID.EQ(Role.ID)).
				LEFT_JOIN(RoleClaim, User.RoleID.EQ(RoleClaim.RoleID)).
				LEFT_JOIN(Organization, User.OrganizationID.EQ(Organization.ID)).
				LEFT_JOIN(Subscription, User.OrganizationID.EQ(Subscription.OrganizationID)).
				LEFT_JOIN(Plan, Subscription.PlanID.EQ(Plan.ID)).
				LEFT_JOIN(Product, Subscription.ProductID.EQ(Product.ID)).
				LEFT_JOIN(ProductFeature, Product.ID.EQ(ProductFeature.ProductID)).
				LEFT_JOIN(Feature, ProductFeature.FeatureID.EQ(Feature.ID)),
		).
		WHERE(User.ID.EQ(Int(id)))

	user := &pkg.UserWithRoleAndSubscription{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) CreateUser(
	ctx context.Context,
	tx *sql.Tx,
	user *model.User,
	provider resource.UserLoginProvider,
) (*model.User, error) {
	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	now := times.Now()

	user.Provider = int32(provider)
	user.LastLoginAt = now
	insertStmt := User.INSERT(User.MutableColumns.Except(User.InsertedAt, User.UpdatedAt)).
		MODEL(user).
		RETURNING(User.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, user)
	return user, err
}

func (u *UserDomain) UpdateProvider(
	ctx context.Context,
	tx *sql.Tx,
	id int64,
	picture *string,
	providerUID string,
	provider resource.UserLoginProvider,
) (*model.User, error) {
	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	now := times.Now()
	user := &model.User{
		ProviderUID: &providerUID,
		Provider:    int32(provider),
		PictureURL:  picture,
		LastLoginAt: now,
		UpdatedAt:   now,
	}

	updateStmt := User.UPDATE(User.PictureURL, User.ProviderUID, User.Provider, User.LastLoginAt, User.UpdatedAt).
		MODEL(user).WHERE(User.ID.EQ(Int(id))).RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, user)
	return user, err
}

// UpdateName does not need transaction
func (u *UserDomain) UpdateName(
	ctx context.Context,
	id int64,
	name string,
) (*model.User, error) {
	now := times.Now()
	user := &model.User{
		Name:      &name,
		UpdatedAt: now,
	}

	updateStmt := User.UPDATE(User.Name, User.UpdatedAt).
		MODEL(user).WHERE(User.ID.EQ(Int(id))).RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) UpdateOrganizationAndRole(
	ctx context.Context,
	tx *sql.Tx,
	id, roleID, organizationID int64,
) (*model.User, error) {
	now := times.Now()
	user := &model.User{
		RoleID:         &roleID,
		OrganizationID: &organizationID,
		UpdatedAt:      now,
	}

	updateStmt := User.UPDATE(User.RoleID, User.OrganizationID, User.UpdatedAt).MODEL(user).WHERE(User.ID.EQ(Int(id))).
		RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, user)
	return user, err
}

func (u *UserDomain) UpdatePaymentID(
	ctx context.Context,
	id int64,
	paymentCustomerID string,
) (*model.User, error) {
	now := times.Now()
	user := &model.User{
		Name:      &paymentCustomerID,
		UpdatedAt: now,
	}

	updateStmt := User.UPDATE(User.PaymentCustomerID, User.UpdatedAt).
		MODEL(user).WHERE(User.ID.EQ(Int(id))).RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

// User Contact

func (u *UserDomain) CreateUserContact(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	email string,
	mode resource.UserContactMode,
	verified bool,
) (*model.UserContact, error) {
	if !mode.Valid() {
		return nil, constant.ErrInvalidUserContactMode
	}
	modeValue := int32(mode)
	userContact := &model.UserContact{
		UserID:   userID,
		Email:    &email,
		Mode:     &modeValue,
		Verified: verified,
	}
	insertStmt := UserContact.INSERT(UserContact.MutableColumns.Except(UserContact.InsertedAt, UserContact.UpdatedAt)).
		MODEL(userContact).RETURNING(UserContact.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, userContact)
	return userContact, err
}
