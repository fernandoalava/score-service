package domain

import "time"

type Rating struct {
	ID uint64
	Rating rune
	TicketID uint64
	RatingCategoryID uint64
	CreatedAt time.Time
}

type RatingWithCategory struct {
	ID uint64
	Rating rune
	TicketID uint64
	RatingCategory RatingCategory
	CreatedAt time.Time
}