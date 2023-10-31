package rabbitmq

import "time"

type NotifyMessage struct {
	ID        int64
	Title     string
	OwnerID   int64
	StartDate time.Time
}
