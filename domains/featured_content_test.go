package domains_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"mbvlabs/domains"
	"mbvlabs/services"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mbvlabs/andurel/pkg/storage"
)

type featuredLimitRow struct{}

func (featuredLimitRow) Scan(dest ...any) error { *dest[0].(*int64) = 3; return nil }

type featuredLimitTx struct {
	query                 string
	lock                  string
	exceptID              int32
	rolledBack, committed bool
}

func (tx *featuredLimitTx) Exec(_ context.Context, query string, _ ...any) (pgconn.CommandTag, error) {
	tx.lock = query
	return pgconn.CommandTag{}, nil
}
func (*featuredLimitTx) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected query")
}
func (tx *featuredLimitTx) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	tx.query, tx.exceptID = query, args[0].(int32)
	return featuredLimitRow{}
}
func (tx *featuredLimitTx) Commit(context.Context) error   { tx.committed = true; return nil }
func (tx *featuredLimitTx) Rollback(context.Context) error { tx.rolledBack = true; return nil }

type featuredLimitDB struct {
	tx *featuredLimitTx
}

func (*featuredLimitDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	panic("unexpected non-transactional write")
}
func (*featuredLimitDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	panic("unexpected non-transactional query")
}
func (*featuredLimitDB) QueryRow(context.Context, string, ...any) pgx.Row {
	panic("unexpected non-transactional query")
}
func (*featuredLimitDB) Health(context.Context) error { return nil }
func (db *featuredLimitDB) BeginTransaction(context.Context) (storage.Transaction, error) {
	return db.tx, nil
}

func TestFeaturedLimitsLockTheSetInTransaction(t *testing.T) {
	for _, tc := range []struct {
		name, table string
		id          int32
		want        error
		run         func(context.Context, *featuredLimitDB) error
	}{
		{"create work", "works", 0, services.ErrFeaturedWorkLimit, func(ctx context.Context, db *featuredLimitDB) error {
			_, err := services.NewWorks(domains.NewWorks(db), storage.NewTransactor(db)).CreateWork(ctx, domains.CreateWorkData{IsFeatured: true})
			return err
		}},
		{"update work", "works", 12, services.ErrFeaturedWorkLimit, func(ctx context.Context, db *featuredLimitDB) error {
			_, err := services.NewWorks(domains.NewWorks(db), storage.NewTransactor(db)).UpdateWork(ctx, domains.UpdateWorkData{ID: 12, IsFeatured: true})
			return err
		}},
		{"create project", "projects", 0, services.ErrFeaturedProjectLimit, func(ctx context.Context, db *featuredLimitDB) error {
			_, err := services.NewProjects(domains.NewProjects(db), storage.NewTransactor(db)).CreateProject(ctx, domains.CreateProjectData{IsFeatured: true})
			return err
		}},
		{"update project", "projects", 13, services.ErrFeaturedProjectLimit, func(ctx context.Context, db *featuredLimitDB) error {
			_, err := services.NewProjects(domains.NewProjects(db), storage.NewTransactor(db)).UpdateProject(ctx, domains.UpdateProjectData{ID: 13, IsFeatured: true})
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tx := &featuredLimitTx{}
			db := &featuredLimitDB{tx: tx}
			if err := tc.run(context.Background(), db); !errors.Is(err, tc.want) {
				t.Fatalf("error = %v, want %v", err, tc.want)
			}
			if !strings.Contains(tx.lock, "pg_advisory_xact_lock") || !strings.Contains(tx.lock, "featured-"+tc.table) ||
				!strings.Contains(tx.query, "FROM "+tc.table) || tx.exceptID != tc.id || !tx.rolledBack || tx.committed {
				t.Fatalf("transaction = %+v", tx)
			}
		})
	}
}
