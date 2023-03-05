package pkg

import (
	"github.com/golang-jwt/jwt/v5"
)

// BearerClaims contains claims
type BearerClaims struct {
	UserID int64
	Email  string
	*jwt.RegisteredClaims
}
