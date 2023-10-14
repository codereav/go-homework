package storage

import "time"

type Event struct {
	ID        int
	Title     string
	Descr     string
	OwnerID   int
	StartDate time.Time
	EndDate   time.Time
	RemindFor time.Duration
	DeletedAt *time.Time
}
