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

type MonitorStatusChange struct {
	ID         int64 `sql:"primary_key"`
	Status     int32
	MonitorID  int64
	InsertedAt time.Time
	UpdatedAt  time.Time
}
