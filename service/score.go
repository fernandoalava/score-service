package service

import (
	"context"
	"time"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
	"github.com/fernandoalava/softwareengineer-test-task/util"
	"github.com/samber/lo"
)

type RatingCategoryRepository interface {
	FetchAll(ctx context.Context) (response []domain.RatingCategory, err error)
}

type ScoreRepository interface {
	FetchScoreByTicketBetween(ctx context.Context, from time.Time, to time.Time) (response []domain.ScoreByTicket, err error)
	FetchAggregateScoreOverPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.ScoreByCategoryWithPeriod, error)
	FetchOverallQuality(ctx context.Context, from, to time.Time) (float64, error)
}

type ScoreService struct {
	ratingCategoryRepository RatingCategoryRepository
	scoreRepository          ScoreRepository
}

type PeriodScoreWithRatings struct {
	From    time.Time
	To      time.Time
	Score   float64
	Ratings uint32
}

type CategoryScoreOverTime struct {
	CategoryName            string
	PeriodScoresWithRatings []PeriodScoreWithRatings
	TotalScore              float64
	TotalRating             uint32
}

type RatingCategoryScore struct {
	RatingCategoryID   uint64
	RatingCategoryName string
	Score              float64
}

type TicketScoreByCategory struct {
	TicketID             uint64
	RatingCategoryScores []RatingCategoryScore
}

type PeriodScore struct {
	From  time.Time
	To    time.Time
	Score float64
}

type GetPeriodOverPeriodScoreChangeResponse struct {
	CurrentPeriod   PeriodScore
	PreviousPeriod  PeriodScore
	ScoreDifference float64
}

type RatingCategoryPeriodScore struct {
	RatingCategory domain.RatingCategory
	DateRange      util.DateRange
	Score          float64
	Rating         uint32
}

func NewScoreService(ratingCategoryRepository RatingCategoryRepository, scoreRepository ScoreRepository) *ScoreService {
	return &ScoreService{
		ratingCategoryRepository: ratingCategoryRepository,
		scoreRepository:          scoreRepository,
	}
}

