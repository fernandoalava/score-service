package domain

import (
	"github.com/fernandoalava/softwareengineer-test-task/util"
)

type ScoreByTicket struct {
	TicketID     uint64
	CategoryID   uint64
	CategoryName string
	Score        float64
}

type ScoreByCategoryWithPeriod struct {
	CategoryID        uint64
	CategoryName      string
	AggregationPeriod util.DateRange
	CategoryScore     float64
	RatingsCount      int
}
