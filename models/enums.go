package models

import "errors"

type StatusEnum string

const (
	Draft     StatusEnum = "draft"
	Published StatusEnum = "published"
)

func (c StatusEnum) IsValid() bool {
	return c == Draft || c == Published
}

func (c StatusEnum) String() string {
	return string(c)
}

func (c StatusEnum) HumanString() string {
	switch c {
	case Draft:
		return "Draft"
	case Published:
		return "Published"
	default:
		return "Invalid Value"
	}
}

func ParseStatusEnum(s string) (StatusEnum, error) {
	switch s {
	case Draft.String():
		return Draft, nil
	case Published.String():
		return Published, nil
	default:
		return "", errors.New("invalid status type")
	}
}
