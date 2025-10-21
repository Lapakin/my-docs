package query

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// UpdateBuilder wraps squirrel.UpdateBuilder with additional functionality
type UpdateBuilder struct {
	ctx     context.Context
	builder sq.UpdateBuilder
}

// NewUpdate creates a new UpdateBuilder with Dollar placeholder format
func NewUpdate(ctx context.Context) *UpdateBuilder {
	return &UpdateBuilder{
		ctx:     ctx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update(""),
	}
}

// Table sets the table to update
func (u *UpdateBuilder) Table(table string) *UpdateBuilder {
	u.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update(table)
	return u
}

// Set sets a column value
func (u *UpdateBuilder) Set(column string, value interface{}) *UpdateBuilder {
	u.builder = u.builder.Set(column, value)
	return u
}

// SetMap sets multiple column values from a map
func (u *UpdateBuilder) SetMap(clauses map[string]interface{}) *UpdateBuilder {
	u.builder = u.builder.SetMap(clauses)
	return u
}

// Where adds a WHERE clause
func (u *UpdateBuilder) Where(pred interface{}, args ...interface{}) *UpdateBuilder {
	u.builder = u.builder.Where(pred, args...)
	return u
}

// WhereID adds a WHERE clause for id = value
func (u *UpdateBuilder) WhereID(id interface{}) *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{ColumnID: id})
	return u
}

// IsActive adds a WHERE clause for is_active = true
func (u *UpdateBuilder) IsActive() *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{ColumnIsActive: true})
	return u
}

// IsDeleted adds a WHERE clause for is_deleted = false (not deleted records by default)
func (u *UpdateBuilder) IsDeleted(deleted bool) *UpdateBuilder {
	u.builder = u.builder.Where(sq.Eq{ColumnIsDeleted: deleted})
	return u
}

// WhereMap adds WHERE clauses from a map
func (u *UpdateBuilder) WhereMap(clauses map[string]interface{}) *UpdateBuilder {
	for key, value := range clauses {
		u.builder = u.builder.Where(sq.Eq{key: value})
	}
	return u
}

// Returning adds RETURNING clause
func (u *UpdateBuilder) Returning(fields ...string) *UpdateBuilder {
	if len(fields) > 0 {
		u.builder = u.builder.Suffix("RETURNING " + fields[0])
		for i := 1; i < len(fields); i++ {
			u.builder = u.builder.Suffix(", " + fields[i])
		}
	}
	return u
}

// ToSql builds and returns the SQL query string and arguments
func (u *UpdateBuilder) ToSql() (query string, args []interface{}, err error) {
	query, args, err = u.builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build UPDATE query: %w", err)
	}
	return query, args, nil
}
