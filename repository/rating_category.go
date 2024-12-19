package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
)

type RatingCategoryRepository struct {
	Conn *sql.DB
}

func NewRatingCategoryRepository(conn *sql.DB) *RatingCategoryRepository {
	return &RatingCategoryRepository{conn}
}

func (repository *RatingCategoryRepository) FetchAll(ctx context.Context) (result []domain.RatingCategory, err error) {
	query := "SELECT id, name, weight FROM rating_categories"
	rows, err := repository.Conn.QueryContext(ctx, query)
	if err != nil {
		log.Println("error while querying rating_categories table", err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println("error trying to close rows", err)
		}
	}()

	result = make([]domain.RatingCategory, 0)
	for rows.Next() {
		ratingCategory := domain.RatingCategory{}
		err = rows.Scan(
			&ratingCategory.ID,
			&ratingCategory.Name,
			&ratingCategory.Weight,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, ratingCategory)
	}

	return result, nil

}
