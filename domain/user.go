package domain

import (
	"time"

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
	user := model.GuestUser{Email: email, Code: pkg.GetUniqueString(), ExpiresAt: now.Add(time.Minute * 10), InsertedAt: now, UpdatedAt: now}
	insertStmt := GuestUser.INSERT(GuestUser.Email, GuestUser.Code, GuestUser.ExpiresAt, GuestUser.InsertedAt, GuestUser.UpdatedAt).MODEL(user).
		RETURNING(GuestUser.AllColumns)
	err := insertStmt.Query(infra.DB, &user)
	return &user, err
}
