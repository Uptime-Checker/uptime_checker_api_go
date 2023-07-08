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

var Organization = newOrganizationTable("public", "organization", "")

type organizationTable struct {
	postgres.Table

	//Columns
	ID                    postgres.ColumnInteger
	Name                  postgres.ColumnString
	Slug                  postgres.ColumnString
	AlarmReminderInterval postgres.ColumnInteger
	AlarmReminderCount    postgres.ColumnInteger
	InsertedAt            postgres.ColumnTimestamp
	UpdatedAt             postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type OrganizationTable struct {
	organizationTable

	EXCLUDED organizationTable
}

// AS creates new OrganizationTable with assigned alias
func (a OrganizationTable) AS(alias string) *OrganizationTable {
	return newOrganizationTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new OrganizationTable with assigned schema name
func (a OrganizationTable) FromSchema(schemaName string) *OrganizationTable {
	return newOrganizationTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new OrganizationTable with assigned table prefix
func (a OrganizationTable) WithPrefix(prefix string) *OrganizationTable {
	return newOrganizationTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new OrganizationTable with assigned table suffix
func (a OrganizationTable) WithSuffix(suffix string) *OrganizationTable {
	return newOrganizationTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newOrganizationTable(schemaName, tableName, alias string) *OrganizationTable {
	return &OrganizationTable{
		organizationTable: newOrganizationTableImpl(schemaName, tableName, alias),
		EXCLUDED:          newOrganizationTableImpl("", "excluded", ""),
	}
}

func newOrganizationTableImpl(schemaName, tableName, alias string) organizationTable {
	var (
		IDColumn                    = postgres.IntegerColumn("id")
		NameColumn                  = postgres.StringColumn("name")
		SlugColumn                  = postgres.StringColumn("slug")
		AlarmReminderIntervalColumn = postgres.IntegerColumn("alarm_reminder_interval")
		AlarmReminderCountColumn    = postgres.IntegerColumn("alarm_reminder_count")
		InsertedAtColumn            = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn             = postgres.TimestampColumn("updated_at")
		allColumns                  = postgres.ColumnList{IDColumn, NameColumn, SlugColumn, AlarmReminderIntervalColumn, AlarmReminderCountColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns              = postgres.ColumnList{NameColumn, SlugColumn, AlarmReminderIntervalColumn, AlarmReminderCountColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return organizationTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                    IDColumn,
		Name:                  NameColumn,
		Slug:                  SlugColumn,
		AlarmReminderInterval: AlarmReminderIntervalColumn,
		AlarmReminderCount:    AlarmReminderCountColumn,
		InsertedAt:            InsertedAtColumn,
		UpdatedAt:             UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
