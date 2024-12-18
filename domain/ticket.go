package domain

import "time"

type Ticket struct {
	ID        uint64
	Subject   string
	CreatedAt time.Time
}