func (scoreService *ScoreService) GetScoreByTicket(ctx context.Context, from time.Time, to time.Time) ([]TicketScoreByCategory, error) {

	err := util.ValidateTimeRange(from, to)
	if err != nil {
		return nil, err
	}

	categoryScoresByTicket, err := scoreService.scoreRepository.FetchScoreByTicketBetween(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return lo.MapToSlice(lo.GroupBy(categoryScoresByTicket, func(scoreByTicket domain.ScoreByTicket) uint64 {
		return scoreByTicket.TicketID
	}), func(ticketId uint64, scoresByTicket []domain.ScoreByTicket) TicketScoreByCategory {
		return TicketScoreByCategory{
			TicketID: ticketId,
			RatingCategoryScores: lo.Map(scoresByTicket, func(scoreByTicket domain.ScoreByTicket, _ int) RatingCategoryScore {
				return RatingCategoryScore{
					RatingCategoryID:   scoreByTicket.CategoryID,
					RatingCategoryName: scoreByTicket.CategoryName,
					Score:              util.FormatScore(scoreByTicket.Score),
				}
			}),
		}
	}), nil

}

func (scoreService *ScoreService) GetAggregatedCategoryScoresOverTime(ctx context.Context, from time.Time, to time.Time) ([]CategoryScoreOverTime, error) {
	err := util.ValidateTimeRange(from, to)
	if err != nil {
		return nil, err
	}

	rangeOfDates := util.GenerateDateRanges(from, to)

	categories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)

	if err != nil {
		return nil, err
	}

	aggregateScoreOverPeriod, err := scoreService.scoreRepository.FetchAggregateScoreOverPeriod(ctx, from, to)
	if err != nil {
		return nil, err
	}

	aggregateScoreOverPeriodGroupedByCategory := lo.MapValues(lo.GroupBy(aggregateScoreOverPeriod, func(score domain.ScoreByCategoryWithPeriod) string { return score.CategoryName }), func(scores []domain.ScoreByCategoryWithPeriod, _ string) map[util.DateRange][]domain.ScoreByCategoryWithPeriod {
		return lo.GroupBy(scores, func(score domain.ScoreByCategoryWithPeriod) util.DateRange {
			return score.AggregationPeriod
		})
	})

	return lo.Map(categories, func(category domain.RatingCategory, _ int) CategoryScoreOverTime {
		groupedByRange, exists := aggregateScoreOverPeriodGroupedByCategory[category.Name]
		scoresWithRating := lo.FlatMap(rangeOfDates, func(currentRange util.DateRange, _ int) []PeriodScoreWithRatings {
			existingScoreInRange, exists := groupedByRange[currentRange]
			if !exists {
				return []PeriodScoreWithRatings{{
					From:    currentRange.From,
					To:      currentRange.To,
					Score:   float64(0),
					Ratings: uint32(0),
				}}
			}
			return lo.Map(existingScoreInRange, func(score domain.ScoreByCategoryWithPeriod, _ int) PeriodScoreWithRatings {
				return PeriodScoreWithRatings{
					From:    score.AggregationPeriod.From,
					To:      score.AggregationPeriod.To,
					Score:   util.FormatScore(score.CategoryScore),
					Ratings: uint32(score.RatingsCount),
				}

			})
		})
		if !exists {
			return CategoryScoreOverTime{
				CategoryName:            category.Name,
				PeriodScoresWithRatings: scoresWithRating,
				TotalRating:             0,
				TotalScore:              0,
			}
		}
		totalRating := lo.SumBy(scoresWithRating, func(period PeriodScoreWithRatings) uint32 {
			return period.Ratings
		})
		totalScore := lo.SumBy(scoresWithRating, func(period PeriodScoreWithRatings) float64 {
			return period.Score
		}) / float64(len(scoresWithRating))
		return CategoryScoreOverTime{
			CategoryName:            category.Name,
			PeriodScoresWithRatings: scoresWithRating,
			TotalRating:             totalRating,
			TotalScore:              util.FormatScore(totalScore),
		}
	}), nil

}

func (scoreService *ScoreService) GetOverAllQualityScore(ctx context.Context, from time.Time, to time.Time) (float64, error) {
	err := util.ValidateTimeRange(from, to)
	if err != nil {
		return 0, err
	}

	score, err := scoreService.scoreRepository.FetchOverallQuality(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return util.FormatScore(score), nil
}

func (scoreService *ScoreService) GetPeriodOverPeriodScoreChange(ctx context.Context, from time.Time, to time.Time) (*GetPeriodOverPeriodScoreChangeResponse, error) {
	err := util.ValidateTimeRange(from, to)
	if err != nil {
		return nil, err
	}

	previousFrom, previousTo := util.CalculatePreviousPeriod(from, to)

	overAllQualityScoreCurrentPeriod, err := scoreService.GetOverAllQualityScore(ctx, from, to)
	if err != nil {
		return nil, err
	}
	overAllQualityScorePreviousPeriod, err := scoreService.GetOverAllQualityScore(ctx, previousFrom, previousTo)
	if err != nil {
		return nil, err
	}
	periodScoreCurrentPeriod := PeriodScore{From: from, To: to, Score: overAllQualityScoreCurrentPeriod}
	periodScorePreviousPeriod := PeriodScore{From: from, To: to, Score: overAllQualityScorePreviousPeriod}
	getPeriodOverPeriodScoreChangeResponse := &GetPeriodOverPeriodScoreChangeResponse{
		CurrentPeriod:   periodScoreCurrentPeriod,
		PreviousPeriod:  periodScorePreviousPeriod,
		ScoreDifference: util.FormatScore((overAllQualityScoreCurrentPeriod - overAllQualityScorePreviousPeriod) / overAllQualityScorePreviousPeriod),
	}

	return getPeriodOverPeriodScoreChangeResponse, nil

}
