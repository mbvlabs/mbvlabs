package services

import (
	"context"
	"testing"
	"time"

	"mbvlabs/domains"

	"github.com/mbvlabs/andurel/pkg/validation"
)

func TestBlogPostScheduleRejectsPastBeforeWriting(t *testing.T) {
	at := time.Now().Add(-time.Hour)
	_, err := (BlogPosts{}).Create(context.Background(), domains.CreateBlogPostData{}, &at)
	if _, ok := validation.As(err); !ok {
		t.Fatalf("expected field validation error, got %v", err)
	}
}
