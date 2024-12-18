package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
)

type RatingRepository struct {
	Conn *sql.DB
}

func NewRatingRepository(conn *sql.DB) *RatingRepository {
	return &RatingRepository{conn}
}

func (repository *RatingRepository) FetchAllInRange(ctx context.Context, from time.Time, to time.Time) (result []domain.Rating, err error) {
	query := "SELECT id, rating, ticket_id, rating_category_id, created_at weight FROM rating WHERE created_at BETWEEN ? AND  ?"
	rows, err := repository.Conn.QueryContext(ctx, query)
	if err != nil {
		// TODO: use some logs
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			// TODO: use some logs
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
