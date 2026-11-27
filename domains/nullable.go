package domains

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

// OptionalString and OptionalTime preserve the legacy API representation of
// nullable content fields while keeping SQL imports in the domains layer.
type OptionalString = sql.NullString
type OptionalTime = sql.NullTime

func contentText(v sql.NullString) pgtype.Text {
	return pgtype.Text{String: v.String, Valid: v.Valid}
}

func contentDate(v sql.NullTime) pgtype.Date {
	return pgtype.Date{Time: v.Time, Valid: v.Valid}
}

func contentTimestamp(v sql.NullTime) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: v.Time, Valid: v.Valid}
}
