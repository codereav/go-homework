package storage

import "time"

type Event struct {
	ID         int64
	Title      string
	Descr      string
	OwnerID    int64
	StartDate  *time.Time
	EndDate    *time.Time
	RemindFor  *time.Time
	NotifySent bool
	DeletedAt  *time.Time
}
