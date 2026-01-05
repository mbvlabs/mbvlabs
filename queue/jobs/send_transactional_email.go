package jobs

import "mbvlabs/email"

type SendTransactionalEmailArgs struct {
	Data email.TransactionalData
}

func (SendTransactionalEmailArgs) Kind() string { return "send_transactional_email" }
