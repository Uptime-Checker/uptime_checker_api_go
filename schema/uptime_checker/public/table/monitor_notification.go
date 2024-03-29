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

var MonitorNotification = newMonitorNotificationTable("public", "monitor_notification", "")

type monitorNotificationTable struct {
	postgres.Table

	//Columns
	ID             postgres.ColumnInteger
	Type           postgres.ColumnInteger
	ExternalID     postgres.ColumnString
	Successful     postgres.ColumnBool
	AlarmID        postgres.ColumnInteger
	MonitorID      postgres.ColumnInteger
	UserContactID  postgres.ColumnInteger
	OrganizationID postgres.ColumnInteger
	IntegrationID  postgres.ColumnInteger
	InsertedAt     postgres.ColumnTimestamp
	UpdatedAt      postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MonitorNotificationTable struct {
	monitorNotificationTable

	EXCLUDED monitorNotificationTable
}

// AS creates new MonitorNotificationTable with assigned alias
func (a MonitorNotificationTable) AS(alias string) *MonitorNotificationTable {
	return newMonitorNotificationTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MonitorNotificationTable with assigned schema name
func (a MonitorNotificationTable) FromSchema(schemaName string) *MonitorNotificationTable {
	return newMonitorNotificationTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MonitorNotificationTable with assigned table prefix
func (a MonitorNotificationTable) WithPrefix(prefix string) *MonitorNotificationTable {
	return newMonitorNotificationTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MonitorNotificationTable with assigned table suffix
func (a MonitorNotificationTable) WithSuffix(suffix string) *MonitorNotificationTable {
	return newMonitorNotificationTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMonitorNotificationTable(schemaName, tableName, alias string) *MonitorNotificationTable {
	return &MonitorNotificationTable{
		monitorNotificationTable: newMonitorNotificationTableImpl(schemaName, tableName, alias),
		EXCLUDED:                 newMonitorNotificationTableImpl("", "excluded", ""),
	}
}

func newMonitorNotificationTableImpl(schemaName, tableName, alias string) monitorNotificationTable {
	var (
		IDColumn             = postgres.IntegerColumn("id")
		TypeColumn           = postgres.IntegerColumn("type")
		ExternalIDColumn     = postgres.StringColumn("external_id")
		SuccessfulColumn     = postgres.BoolColumn("successful")
		AlarmIDColumn        = postgres.IntegerColumn("alarm_id")
		MonitorIDColumn      = postgres.IntegerColumn("monitor_id")
		UserContactIDColumn  = postgres.IntegerColumn("user_contact_id")
		OrganizationIDColumn = postgres.IntegerColumn("organization_id")
		IntegrationIDColumn  = postgres.IntegerColumn("integration_id")
		InsertedAtColumn     = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn      = postgres.TimestampColumn("updated_at")
		allColumns           = postgres.ColumnList{IDColumn, TypeColumn, ExternalIDColumn, SuccessfulColumn, AlarmIDColumn, MonitorIDColumn, UserContactIDColumn, OrganizationIDColumn, IntegrationIDColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns       = postgres.ColumnList{TypeColumn, ExternalIDColumn, SuccessfulColumn, AlarmIDColumn, MonitorIDColumn, UserContactIDColumn, OrganizationIDColumn, IntegrationIDColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return monitorNotificationTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		Type:           TypeColumn,
		ExternalID:     ExternalIDColumn,
		Successful:     SuccessfulColumn,
		AlarmID:        AlarmIDColumn,
		MonitorID:      MonitorIDColumn,
		UserContactID:  UserContactIDColumn,
		OrganizationID: OrganizationIDColumn,
		IntegrationID:  IntegrationIDColumn,
		InsertedAt:     InsertedAtColumn,
		UpdatedAt:      UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
