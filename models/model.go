// Package models contains data models and validation logic.
package models

import (
	"github.com/go-playground/validator/v10"
)

var Validate = setupValidator()

func setupValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	return v
}

type (
	user           struct{}
	token          struct{}
	project        struct{}
	blogPost       struct{}
	diaryEntry     struct{}
	projectInquiry struct{}
	work           struct{}
)

var (
	User           user
	Token          token
	Project        project
	BlogPost       blogPost
	DiaryEntry     diaryEntry
	ProjectInquiry projectInquiry
	Work           work
)
