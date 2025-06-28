package adapter

import (
	"database/sql"
	"dbaker/pkg/config"
	"dbaker/pkg/generator"
	"dbaker/pkg/model"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgreSQLAdapter struct {
	config config.Config
	db     *sql.DB
	gen    *generator.ValueGenerator
}

func NewPostgreSQLAdapter(config config.Config) PostgreSQLAdapter {
	return PostgreSQLAdapter{
		config: config,
		db:     nil,
		gen:    nil,
	}
}

func (p *PostgreSQLAdapter) Init() error {
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.config.Host,
		p.config.Port,
		p.config.Username,
		p.config.Password,
		p.config.Database,
		p.config.SSLMode,
	)

	db, err := sql.Open("pgx", connection)
	if err != nil {
		return fmt.Errorf("failed to init a database connection: %w", err)
	}

	p.db = db
	return nil
}

func (p *PostgreSQLAdapter) Close() error {
	return p.db.Close()
}

func (p *PostgreSQLAdapter) IntrospectTable(name string, schema string) (*model.Table, error) {
	tbl, err := p.findTable(name, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to find table: %w", err)
	}

	infoSchemaColumns, err := p.findTableColumns(name, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to find table columns: %w", err)
	}

	constraints, err := p.findTableConstraints(name, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to find table constraints: %w", err)
	}

	var columns []model.Column
	for _, infoSchemaColumn := range infoSchemaColumns {
		column := infoSchemaColumn.mapToColumn()

		for _, constraint := range constraints {
			if isUnique(column, constraint) {
				column.IsUnique = true
				break
			}
		}

		columns = append(columns, column)
	}

	table := model.Table{
		Name:    *tbl.TableName,
		Schema:  *tbl.TableSchema,
		Columns: columns,
	}
	return &table, nil
}

const FIND_TABLE_BY_NAME_AND_SCHEMA_QUERY = `
select
	table_schema,
	table_name
from
	information_schema.tables
where
	table_schema = $1
and
	table_name = $2
`

type InfoSchemaTable struct {
	TableName   *string
	TableSchema *string
}

func (p *PostgreSQLAdapter) findTable(table string, schema string) (*InfoSchemaTable, error) {
	statement, err := p.db.Prepare(FIND_TABLE_BY_NAME_AND_SCHEMA_QUERY)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a query statement: %w", err)
	}

	row := statement.QueryRow(schema, table)

	var tbl InfoSchemaTable
	err = row.Scan(
		&tbl.TableName,
		&tbl.TableSchema,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table row: %w", err)
	}

	return &tbl, nil
}

const FIND_TABLE_COLUMNS_BY_NAME_AND_SCHEMA_QUERY = `
select
	column_name,
	udt_name,
	character_maximum_length,
	is_nullable,
	is_identity
from
	information_schema.columns
where
	table_schema = $1
and
	table_name = $2;
`

type InfoSchemaColumn struct {
	ColumnName             *string
	UdtName                *string
	CharacterMaximumLength *uint
	IsNullable             *string
	IsIdentity             *string
}

func (c InfoSchemaColumn) mapToColumn() model.Column {
	var column model.Column
	column.Name = *c.ColumnName

	// TODO: this might require a proper mapping from UDT to model column type
	column.Typ = model.ColumnType(*c.UdtName)

	if c.CharacterMaximumLength != nil {
		column.MaxLength = *c.CharacterMaximumLength
	}

	if c.IsIdentity != nil && *c.IsIdentity == "YES" {
		column.IsGenerated = true
	}

	if c.IsNullable != nil && *c.IsNullable == "YES" {
		column.IsNullable = true
	}

	return column
}

func (p *PostgreSQLAdapter) findTableColumns(table string, schema string) ([]InfoSchemaColumn, error) {
	statement, err := p.db.Prepare(FIND_TABLE_COLUMNS_BY_NAME_AND_SCHEMA_QUERY)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a query statement: %w", err)
	}

	rows, err := statement.Query(schema, table)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema.columns table: %w", err)
	}
	defer rows.Close()

	var columns []InfoSchemaColumn
	for rows.Next() {
		var column InfoSchemaColumn
		if err := rows.Scan(
			&column.ColumnName,
			&column.UdtName,
			&column.CharacterMaximumLength,
			&column.IsNullable,
			&column.IsIdentity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		columns = append(columns, column)
	}

	return columns, nil
}

const FIND_TABLE_CONSTRAINTS_BY_NAME_AND_SCHEMA_QUERY = `
select
	ku.constraint_name,
	tc.constraint_type,
	ku.column_name,
	ku.ordinal_position
from
	information_schema.key_column_usage as ku
left join
	information_schema.table_constraints as tc
on
	ku.constraint_name = tc.constraint_name
and
	ku.table_schema = tc.table_schema
and
	ku.table_name = tc.table_name
where
	ku.table_schema = $1
and
	ku.table_name = $2;
`

type InfoSchemaConstraint struct {
	ConstraintName  *string
	ConstraintType  *string
	ColumnName      *string
	OrdinalPosition uint
}

func (p *PostgreSQLAdapter) findTableConstraints(name string, schema string) ([]InfoSchemaConstraint, error) {
	statement, err := p.db.Prepare(FIND_TABLE_CONSTRAINTS_BY_NAME_AND_SCHEMA_QUERY)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a query statement: %w", err)
	}

	rows, err := statement.Query(schema, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query information_schema table_constraints or key_column_usage tables: %w", err)
	}
	defer rows.Close()

	var constraints []InfoSchemaConstraint
	for rows.Next() {
		var constraint InfoSchemaConstraint
		if err := rows.Scan(
			&constraint.ConstraintName,
			&constraint.ConstraintType,
			&constraint.ColumnName,
			&constraint.OrdinalPosition,
		); err != nil {
			return nil, fmt.Errorf("failed to scan constraint: %w", err)
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

func isUnique(column model.Column, constraint InfoSchemaConstraint) bool {
	return column.Name == *constraint.ColumnName &&
		(*constraint.ConstraintType == "UNIQUE" ||
			*constraint.ConstraintType == "PRIMARY KEY")
}

// no batch support yet
// generated values (infered, identities etc should not be present at this point)
// insert into <schema>.<table> (<for-earch column.Name>,) values (for-each column '?')
func (p *PostgreSQLAdapter) WriteRow(table string, schema string, columns []model.Column, iter uint32) error {
	columnNames := inferColNames(columns)
	placeholders := inferPgValPlaceholders(len(columns))
	insertQuery := fmt.Sprintf("insert into %s.%s (%s) values (%s);",
		table, schema, columnNames, placeholders)

	columnValues, err := p.gen.GenVals(columns, iter)
	if err != nil {
		return fmt.Errorf("failed to generate row values for iteration '%d': %w", iter, err)
	}

	fmt.Printf("Generated insert query: %s\n", insertQuery)
	fmt.Printf("Generated values: %v\n", columnValues)

	stmt, err := p.db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(columnValues...)
	if err != nil {
		return fmt.Errorf("failed to insert data to table '%s.%s' on iteration '%d': %w",
			schema, table, iter, err)
	}

	return nil
}

func inferColNames(columns []model.Column) string {
	builder := strings.Builder{}
	for index, column := range columns {
		builder.WriteString(column.Name)

		if index < len(columns)-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func inferPgValPlaceholders(columnLen int) string {
	builder := strings.Builder{}
	for index := range columnLen {
		builder.WriteString(fmt.Sprintf("$%d", index+1))

		if index < columnLen-1 {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}
