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

var Feature = newFeatureTable("public", "feature", "")

type featureTable struct {
	postgres.Table

	//Columns
	ID         postgres.ColumnInteger
	Name       postgres.ColumnString
	Type       postgres.ColumnInteger
	InsertedAt postgres.ColumnTimestamp
	UpdatedAt  postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type FeatureTable struct {
	featureTable

	EXCLUDED featureTable
}

// AS creates new FeatureTable with assigned alias
func (a FeatureTable) AS(alias string) *FeatureTable {
	return newFeatureTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FeatureTable with assigned schema name
func (a FeatureTable) FromSchema(schemaName string) *FeatureTable {
	return newFeatureTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FeatureTable with assigned table prefix
func (a FeatureTable) WithPrefix(prefix string) *FeatureTable {
	return newFeatureTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FeatureTable with assigned table suffix
func (a FeatureTable) WithSuffix(suffix string) *FeatureTable {
	return newFeatureTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFeatureTable(schemaName, tableName, alias string) *FeatureTable {
	return &FeatureTable{
		featureTable: newFeatureTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newFeatureTableImpl("", "excluded", ""),
	}
}

func newFeatureTableImpl(schemaName, tableName, alias string) featureTable {
	var (
		IDColumn         = postgres.IntegerColumn("id")
		NameColumn       = postgres.StringColumn("name")
		TypeColumn       = postgres.IntegerColumn("type")
		InsertedAtColumn = postgres.TimestampColumn("inserted_at")
		UpdatedAtColumn  = postgres.TimestampColumn("updated_at")
		allColumns       = postgres.ColumnList{IDColumn, NameColumn, TypeColumn, InsertedAtColumn, UpdatedAtColumn}
		mutableColumns   = postgres.ColumnList{NameColumn, TypeColumn, InsertedAtColumn, UpdatedAtColumn}
	)

	return featureTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:         IDColumn,
		Name:       NameColumn,
		Type:       TypeColumn,
		InsertedAt: InsertedAtColumn,
		UpdatedAt:  UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
