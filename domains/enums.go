package domains

import "errors"

type StatusEnum string

const (
	Draft     StatusEnum = "draft"
	Published StatusEnum = "published"
)

func (s StatusEnum) IsValid() bool  { return s == Draft || s == Published }
func (s StatusEnum) String() string { return string(s) }
func (s StatusEnum) HumanString() string {
	switch s {
	case Draft:
		return "Draft"
	case Published:
		return "Published"
	default:
		return "Invalid Value"
	}
}
func ParseStatusEnum(s string) (StatusEnum, error) {
	status := StatusEnum(s)
	if !status.IsValid() {
		return "", errors.New("invalid status type")
	}
	return status, nil
}
