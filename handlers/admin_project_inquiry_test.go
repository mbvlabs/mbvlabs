package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

type adminInquiryStub struct {
	entry        domains.ProjectInquiryEntity
	banned       bool
	checkedEmail string
}

func (s *adminInquiryStub) Find(context.Context, int32) (domains.ProjectInquiryEntity, error) {
	return s.entry, nil
}
func (s *adminInquiryStub) Paginate(context.Context, int64, int64) (domains.PaginatedProjectInquiries, error) {
	return domains.PaginatedProjectInquiries{}, nil
}
func (s *adminInquiryStub) Create(context.Context, domains.CreateProjectInquiryData) (domains.ProjectInquiryEntity, error) {
	return s.entry, nil
}
func (s *adminInquiryStub) Update(context.Context, domains.UpdateProjectInquiryData) (domains.ProjectInquiryEntity, error) {
	return s.entry, nil
}
func (s *adminInquiryStub) Destroy(context.Context, int32) error { return nil }
func (s *adminInquiryStub) IsDomainBanned(_ context.Context, email string) (bool, error) {
	s.checkedEmail = email
	return s.banned, nil
}
func (s *adminInquiryStub) BanDomain(context.Context, int32) (string, error) {
	return "example.com", nil
}
func (s *adminInquiryStub) UnbanDomain(context.Context, int32) (string, error) {
	return "example.com", nil
}

func TestAdminInquiryShowIncludesDomainBan(t *testing.T) {
	store := &adminInquiryStub{entry: domains.ProjectInquiryEntity{ID: 8, Email: "a@example.com", Metadata: json.RawMessage(`{"key":"value"}`)}, banned: true}
	rt := anduin.NewRouter()
	(AdminProjectInquiryHandler{inquiries: store, renderer: adminTestRenderer(t)}).Register(rt)
	req := httptest.NewRequest(http.MethodGet, "/admin/project-inquiries/8", nil)
	req.Header.Set("X-Inertia", "true")
	w := httptest.NewRecorder()
	rt.Handler().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var page struct {
		Component string `json:"component"`
		Props     struct {
			Item struct {
				DomainBanned bool
				Metadata     string
			} `json:"item"`
		} `json:"props"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if page.Component != "Admin/ProjectInquiry/Show" || !page.Props.Item.DomainBanned || store.checkedEmail != "a@example.com" || page.Props.Item.Metadata == "" {
		t.Fatalf("unexpected page: %+v", page)
	}
}
