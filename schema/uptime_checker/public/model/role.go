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

type Role struct {
	ID         int64 `sql:"primary_key"`
	Name       string
	Type       *int32
	InsertedAt time.Time
	UpdatedAt  time.Time
}
