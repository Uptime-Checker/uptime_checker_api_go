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

var Check = newCheckTable("public", "check", "")

type checkTable struct {
	postgres.Table

	//Columns
	ID             postgres.ColumnInteger
	Body           postgres.ColumnString
	Traces         postgres.ColumnString
	Headers        postgres.ColumnString
	StatusCode     postgres.ColumnInteger
	ContentSize    postgres.ColumnInteger
	ContentType    postgres.ColumnString
	Duration       postgres.ColumnInteger
	Success        postgres.ColumnBool
	RegionID       postgres.ColumnInteger
	MonitorID      postgres.ColumnInteger
	OrganizationID postgres.ColumnInteger
	InsertedAt     postgres.ColumnTimestamp
	UpdatedAt      postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type CheckTable struct {
	checkTable

	EXCLUDED checkTable
}

// AS creates new CheckTable with assigned alias
func (a CheckTable) AS(alias string) *CheckTable {
	return newCheckTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new CheckTable with assigned schema name
func (a CheckTable) FromSchema(schemaName string) *CheckTable {
	return newCheckTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new CheckTable with assigned table prefix
func (a CheckTable) WithPrefix(prefix string) *CheckTable {
	return newCheckTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new CheckTable with assigned table suffix
func (a CheckTable) WithSuffix(suffix string) *CheckTable {
	return newCheckTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newCheckTable(schemaName, tableName, alias string) *CheckTable {
	return &CheckTable{
		checkTable: newCheckTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newCheckTableImpl("", "excluded", ""),
	}
}

func newCheckTableImpl(schemaName, tableName, alias string) checkTable {
	var (
		IDColumn             = postgres.IntegerColumn("id")
		BodyColumn           = postgres.StringColumn("body")
		TracesColumn         = postgres.StringColumn("traces")
		HeadersColumn        = postgres.StringColumn("headers")
		StatusCodeColumn     = postgres.IntegerColumn("status_code")
		ContentSizeColumn    = postgres.IntegerColumn("content_size")
		ContentTypeColumn    = postgres.StringColumn("content_type")
		DurationColumn       = postgres.IntegerColumn("duration")
		SuccessColumn        = postgres.BoolColumn("success")
		RegionIDColumn       = postgres.IntegerColumn("region_id")
		MonitorIDColumn      = postgres.IntegerColumn("monitor_id")
		OrganizationIDColumn = postgres.IntegerColumn("organization_id")
		InsertedAtColumn     = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn      = postgres.TimestampColumn("updated_at")
		allColumns           = postgres.ColumnList{IDColumn, BodyColumn, TracesColumn, HeadersColumn, StatusCodeColumn, ContentSizeColumn, ContentTypeColumn, DurationColumn, SuccessColumn, RegionIDColumn, MonitorIDColumn, OrganizationIDColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns       = postgres.ColumnList{BodyColumn, TracesColumn, HeadersColumn, StatusCodeColumn, ContentSizeColumn, ContentTypeColumn, DurationColumn, SuccessColumn, RegionIDColumn, MonitorIDColumn, OrganizationIDColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return checkTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:             IDColumn,
		Body:           BodyColumn,
		Traces:         TracesColumn,
		Headers:        HeadersColumn,
		StatusCode:     StatusCodeColumn,
		ContentSize:    ContentSizeColumn,
		ContentType:    ContentTypeColumn,
		Duration:       DurationColumn,
		Success:        SuccessColumn,
		RegionID:       RegionIDColumn,
		MonitorID:      MonitorIDColumn,
		OrganizationID: OrganizationIDColumn,
		InsertedAt:     InsertedAtColumn,
		UpdatedAt:      UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
