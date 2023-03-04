package domain

import (
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

func (u *UserDomain) CreateGuest(email, code string) (*model.GuestUser, error) {
	now := times.Now()
	user := &model.GuestUser{Email: email, Code: code, ExpiresAt: now.Add(time.Minute * 10)}
	insertStmt := GuestUser.INSERT(GuestUser.Email, GuestUser.Code, GuestUser.ExpiresAt).MODEL(user).
		RETURNING(GuestUser.AllColumns)
	err := insertStmt.Query(infra.DB, user)
	return user, err
}

func (u *UserDomain) GetLatestGuestUser(email string) (*model.GuestUser, error) {
	stmt := SELECT(GuestUser.AllColumns).FROM(GuestUser).WHERE(GuestUser.Email.EQ(String(email))).
		ORDER_BY(GuestUser.ExpiresAt.DESC()).LIMIT(1)

	user := &model.GuestUser{}
	err := stmt.Query(infra.DB, user)
	return user, err
}

func (u *UserDomain) GetGuestUser(email, code string) (*model.GuestUser, error) {
	stmt := SELECT(GuestUser.AllColumns).FROM(GuestUser).WHERE(GuestUser.Email.EQ(String(email))).
		WHERE(GuestUser.Code.EQ(String(code))).
		ORDER_BY(GuestUser.ExpiresAt.DESC()).LIMIT(1)

	user := &model.GuestUser{}
	err := stmt.Query(infra.DB, user)
	return user, err
}

func (u *UserDomain) GetUser(email string) (*model.User, error) {
	stmt := SELECT(User.AllColumns).FROM(User).WHERE(User.Email.EQ(String(email))).LIMIT(1)

	user := &model.User{}
	err := stmt.Query(infra.DB, user)
	return user, err
}

func (u *UserDomain) CreateUser(email string, provider resource.UserLoginProvider) (*model.User, error) {
	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	providerValue := provider.Value()
	now := times.Now()
	user := &model.User{
		Email:       email,
		ProviderUID: &email,
		Provider:    &providerValue,
		LastLoginAt: &now,
		UpdatedAt:   time.Time{},
	}
	insertStmt := User.INSERT(User.Email, User.ProviderUID, User.Provider, User.LastLoginAt, User.UpdatedAt).MODEL(user).
		RETURNING(User.AllColumns)
	err := insertStmt.Query(infra.DB, user)
	return user, err
}

func (u *UserDomain) UpdateProvider(email string, provider resource.UserLoginProvider) (*model.User, error) {
	if !provider.Valid() {
		return nil, constant.ErrInvalidProvider
	}
	providerValue := provider.Value()
	now := times.Now()
	user := &model.User{
		ProviderUID: &email,
		Provider:    &providerValue,
		LastLoginAt: &now,
		UpdatedAt:   now,
	}

	updateStmt := User.UPDATE(User.ProviderUID, User.Provider, User.LastLoginAt, User.UpdatedAt).
		MODEL(user).WHERE(User.Email.EQ(String(email))).RETURNING(User.AllColumns)

	err := updateStmt.Query(infra.DB, user)
	return user, err
}
