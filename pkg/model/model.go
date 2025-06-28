package model

type Table struct {
	Name    string   `json:"tableName"`
	Schema  string   `json:"tableSchema,omitempty"`
	Columns []Column `json:"tableColumns"`
}

type ColumnType string

const (
	// Numbers
	SmallInt ColumnType = "smallint"
	Int      ColumnType = "int4"
	BigInt   ColumnType = "bigint"
	Real     ColumnType = "real"
	Double   ColumnType = "double"
	Decimal  ColumnType = "decimal"

	// Text
	Char    ColumnType = "char"
	Varchar ColumnType = "varchar"
	Text    ColumnType = "text"

	// Special
	UUID    ColumnType = "uuid"
	Boolean ColumnType = "bool"

	// Date & Time
	Date        ColumnType = "date"
	Time        ColumnType = "time"
	Timestamp   ColumnType = "timestamp"
	TimestampTZ ColumnType = "timestamptz"
)

// TODO: missing a few fields like precision, scale, etc.
type Column struct {
	Name      string     `json:"columnName"`
	Typ       ColumnType `json:"columnType"`
	MaxLength uint       `json:"maxLength,omitempty"`

	IsUnique    bool   `json:"isUnique"`
	IsGenerated bool   `json:"isGenerated"`
	IsNullable  bool   `json:"isNullable"`
	ForeignKey  string `json:"foreginKey,omitempty"`

	Annotation string `json:"annotation,omitempty"`
}
