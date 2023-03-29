package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"

	"github.com/Uptime-Checker/uptime_checker_api_go/cache"
	"github.com/Uptime-Checker/uptime_checker_api_go/config"
	"github.com/Uptime-Checker/uptime_checker_api_go/constant"
	"github.com/Uptime-Checker/uptime_checker_api_go/domain"
	"github.com/Uptime-Checker/uptime_checker_api_go/infra/lgr"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg"
	"github.com/Uptime-Checker/uptime_checker_api_go/pkg/times"
	"github.com/Uptime-Checker/uptime_checker_api_go/schema/uptime_checker/public/model"
)

type AuthService struct {
	userDomain *domain.UserDomain
}

func NewAuthService(userDomain *domain.UserDomain) *AuthService {
	return &AuthService{userDomain: userDomain}
}

func (a *AuthService) VerifyGuestUser(ctx context.Context, email, code string) (*model.GuestUser, error) {
	now := times.Now()
	tracingID := pkg.GetTracingID(ctx)

	guestUser, err := a.userDomain.GetGuestUser(ctx, email, code)
	if err != nil {
		lgr.Error(tracingID, 1, "no guest user with", email, "code:", code)
		return nil, err
	}
	if times.CompareDate(now, guestUser.ExpiresAt) == constant.Date1AfterDate2 {
		lgr.Print(tracingID, 2, "no guest expired", guestUser.ExpiresAt, "now:", now)
		return nil, constant.ErrGuestUserCodeExpired
	}
	return guestUser, nil
}

func (a *AuthService) GenerateUserToken(user *model.User) (string, error) {
	now := times.Now()
	expirationTime := now.Add(constant.BearerTokenExpirationInDays * 24 * time.Hour)
	claims := &pkg.BearerClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Audience:  jwt.ClaimStrings{user.Email},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTKey)
}

// GetUserByToken returns user by token
func (a *AuthService) GetUserByToken(ctx context.Context, tok string) (*pkg.UserWithRoleAndSubscription, error) {
	token, err := jwt.Parse(tok, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constant.ErrInvalidUserToken
		}
		return config.JWTKey, nil
	})
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, constant.ErrInvalidUserToken
	}

	usr := pkg.BearerClaims{}
	if err := mapstructure.Decode(token.Claims, &usr); err != nil {
		return nil, err
	}

	cachedUser := cache.GetUserWithRoleAndSubscription(usr.UserID)
	if cachedUser == nil {
		user, err := a.userDomain.GetUserWithRoleAndSubscription(ctx, usr.UserID)
		if err != nil {
			return nil, err
		}
		cache.SetUserWithRoleAndSubscription(user)
		return user, nil
	}
	return cachedUser, nil
}
