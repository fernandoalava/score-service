package repository

import (
	"context"
	"database/sql"
	"time"

	"log"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
	"github.com/fernandoalava/softwareengineer-test-task/util"
)

type ScoreRepository struct {
	Conn *sql.DB
}

func NewScoreRepository(conn *sql.DB) *ScoreRepository {
	return &ScoreRepository{conn}
}

func (repository *ScoreRepository) FetchScoreByTicketBetween(ctx context.Context, from time.Time, to time.Time) ([]domain.ScoreByTicket, error) {
	query := `
		WITH FilteredRatings AS (
		SELECT
			t.id as ticket_id,
			r.rating_category_id,
			c.name as rating_category_name,
			r.rating,
			c.weight
		FROM
			ratings r
		JOIN
			rating_categories c ON r.rating_category_id = c.id
		JOIN
			tickets t ON r.ticket_id = t.id
		WHERE
			r.created_at BETWEEN ? AND ?
	),
	WeightedAverages AS (
		SELECT
			ticket_id,
			rating_category_id,
			rating_category_name,
			SUM(rating * weight) / SUM(weight) AS weighted_average
		FROM
			FilteredRatings
		GROUP BY
			ticket_id,
			rating_category_id,
			rating_category_name
	)
	SELECT
		ticket_id,
		rating_category_id,
		rating_category_name,
		COALESCE(ROUND(weighted_average / 5 * 100, 2),0) AS category_score
	FROM
		WeightedAverages;
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

func (repository *ScoreRepository) FetchAggregateScoreOverPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.ScoreByCategoryWithPeriod, error) {
	query := `
		WITH FilteredRatings AS (
		SELECT
			t.id as ticket_id,
			r.rating_category_id,
			c.name as rating_category_name,
			r.rating,
			c.weight,
			r.created_at as review_date
		FROM
			ratings r
		JOIN
			rating_categories c ON r.rating_category_id = c.id
		JOIN
			tickets t ON r.ticket_id = t.id
		WHERE
			r.created_at BETWEEN  ? AND ?
	),
	DailyAverages AS (
		SELECT
			rating_category_id,
			rating_category_name,
			date(review_date) AS review_date,
			AVG(rating) AS daily_average_rating,
			COUNT(*) AS daily_rating_count
		FROM
			FilteredRatings
		GROUP BY
			rating_category_id,
			rating_category_name,
			date(review_date)
	)
	,
	WeeklyAverages AS (
		SELECT
			rating_category_id,
			rating_category_name,
			date(review_date, 'weekday 1') AS week_start,
			date(review_date, 'weekday 1', '+6 days') AS week_end,
			AVG(daily_average_rating) AS weekly_average_rating,
			SUM(daily_rating_count) AS weekly_rating_count
		FROM
			DailyAverages
		GROUP BY
			rating_category_id,
			rating_category_name,
			date(review_date, 'weekday 0'),
			date(review_date, 'weekday 1', '+6 days') 
	)
	,
	AggregatedScores AS (
		SELECT
			DailyAverages.rating_category_id,
			DailyAverages.rating_category_name,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN review_date 
				ELSE week_start || '/' || week_end 
			END AS aggregation_period,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN daily_average_rating 
				ELSE weekly_average_rating 
			END AS average_rating,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN daily_rating_count 
				ELSE weekly_rating_count 
			END AS rating_count
		FROM
			DailyAverages
		LEFT JOIN
			WeeklyAverages ON DailyAverages.rating_category_id = WeeklyAverages.rating_category_id
				AND DailyAverages.review_date = WeeklyAverages.week_start
		UNION ALL
		SELECT
			WeeklyAverages.rating_category_id,
			WeeklyAverages.rating_category_name,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN review_date 
				ELSE week_start || '/' || week_end 
			END AS aggregation_period,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN daily_average_rating 
				ELSE weekly_average_rating 
			END AS average_rating,
			CASE 
				WHEN (julianday(?) - julianday(?)) <= 31 THEN daily_rating_count 
				ELSE weekly_rating_count 
			END AS rating_count
		FROM
			WeeklyAverages
		LEFT JOIN
			DailyAverages ON WeeklyAverages.rating_category_id = DailyAverages.rating_category_id
				AND WeeklyAverages.week_start = DailyAverages.review_date
		WHERE
			DailyAverages.review_date IS NULL
	)
	SELECT
		rating_category_id,
		rating_category_name,
		aggregation_period,
		ROUND(average_rating / 5 * 100, 2) AS category_score,
		rating_count
	FROM
		AggregatedScores
	WHERE aggregation_period IS NOT NULL
	ORDER BY
		rating_category_id,
		rating_category_name,
		aggregation_period;
	`
	fromStringValue := util.TimeToString(from)
	toStringValue := util.TimeToString(to)
	rows, err := repository.Conn.QueryContext(ctx, query, fromStringValue, toStringValue, toStringValue, fromStringValue, toStringValue, fromStringValue, toStringValue, fromStringValue, toStringValue, fromStringValue, toStringValue, fromStringValue, toStringValue, fromStringValue)
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
	var result []domain.ScoreByCategoryWithPeriod
	for rows.Next() {
		scoreByCategoryWithPeriod := domain.ScoreByCategoryWithPeriod{}
		var aggregatePeriod string
		err = rows.Scan(
			&scoreByCategoryWithPeriod.CategoryID,
			&scoreByCategoryWithPeriod.CategoryName,
			&aggregatePeriod,
			&scoreByCategoryWithPeriod.CategoryScore,
			&scoreByCategoryWithPeriod.RatingsCount,
		)

		period, err := util.ParsePeriodFromString(aggregatePeriod, "/")
		if err != nil {
			return nil, err
		}
		scoreByCategoryWithPeriod.AggregationPeriod = *period
		result = append(result, scoreByCategoryWithPeriod)
	}

	return result, nil
}

func (repository *ScoreRepository) FetchOverallQuality(ctx context.Context, from, to time.Time) (float64, error) {

	query := `
		WITH FilteredRatings AS (
			SELECT
				r.rating,
				c.weight
			FROM
				ratings r
			JOIN
				rating_categories c ON r.rating_category_id = c.id
			WHERE
				r.created_at BETWEEN ? AND ?
		),
		WeightedAverage AS (
			SELECT 
				SUM(rating * weight) / SUM(weight) AS overall_average_rating
			FROM 
				FilteredRatings r
		)
		SELECT 
			ROUND(AVG(overall_average_rating) / 5 * 100, 2) AS overall_score 
		FROM 
			WeightedAverage;
	`

	fromStringValue := util.TimeToString(from)
	toStringValue := util.TimeToString(to)
	rows, err := repository.Conn.QueryContext(ctx, query, fromStringValue, toStringValue)
	if err != nil {
		log.Println("error while querying ratings table", err)
		return 0, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println("error trying to close rows", err)
		}
	}()
	var overallScore float32
	for rows.Next() {
		err = rows.Scan(
			&overallScore,
		)
		if err != nil {
			return 0, err
		}
	}

	return float64(overallScore), nil

}
