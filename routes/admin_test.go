package routes

import "testing"

func TestAdminURLs(t *testing.T) {
	tests := []struct {
		name, got, want string
	}{
		{"home", AdminHomePage.URL(), "/admin"},
		{"work", AdminWorkEdit.URL(12), "/admin/works/12/edit"},
		{"project", AdminProjectShow.URL(12), "/admin/projects/12"},
		{"post", AdminBlogPostNew.URL(), "/admin/blog-posts/new"},
		{"ban", AdminProjectInquiryBanDomain.URL(12), "/admin/project-inquiries/12/domain-ban"},
		{"today", AdminDiaryEntryToday.URL(), "/admin/diary/today"},
		{"queue", AdminQueueJobRetry.URL(123), "/admin/queue/jobs/123/retry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("URL = %q, want %q", tt.got, tt.want)
			}
		})
	}
}
