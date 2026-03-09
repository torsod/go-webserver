package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// updateFields dynamically builds and executes an UPDATE statement for the given fields
func updateFields(ctx context.Context, pool *pgxpool.Pool, table string, id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	var setClauses []string
	var args []interface{}
	i := 1

	for col, val := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", table, strings.Join(setClauses, ", "), i)

	_, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update %s fields: %w", table, err)
	}
	return nil
}
