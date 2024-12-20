package domain

type ScoreByTicket struct {
	TicketID     uint64
	CategoryID   uint64
	CategoryName string
	Score        float64
}
