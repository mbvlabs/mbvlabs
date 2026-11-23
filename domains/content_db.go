package domains

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// OptionalString and OptionalTime preserve the legacy API representation of
// nullable content fields while keeping SQL imports in the domains layer.
type OptionalString = sql.NullString
type OptionalTime = sql.NullTime

// contentDB is the query surface shared by the connection and a transaction.
type contentDB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// LockFeatured serializes changes to the featured-work set for the lifetime
// of the transaction. Call it on a transaction-bound receiver before counting.
func (w *Works) LockFeatured(ctx context.Context) error {
	_, err := w.db.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('mbvlabs'), hashtext('featured-works'))`)
	return err
}

// LockFeatured serializes changes to the featured-project set.
func (p *Projects) LockFeatured(ctx context.Context) error {
	_, err := p.db.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext('mbvlabs'), hashtext('featured-projects'))`)
	return err
}
