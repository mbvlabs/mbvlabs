package controllers_test

import (
	"encoding/gob"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"mbvlabs/controllers"
	admincontrollers "mbvlabs/controllers/admin"
	"mbvlabs/database"
	"mbvlabs/models"
	"mbvlabs/models/factories"
	"mbvlabs/router/cookies"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/v5/session"
	"github.com/labstack/echo/v5"
)

func TestProjectInquiries_DomainBan(t *testing.T) {
	gob.Register(cookies.FlashMessage{})

	ctx := t.Context()
	db := testCluster.NewTestDB(t, database.Migrations, "migrations")
	inquiry, err := factories.CreateProjectInquiry(
		ctx,
		db.Executor(),
		factories.WithProjectInquiriesEmail("Sales@Example.com"),
	)
	if err != nil {
		t.Fatalf("create project inquiry: %v", err)
	}

	store := sessions.NewCookieStore([]byte("01234567890123456789012345678901"))
	invoke := func(handler echo.HandlerFunc, method, path string, values url.Values, id int32) *httptest.ResponseRecorder {
		t.Helper()
		var body *strings.Reader
		if values == nil {
			body = strings.NewReader("")
		} else {
			body = strings.NewReader(values.Encode())
		}
		req := httptest.NewRequest(method, path, body)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		rec := httptest.NewRecorder()
		etx := echo.New().NewContext(req, rec)
		if id != 0 {
			etx.SetPathValues(echo.PathValues{{
				Name:  "id",
				Value: strconv.FormatInt(int64(id), 10),
			}})
		}
		if err := session.Middleware(store)(handler)(etx); err != nil {
			t.Fatalf("%s %s: %v", method, path, err)
		}
		return rec
	}

	adminController := admincontrollers.NewProjectInquiries(db)
	if rec := invoke(adminController.BanDomain, http.MethodPost, "/admin/project-inquiries/1/domain-ban", nil, inquiry.ID); rec.Code != http.StatusSeeOther {
		t.Fatalf("ban domain status = %d, want %d", rec.Code, http.StatusSeeOther)
	}

	var domain string
	if err := db.Executor().NewSelect().
		Table("project_inquiry_domain_bans").
		Column("domain").
		Scan(ctx, &domain); err != nil {
		t.Fatalf("find domain ban: %v", err)
	}
	if domain != "example.com" {
		t.Fatalf("banned domain = %q, want example.com", domain)
	}

	publicController := controllers.NewProjectInquiries(db)
	submit := func(email string) {
		t.Helper()
		rec := invoke(publicController.Create, http.MethodPost, "/contact", url.Values{
			"name":    {"Prospect"},
			"email":   {email},
			"message": {"A valid project inquiry message."},
		}, 0)
		if rec.Code != http.StatusSeeOther {
			t.Fatalf("submit %s status = %d, want %d", email, rec.Code, http.StatusSeeOther)
		}
	}

	submit("another@EXAMPLE.com")
	submit("client@allowed.test")

	count, err := db.Executor().NewSelect().Model((*models.ProjectInquiryEntity)(nil)).Count(ctx)
	if err != nil {
		t.Fatalf("count inquiries: %v", err)
	}
	if count != 2 {
		t.Fatalf("inquiry count after blocked and allowed submissions = %d, want 2", count)
	}

	if rec := invoke(adminController.UnbanDomain, http.MethodDelete, "/admin/project-inquiries/1/domain-ban", nil, inquiry.ID); rec.Code != http.StatusSeeOther {
		t.Fatalf("unban domain status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	submit("welcome@example.com")

	count, err = db.Executor().NewSelect().Model((*models.ProjectInquiryEntity)(nil)).Count(ctx)
	if err != nil {
		t.Fatalf("count inquiries after unban: %v", err)
	}
	if count != 3 {
		t.Fatalf("inquiry count after unban = %d, want 3", count)
	}
}
