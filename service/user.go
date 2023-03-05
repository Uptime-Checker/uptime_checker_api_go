package service

import (
	"context"
	"database/sql"

	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain/resource"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type UserService struct {
	userDomain *domain.UserDomain
}

func NewUserService(userDomain *domain.UserDomain) *UserService {
	return &UserService{userDomain: userDomain}
}

func (u *UserService) CreateNewUserAndContact(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	user, err := u.userDomain.CreateUser(ctx, tx, email, resource.UserLoginProviderEmail)
	if err != nil {
		return nil, err
	}
	_, err = u.userDomain.CreateUserContact(ctx, tx, user.ID, email, resource.UserContactModeEmail, true)
	if err != nil {
		return nil, err
	}
	return user, nil
}
