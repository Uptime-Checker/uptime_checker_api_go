//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var AlarmPolicy = newAlarmPolicyTable("public", "alarm_policy", "")

type alarmPolicyTable struct {
	postgres.Table

	//Columns
	ID             postgres.ColumnInteger
	Reason         postgres.ColumnString
	Threshold      postgres.ColumnInteger
	MonitorID      postgres.ColumnInteger
	OrganizationID postgres.ColumnInteger
	InsertedAt     postgres.ColumnTimestamp
	UpdatedAt      postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AlarmPolicyTable struct {
	alarmPolicyTable

	EXCLUDED alarmPolicyTable
}

// AS creates new AlarmPolicyTable with assigned alias
func (a AlarmPolicyTable) AS(alias string) *AlarmPolicyTable {
	return newAlarmPolicyTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AlarmPolicyTable with assigned schema name
func (a AlarmPolicyTable) FromSchema(schemaName string) *AlarmPolicyTable {
	return newAlarmPolicyTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AlarmPolicyTable with assigned table prefix
func (a AlarmPolicyTable) WithPrefix(prefix string) *AlarmPolicyTable {
	return newAlarmPolicyTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AlarmPolicyTable with assigned table suffix
func (a AlarmPolicyTable) WithSuffix(suffix string) *AlarmPolicyTable {
	return newAlarmPolicyTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAlarmPolicyTable(schemaName, tableName, alias string) *AlarmPolicyTable {
	return &AlarmPolicyTable{
		alarmPolicyTable: newAlarmPolicyTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newAlarmPolicyTableImpl("", "excluded", ""),
	}
}

func newAlarmPolicyTableImpl(schemaName, tableName, alias string) alarmPolicyTable {
	var (
		IDColumn             = postgres.IntegerColumn("id")
		ReasonColumn         = postgres.StringColumn("reason")
		ThresholdColumn      = postgres.IntegerColumn("threshold")
		MonitorIDColumn      = postgres.IntegerColumn("monitor_id")
		OrganizationIDColumn = postgres.IntegerColumn("organization_id")
		InsertedAtColumn     = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn      = postgres.TimestampColumn("updated_at")
		allColumns           = postgres.ColumnList{IDColumn, ReasonColumn, ThresholdColumn, MonitorIDColumn, OrganizationIDColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns       = postgres.ColumnList{ReasonColumn, ThresholdColumn, MonitorIDColumn, OrganizationIDColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return alarmPolicyTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		Reason:         ReasonColumn,
		Threshold:      ThresholdColumn,
		MonitorID:      MonitorIDColumn,
		OrganizationID: OrganizationIDColumn,
		InsertedAt:     InsertedAtColumn,
		UpdatedAt:      UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}