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

func (repository *RatingRepository) FetchScoreByTicketBetween(ctx context.Context, from time.Time, to time.Time) ([]domain.ScoreByTicket, error) {
	query := `
		WITH normalized_ticket_rating AS (
			SELECT
			t.id,
			r.rating_category_id,
			rc.name as rating_category_name,
			((r.rating / 5.0) * 100) AS normalized_rating
			FROM tickets t
			LEFT JOIN ratings r ON (t.id = r.ticket_id)
			JOIN rating_categories rc ON (r.rating_category_id = rc.id)
			AND r.created_at BETWEEN ? AND ?
		)
		SELECT id, rating_category_id, rating_category_name, AVG(normalized_rating) score
		FROM normalized_ticket_rating
		GROUP BY id, rating_category_id, rating_category_name
		`
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

	var result []domain.ScoreByTicket
	for rows.Next() {
		scoreByTicket := domain.ScoreByTicket{}
		err = rows.Scan(
			&scoreByTicket.TicketID,
			&scoreByTicket.CategoryID,
			&scoreByTicket.CategoryName,
			&scoreByTicket.Score,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, scoreByTicket)
	}

	return result, nil
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
