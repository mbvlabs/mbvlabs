package jobs

type SendDiaryReminderArgs struct {
	Period string
}

func (SendDiaryReminderArgs) Kind() string { return "send_diary_reminder" }
