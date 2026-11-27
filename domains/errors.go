// Package domains owns business rules and persistence. Domain writes validate
// before querying; services coordinate work spanning multiple receivers.
//
// A receiver holds a storage.Transactor when one of its own methods must be
// atomic by itself, either because it writes more than once or because it
// reads, modifies and writes a row.
package domains

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/mbvlabs/andurel/pkg/storage"
)

var (
	ErrNotFound         = pgx.ErrNoRows
	ErrConflict         = storage.ErrConflict
	ErrDomainValidation = errors.New("the provided payload failed validations")
)

func IsDuplicateWorkSlug(err error) bool {
	return errors.Is(storage.Conflict(err, "works_slug_key"), ErrConflict)
}

func IsDuplicateBlogPostSlug(err error) bool {
	return errors.Is(storage.Conflict(err, "blog_posts_slug_key"), ErrConflict)
}
