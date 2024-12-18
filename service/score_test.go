package service_test

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

func TestCorruptDbErrors(t *testing.T) {
	db, err := sql.Open("sqlite", "../database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ratingCategoryRepository := repository.NewRatingCategoryRepository(db)
	ratingRepository := repository.NewRatingRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	scoreService := service.NewScoreService(ticketRepository,ratingCategoryRepository, ratingRepository)

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	results, err := scoreService.GetScoreByTicket(context.TODO(), from, to)
	assert.Nil(t, err)
	assert.Len(t, results, 108)
}
