package domain

import (
	"context"
	"database/sql"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"

	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
	. "github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/table"
)

type UserDomain struct{}

func NewUserDomain() *UserDomain {
	return &UserDomain{}
}

func (u *UserDomain) CreateGuest(ctx context.Context, email, code string) (*model.GuestUser, error) {
	now := times.Now()
	user := &model.GuestUser{Email: email, Code: code, ExpiresAt: now.Add(time.Minute * 10)}
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
	stmt := SELECT(GuestUser.AllColumns).FROM(GuestUser).WHERE(GuestUser.Email.EQ(String(email))).
		WHERE(GuestUser.Code.EQ(String(code))).
		ORDER_BY(GuestUser.ExpiresAt.DESC()).LIMIT(1)

	user := &model.GuestUser{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) GetUser(ctx context.Context, email string) (*model.User, error) {
	stmt := SELECT(User.AllColumns).FROM(User).WHERE(User.Email.EQ(String(email))).LIMIT(1)

	user := &model.User{}
	err := stmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) CreateUser(
	ctx context.Context,
	tx *sql.Tx,
	email string,
	provider resource.UserLoginProvider,
) (*model.User, error) {

	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	providerValue := int32(provider)
	now := times.Now()
	user := &model.User{
		Email:       email,
		ProviderUID: &email,
		Provider:    &providerValue,
		LastLoginAt: &now,
	}
	insertStmt := User.INSERT(User.Email, User.ProviderUID, User.Provider, User.LastLoginAt).MODEL(user).
		RETURNING(User.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, user)
	return user, err
}

func (u *UserDomain) UpdateProvider(
	ctx context.Context,
	tx *sql.Tx,
	email string,
	provider resource.UserLoginProvider,
) (*model.User, error) {

	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	providerValue := int32(provider)
	now := times.Now()
	user := &model.User{
		ProviderUID: &email,
		Provider:    &providerValue,
		LastLoginAt: &now,
		UpdatedAt:   now,
	}

	updateStmt := User.UPDATE(User.ProviderUID, User.Provider, User.LastLoginAt, User.UpdatedAt).
		MODEL(user).WHERE(User.Email.EQ(String(email))).RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, tx, user)
	return user, err
}

func (u *UserDomain) UpdateName(
	ctx context.Context,
	id int64,
	name string,
) (*model.User, error) {

	now := times.Now()
	user := &model.User{
		Name:      name,
		UpdatedAt: now,
	}

	updateStmt := User.UPDATE(User.Name, User.UpdatedAt).
		MODEL(user).WHERE(User.ID.EQ(Int(id))).RETURNING(User.AllColumns)

	err := updateStmt.QueryContext(ctx, infra.DB, user)
	return user, err
}

func (u *UserDomain) DeleteGuestUser(ctx context.Context, tx *sql.Tx, id int64) error {
	deleteStmt := GuestUser.DELETE().WHERE(GuestUser.ID.EQ(Int(id)))
	_, err := deleteStmt.ExecContext(ctx, tx)
	return err
}

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
		UserID:   &userID,
		Email:    &email,
		Mode:     &modeValue,
		Verified: verified,
	}
	insertStmt := UserContact.INSERT(UserContact.MutableColumns).MODEL(userContact).
		RETURNING(UserContact.AllColumns)
	err := insertStmt.QueryContext(ctx, tx, userContact)
	return userContact, err
}
