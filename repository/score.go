package repository

import (
	"context"
	"database/sql"
	"github.com/fernandoalava/softwareengineer-test-task/domain"
	"time"
)

type ScoreRepository struct {
	Conn *sql.DB
}

func NewScoreRepository(conn *sql.DB) *ScoreRepository {
	return &ScoreRepository{conn}
}

func (repository *ScoreRepository) FetchAggregateScoreByTicketInRange(ctx context.Context, from time.Time, to time.Time) (result []domain.ScoreByTicket, err error) {

	fromStr := from.Format("2006-01-02T15:04:05")
	toStr := to.Format("2006-01-02T15:04:05")

	query := `WITH ratings_with_score AS(
	    SELECT r.ticket_id,
            r.rating_category_id AS category_id,
            rc.name AS category_name,
            COALESCE(ROUND((r.rating * rc.weight * 1.0) / (5 * rc.weight) * 100, 2), 0) AS score
		FROM ratings r
		INNER JOIN rating_categories rc ON r.rating_category_id = rc.id
		WHERE r.created_at BETWEEN ? AND ?)
		SELECT ticket_id,
		category_name,
		score
		FROM ratings_with_score;`

	rows, err := repository.Conn.QueryContext(ctx, query, fromStr, toStr)
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

	result = make([]domain.ScoreByTicket, 0)
	for rows.Next() {
		rating := domain.ScoreByTicket{}
		err = rows.Scan(
			&rating.TicketID,
			&rating.CategoryName,
			&rating.Score,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, rating)
	}

	return result, nil

}
