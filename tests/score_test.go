package tests

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/fernandoalava/softwareengineer-test-task/repository"
	"github.com/fernandoalava/softwareengineer-test-task/service"
	"github.com/fernandoalava/softwareengineer-test-task/util"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func getScoreService() (*service.ScoreService, func()) {
	db, err := sql.Open("sqlite", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	closer := func() {
		err := db.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}

	ratingCategoryRepository := repository.NewRatingCategoryRepository(db)
	ratingRepository := repository.NewRatingRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	scoreService := service.NewScoreService(ticketRepository, ratingCategoryRepository, ratingRepository)
	return scoreService, closer
}

func TestGetScoreByTicket(t *testing.T) {
	scoreService, closer := getScoreService()
	defer closer()
	from, _ := util.StringToTime("2019-07-05T00:00:00")
	to, _ := util.StringToTime("2019-07-06T23:59:00")

	results, err := scoreService.GetScoreByTicket(context.TODO(), from, to)
	assert.Nil(t, err)
	assert.NotEmpty(t, results)
}

func TestGetOverAllQualityScore(t *testing.T) {
	scoreService, closer := getScoreService()
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	results, err := scoreService.GetOverAllQualityScore(context.TODO(), from, to)
	assert.Nil(t, err)
	assert.Equal(t, float64(36.61), results)
}

func TestGetAggregatedCategoryScoresOverTime(t *testing.T) {
	scoreService, closer := getScoreService()
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	results, err := scoreService.GetAggregatedCategoryScoresOverTime(context.TODO(), from, to)
	assert.Nil(t, err)
	assert.NotEmpty(t, results)
}

func TestGetPeriodOverPeriodScoreChange(t *testing.T) {
	scoreService, closer := getScoreService()
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	results, err := scoreService.GetPeriodOverPeriodScoreChange(context.TODO(), from, to)
	assert.Nil(t, err)
	assert.Equal(t, 0.68, results.ScoreDifference)
}
