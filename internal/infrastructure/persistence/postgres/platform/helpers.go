package platform

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// buildFilterQuery returns a FROM clause and args for an optional activeOnly filter.
// Example: buildFilterQuery("crypto_assets", "is_active", true)
//
//	-> ("crypto_assets WHERE is_active = $1", []interface{}{true})
func buildFilterQuery(table, activeCol string, activeOnly bool, searchQuery string, searchCols []string) (string, []interface{}) {
	where := []string{}
	args := []interface{}{}

	if activeOnly {
		args = append(args, true)
		where = append(where, fmt.Sprintf("%s = $%d", activeCol, len(args)))
	}

	if searchQuery != "" {
		orParts := []string{}
		args = append(args, "%"+searchQuery+"%")
		searchArgIdx := len(args)
		for _, col := range searchCols {
			orParts = append(orParts, fmt.Sprintf("%s ILIKE $%d", col, searchArgIdx))
		}
		where = append(where, "("+strings.Join(orParts, " OR ")+")")
	}

	query := table
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	return query, args
}

// buildPagedQuery appends LIMIT / OFFSET placeholders to a base SELECT query.
func buildPagedQuery(base string, args []interface{}, limit, offset int) (string, []interface{}) {
	ph := len(args) + 1
	if limit > 0 {
		base += fmt.Sprintf(" LIMIT $%d", ph)
		args = append(args, limit)
		ph++
	}
	if offset > 0 {
		base += fmt.Sprintf(" OFFSET $%d", ph)
		args = append(args, offset)
	}
	return base, args
}

// deleteByID executes a DELETE and returns an error if the row was not found.
func deleteByID(db *sql.DB, ctx context.Context, table string, id int, entityName string) error {
	res, err := db.ExecContext(ctx, "DELETE FROM "+table+" WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete %s: %w", entityName, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("%s not found", entityName)
	}
	return nil
}
