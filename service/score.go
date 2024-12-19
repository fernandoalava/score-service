package service

import (
	"context"
	"math"
	"time"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
	"github.com/fernandoalava/softwareengineer-test-task/util"
	"github.com/samber/lo"
)

type TicketRepository interface {
	FetchAll(ctx context.Context) (response []domain.Ticket, err error)
}

type RatingCategoryRepository interface {
	FetchAll(ctx context.Context) (response []domain.RatingCategory, err error)
}

type RatingRepository interface {
	FindByCreatedAtBetween(ctx context.Context, from time.Time, to time.Time) (response []domain.Rating, err error)
}

type ScoreService struct {
	ticketRepository         TicketRepository
	ratingCategoryRepository RatingCategoryRepository
	ratingRepository         RatingRepository
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

func NewScoreService(ticketRepository TicketRepository, ratingCategoryRepository RatingCategoryRepository, ratingRepository RatingRepository) *ScoreService {
	return &ScoreService{
		ticketRepository:         ticketRepository,
		ratingCategoryRepository: ratingCategoryRepository,
		ratingRepository:         ratingRepository,
	}
}

func calculateScore(rating int, weight float32) float64 {
	return util.FormatScore(float64((float32(rating) * weight * 1 / 5 * weight) * 100))
}

func getRatingsWithRatingCategory(ratings []domain.Rating, ratingCategories []domain.RatingCategory) []domain.RatingWithCategory {
	ratingCategoriesMap := lo.Associate(ratingCategories, func(ratingCategory domain.RatingCategory) (uint64, domain.RatingCategory) {
		return ratingCategory.ID, ratingCategory
	})
	ratingsWithCategories := lo.Map(ratings, func(rating domain.Rating, _ int) domain.RatingWithCategory {
		return domain.RatingWithCategory{
			ID:             rating.ID,
			Rating:         rating.Rating,
			TicketID:       rating.TicketID,
			RatingCategory: ratingCategoriesMap[rating.RatingCategoryID],
			CreatedAt:      rating.CreatedAt,
		}
	})
	return ratingsWithCategories
}

func (scoreService *ScoreService) GetScoreByTicket(ctx context.Context, from time.Time, to time.Time) (scoresByTicket []TicketScoreByCategory, err error) {

	err = util.ValidateTimeRange(from, to)
	if err != nil {
		return nil, err
	}

	tickets, err := scoreService.ticketRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	ratingCategories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	ratings, err := scoreService.ratingRepository.FindByCreatedAtBetween(ctx, from, to)
	if err != nil {
		return nil, err
	}

	ratingsWithCategories := getRatingsWithRatingCategory(ratings, ratingCategories)

	ratingsWithCategoryGroupedByTicketId := lo.GroupBy(ratingsWithCategories, func(ratingWithCategory domain.RatingWithCategory) uint64 {
		return ratingWithCategory.ID
	})

	ticketsWithRatingCategoryScore := lo.MapValues(ratingsWithCategoryGroupedByTicketId, func(r []domain.RatingWithCategory, _ uint64) map[domain.RatingCategory]float64 {
		ratingsGroupedByCategory := lo.GroupBy(r, func(rr domain.RatingWithCategory) domain.RatingCategory {
			return rr.RatingCategory
		})
		ratingsWithScore := lo.MapValues(ratingsGroupedByCategory, func(ratingsWithCategory []domain.RatingWithCategory, _ domain.RatingCategory) float64 {
			scores := lo.Map(ratingsWithCategory, func(ratingWithCategory domain.RatingWithCategory, _ int) float64 {
				return calculateScore(ratingWithCategory.Rating, ratingWithCategory.RatingCategory.Weight)
			})
			return lo.Sum(scores) / float64(len(scores))
		})
		return ratingsWithScore
	})

	scoresByTicket = lo.Map(tickets, func(ticket domain.Ticket, _ int) TicketScoreByCategory {
		ratingCategoryScores := lo.Map(ratingCategories, func(ratingCategory domain.RatingCategory, _ int) RatingCategoryScore {
			score, exists := ticketsWithRatingCategoryScore[ticket.ID][ratingCategory]
			if !exists {
				score = float64(0)
			}
			return RatingCategoryScore{
				RatingCategoryID:   ratingCategory.ID,
				RatingCategoryName: ratingCategory.Name,
				Score:              util.FormatScore(score),
			}
		})
		return TicketScoreByCategory{
			TicketID:             ticket.ID,
			RatingCategoryScores: ratingCategoryScores,
		}
	})

	return
}

func (scoreService *ScoreService) GetAggregatedCategoryScoresOverTime(ctx context.Context, from time.Time, to time.Time) (categoryScoresOverTime []CategoryScoreOverTime, err error) {
	err = util.ValidateTimeRange(from, to)
	if err != nil {
		return nil, err
	}

	ratingCategories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	ratings, err := scoreService.ratingRepository.FindByCreatedAtBetween(ctx, from, to)
	if err != nil {
		return nil, err
	}

	ratingsWithRatingCategory := getRatingsWithRatingCategory(ratings, ratingCategories)
	rangeOfDates := util.GenerateDateRanges(from, to)

	ratingsGroupedByDate := lo.GroupBy(ratingsWithRatingCategory, func(rating domain.RatingWithCategory) time.Time {
		return time.Date(rating.CreatedAt.Year(), rating.CreatedAt.Month(), rating.CreatedAt.Day(), 0, 0, 0, 0, rating.CreatedAt.Location())
	})

	ratingsGroupedByDateAndCategory := lo.MapValues(ratingsGroupedByDate, func(ratings []domain.RatingWithCategory, _ time.Time) map[uint64][]lo.Tuple3[uint64, string, float64] {
		ratingsGroupedByCategory := lo.GroupBy(ratings, func(rating domain.RatingWithCategory) uint64 {
			return rating.RatingCategory.ID
		})
		return lo.MapValues(ratingsGroupedByCategory, func(r []domain.RatingWithCategory, _ uint64) []lo.Tuple3[uint64, string, float64] {
			return lo.Map(r, func(ratingWithCategory domain.RatingWithCategory, _ int) lo.Tuple3[uint64, string, float64] {
				score := calculateScore(ratingWithCategory.Rating, ratingWithCategory.RatingCategory.Weight)
				return lo.T3(ratingWithCategory.RatingCategory.ID, ratingWithCategory.RatingCategory.Name, score)
			})
		})

	})

	scorePeriods := lo.FlatMap(ratingCategories, func(ratingCategory domain.RatingCategory, _ int) []RatingCategoryPeriodScore {

		return lo.FlatMap(rangeOfDates, func(dateRange util.DateRange, _ int) []RatingCategoryPeriodScore {
			ratingsGroupedByDateAndCategoryFiltered := lo.PickBy(ratingsGroupedByDateAndCategory, func(date time.Time, _ map[uint64][]lo.Tuple3[uint64, string, float64]) bool {
				return util.IsDateInRange(date, dateRange)
			})
			return lo.MapToSlice(lo.MapValues(ratingsGroupedByDateAndCategoryFiltered, func(ratingsGroupedByCategoryId map[uint64][]lo.Tuple3[uint64, string, float64], _ time.Time) RatingCategoryPeriodScore {
				ratings, exists := ratingsGroupedByCategoryId[ratingCategory.ID]

				if !exists {
					return RatingCategoryPeriodScore{
						RatingCategory: ratingCategory,
						DateRange:      dateRange,
						Score:          float64(0),
						Rating:         uint32(0),
					}
				}
				totalScore := lo.SumBy(ratings, func(rating lo.Tuple3[uint64, string, float64]) float64 {
					_, _, score := lo.Unpack3(rating)
					return score
				})
				return RatingCategoryPeriodScore{
					RatingCategory: ratingCategory,
					DateRange:      dateRange,
					Score:          util.FormatScore(totalScore / float64(len(ratings))),
					Rating:         uint32(len(ratings)),
				}
			}), func(date time.Time, ratingCategoryPeriodScore RatingCategoryPeriodScore) RatingCategoryPeriodScore {
				return ratingCategoryPeriodScore
			})
		})

	})

	scorePeriodsGroupedByRatingCategory := lo.GroupBy(scorePeriods, func(ratingCategoryPeriodScore RatingCategoryPeriodScore) domain.RatingCategory {
		return ratingCategoryPeriodScore.RatingCategory
	})

	aggregatedScorePeriodsGroupedByCategory := lo.MapValues(scorePeriodsGroupedByRatingCategory, func(scorePeriodWithCategory []RatingCategoryPeriodScore, key domain.RatingCategory) lo.Tuple3[[]PeriodScoreWithRatings, float64, uint32] {
		periods := lo.Map(scorePeriodWithCategory, func(ratingCategoryPeriodScore RatingCategoryPeriodScore, _ int) PeriodScoreWithRatings {
			return PeriodScoreWithRatings{
				From:    ratingCategoryPeriodScore.DateRange.From,
				To:      ratingCategoryPeriodScore.DateRange.To,
				Score:   util.FormatScore(ratingCategoryPeriodScore.Score),
				Ratings: ratingCategoryPeriodScore.Rating,
			}
		})
		scoresForPeriods := lo.Map(periods, func(periodScore PeriodScoreWithRatings, _ int) float64 {
			return periodScore.Score
		})
		ratingsForPeriod := lo.SumBy(periods, func(periodScore PeriodScoreWithRatings) uint32 {
			return periodScore.Ratings
		})

		return lo.T3(periods, lo.Sum(scoresForPeriods)/float64(len(scoresForPeriods)), ratingsForPeriod)
	})

	categoryScoresOverTime = lo.MapToSlice(aggregatedScorePeriodsGroupedByCategory, func(ratingCategory domain.RatingCategory, scores lo.Tuple3[[]PeriodScoreWithRatings, float64, uint32]) CategoryScoreOverTime {
		periodScores, totalScore, totalRatings := lo.Unpack3(scores)
		return CategoryScoreOverTime{
			CategoryName:            ratingCategory.Name,
			PeriodScoresWithRatings: periodScores,
			TotalRating:             totalRatings,
			TotalScore:              util.FormatScore(totalScore),
		}
	})

	return

}

func (scoreService *ScoreService) GetOverAllQualityScore(ctx context.Context, from time.Time, to time.Time) (overAllScore float64, err error) {
	err = util.ValidateTimeRange(from, to)
	if err != nil {
		return 0, err
	}

	ratingCategories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)
	if err != nil {
		return 0, err
	}

	ratings, err := scoreService.ratingRepository.FindByCreatedAtBetween(ctx, from, to)
	if err != nil {
		return 0, err
	}

	ratingsWithRatingCategory := getRatingsWithRatingCategory(ratings, ratingCategories)
	ratingsWithRatingCategoryAndScores := lo.Map(ratingsWithRatingCategory, func(ratingWithCategory domain.RatingWithCategory, _ int) float64 {
		return calculateScore(ratingWithCategory.Rating, ratingWithCategory.RatingCategory.Weight)
	})

	overAllScore = util.FormatScore(lo.Sum(ratingsWithRatingCategoryAndScores) / float64(len(ratingsWithRatingCategoryAndScores)))

	return
}

func (scoreService *ScoreService) GetPeriodOverPeriodScoreChange(ctx context.Context, from time.Time, to time.Time) (getPeriodOverPeriodScoreChangeResponse *GetPeriodOverPeriodScoreChangeResponse, err error) {
	err = util.ValidateTimeRange(from, to)
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
	getPeriodOverPeriodScoreChangeResponse = &GetPeriodOverPeriodScoreChangeResponse{
		CurrentPeriod:   periodScoreCurrentPeriod,
		PreviousPeriod:  periodScorePreviousPeriod,
		ScoreDifference: util.FormatScore(math.Abs(float64(overAllQualityScoreCurrentPeriod - overAllQualityScorePreviousPeriod))),
	}

	return

}
