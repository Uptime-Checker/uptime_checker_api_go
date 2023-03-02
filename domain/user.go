package domain

import (
	"time"

	. "github.com/go-jet/jet/v2/postgres"

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

func (u *UserDomain) CreateGuest(email string) (*model.GuestUser, error) {
	now := times.Now()
	code := pkg.GetUniqueString()
	user := &model.GuestUser{Email: email, Code: pkg.HashSha(code), ExpiresAt: now.Add(time.Minute * 10)}
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
