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

type Assertion struct {
	ID         int64 `sql:"primary_key"`
	Source     int32
	Property   *string
	Comparison int32
	Value      *string
	MonitorID  int64
	InsertedAt time.Time
	UpdatedAt  time.Time
}
