package domain

import (
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/validate"
	"time"

	. "github.com/go-jet/jet/v2/postgres"

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

func (u *UserDomain) CreateUser(email string, provider validate.UserLoginProvider) (*model.User, error) {
	now := times.Now()
	user := &model.User{Email: email, Code: code, ExpiresAt: now.Add(time.Minute * 10)}
	insertStmt := User.INSERT(GuestUser.Email, GuestUser.Code, GuestUser.ExpiresAt).MODEL(user).
		RETURNING(User.AllColumns)
	err := insertStmt.Query(infra.DB, user)
	return user, err
}
