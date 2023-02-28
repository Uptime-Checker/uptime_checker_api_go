//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type User struct {
	ID                int64 `sql:"primary_key"`
	Name              string
	Email             string
	PictureURL        *string
	Password          *string
	PaymentCustomerID *string
	ProviderUID       *string
	Provider          *int32
	LastLoginAt       *time.Time
	RoleID            *int64
	OrganizationID    *int64
	InsertedAt        time.Time
	UpdatedAt         time.Time
}
