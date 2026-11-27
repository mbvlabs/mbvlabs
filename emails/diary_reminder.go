package emails

import "github.com/mbvlabs/andurel/pkg/postmand"

type DiaryReminder struct {
	Title    string
	Period   string
	DiaryURL string
}

var _ postmand.Transformer = DiaryReminder{}

func (r DiaryReminder) ToHTML() (string, error) {
	return templates.HTML("diary_reminder", r)
}

func (r DiaryReminder) ToText() (string, error) {
	return templates.Text("diary_reminder", r)
}
