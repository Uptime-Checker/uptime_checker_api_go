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

var Property = newPropertyTable("public", "property", "")

type propertyTable struct {
	postgres.Table

	//Columns
	ID         postgres.ColumnInteger
	Key        postgres.ColumnString
	Value      postgres.ColumnString
	InsertedAt postgres.ColumnTimestamp
	UpdatedAt  postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type PropertyTable struct {
	propertyTable

	EXCLUDED propertyTable
}

// AS creates new PropertyTable with assigned alias
func (a PropertyTable) AS(alias string) *PropertyTable {
	return newPropertyTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PropertyTable with assigned schema name
func (a PropertyTable) FromSchema(schemaName string) *PropertyTable {
	return newPropertyTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PropertyTable with assigned table prefix
func (a PropertyTable) WithPrefix(prefix string) *PropertyTable {
	return newPropertyTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PropertyTable with assigned table suffix
func (a PropertyTable) WithSuffix(suffix string) *PropertyTable {
	return newPropertyTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPropertyTable(schemaName, tableName, alias string) *PropertyTable {
	return &PropertyTable{
		propertyTable: newPropertyTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newPropertyTableImpl("", "excluded", ""),
	}
}

func newPropertyTableImpl(schemaName, tableName, alias string) propertyTable {
	var (
		IDColumn         = postgres.IntegerColumn("id")
		KeyColumn        = postgres.StringColumn("key")
		ValueColumn      = postgres.StringColumn("value")
		InsertedAtColumn = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn  = postgres.TimestampColumn("updated_at")
		allColumns       = postgres.ColumnList{IDColumn, KeyColumn, ValueColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns   = postgres.ColumnList{KeyColumn, ValueColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return propertyTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		Key:        KeyColumn,
		Value:      ValueColumn,
		InsertedAt: InsertedAtColumn,
		UpdatedAt:  UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
