package service

import (
	"context"

	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/log"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type UserService struct {
	userDomain *domain.UserDomain
}

func NewUserService(userDomain *domain.UserDomain) *UserService {
	return &UserService{userDomain: userDomain}
}

func (u *UserService) VerifyGuestUser(ctx context.Context, email, code string) (*model.GuestUser, error) {
	now := times.Now()
	tracingID := pkg.GetTracingID(ctx)

	guestUser, err := u.userDomain.GetGuestUser(email, code)
	if err != nil {
		log.Default.Error(tracingID, 1, "no guest user with", email, "code:", code)
		return nil, err
	}
	if times.CompareDate(now, guestUser.ExpiresAt) == constant.Date1AfterDate2 {
		log.Default.Print(tracingID, 2, "no guest expired", guestUser.ExpiresAt, "now:", now)
		return nil, constant.ErrGuestUserCodeExpired
	}
	return guestUser, nil
}
