package repository

import (
	"context"
	"database/sql"
	"time"

	"log"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
	"github.com/fernandoalava/softwareengineer-test-task/util"
)

type RatingRepository struct {
	Conn *sql.DB
}

func NewRatingRepository(conn *sql.DB) *RatingRepository {
	return &RatingRepository{conn}
}

func (repository *RatingRepository) FindByCreatedAtBetween(ctx context.Context, from time.Time, to time.Time) (result []domain.Rating, err error) {
	query := "SELECT id, rating, ticket_id, rating_category_id, created_at weight FROM ratings WHERE created_at BETWEEN ? AND  ?"
	rows, err := repository.Conn.QueryContext(ctx, query, util.TimeToString(from), util.TimeToString(to))
	if err != nil {
		log.Println("error while querying ratings table", err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println("error trying to close rows", err)
		}
	}()

	result = make([]domain.Rating, 0)
	for rows.Next() {
		rating := domain.Rating{}
		err = rows.Scan(
			&rating.ID,
			&rating.Rating,
			&rating.TicketID,
			&rating.RatingCategoryID,
			&rating.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, rating)
	}

	return result, nil

}
