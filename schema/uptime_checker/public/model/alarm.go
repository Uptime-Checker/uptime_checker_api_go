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

type Alarm struct {
	ID                 int64 `sql:"primary_key"`
	Ongoing            *bool
	ResolvedAt         *time.Time
	TriggeredByCheckID *int64
	ResolvedByCheckID  *int64
	MonitorID          *int64
	OrganizationID     *int64
	InsertedAt         time.Time
	UpdatedAt          time.Time
}
