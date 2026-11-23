package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/anduin"
	"github.com/mbvlabs/andurel/pkg/inertia"
	"github.com/mbvlabs/andurel/pkg/kiks"
)

func adminContentID(ctx context.Context) (int32, error) {
	id, err := strconv.ParseInt(anduin.PathValue(ctx, "id"), 10, 32)
	if err != nil {
		return 0, anduin.BadRequest(err)
	}
	return int32(id), nil
}

func adminContentPagination(ctx context.Context) (int64, int64) {
	q := anduin.Request(ctx).URL.Query()
	page, err := strconv.ParseInt(q.Get("page"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.ParseInt(q.Get("per_page"), 10, 64)
	if err != nil || size < 1 || size > 100 {
		size = 25
	}
	return page, size
}

func adminContentDate(value domains.OptionalTime) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format("2006-01-02")
}

func adminContentParseDate(value, field string, errs map[string]string) time.Time {
	if value == "" {
		return time.Time{}
	}
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		errs[field] = "must be a valid date"
	}
	return date
}

func adminContentFlashRedirect(ctx context.Context, renderer *inertia.Renderer, level kiks.FlashType, message, target string) error {
	if err := kiks.AddFlash(ctx, level, message); err != nil {
		return err
	}
	return renderer.Redirect(ctx, target, http.StatusSeeOther)
}
