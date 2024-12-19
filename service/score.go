package service

import (
	"context"
	"time"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
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

func NewScoreService(ticketRepository TicketRepository, ratingCategoryRepository RatingCategoryRepository, ratingRepository RatingRepository) *ScoreService {
	return &ScoreService{
		ticketRepository:         ticketRepository,
		ratingCategoryRepository: ratingCategoryRepository,
		ratingRepository:         ratingRepository,
	}
}

func calculateScore(rating float32, weight float32) float32 {
	return (rating * weight * 1 / 5 * weight) * 100
}

type RatingCategoryScore struct {
	RatingCategoryID   uint64
	RatingCategoryName string
	Score              float32
}

type TicketScoreByCategory struct {
	TicketID             uint64
	RatingCategoryScores []RatingCategoryScore
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

	ticketsWithRatingCategoryScore := lo.MapValues(ratingsWithCategoryGroupedByTicketId, func(r []domain.RatingWithCategory, _ uint64) map[domain.RatingCategory]float32 {
		ratingsGroupedByCategory := lo.GroupBy(r, func(rr domain.RatingWithCategory) domain.RatingCategory {
			return rr.RatingCategory
		})
		ratingsWithScore := lo.MapValues(ratingsGroupedByCategory, func(ratingsWithCategory []domain.RatingWithCategory, _ domain.RatingCategory) float32 {
			scores := lo.Map(ratingsWithCategory, func(ratingWithCategory domain.RatingWithCategory, _ int) float32 {
				return calculateScore(float32(ratingWithCategory.Rating), float32(ratingWithCategory.RatingCategory.Weight))
			})
			return lo.Sum(scores) / float32(len(scores))
		})
		return ratingsWithScore
	})

	scoresByTicket = lo.Map(tickets, func(ticket domain.Ticket, _ int) TicketScoreByCategory {
		ratingCategoryScores := lo.Map(ratingCategories, func(ratingCategory domain.RatingCategory, _ int) RatingCategoryScore {
			score, exists := ticketsWithRatingCategoryScore[ticket.ID][ratingCategory]
			if !exists {
				score = float32(0)
			}
			return RatingCategoryScore{
				RatingCategoryID:   ratingCategory.ID,
				RatingCategoryName: ratingCategory.Name,
				Score:              score,
			}
		})
		return TicketScoreByCategory{
			TicketID:             ticket.ID,
			RatingCategoryScores: ratingCategoryScores,
		}
	})

	return
}

func (scoreService *ScoreService) GetOverAllQualityScore(ctx context.Context, from time.Time, to time.Time) (overAllScore float32, err error) {
	ratingCategories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)
	if err != nil {
		return 0, err
	}

	ratings, err := scoreService.ratingRepository.FindByCreatedAtBetween(ctx, from, to)
	if err != nil {
		return 0, err
	}

	ratingsWithRatingCategory := getRatingsWithRatingCategory(ratings, ratingCategories)
	ratingsWithRatingCategoryAndScores := lo.Map(ratingsWithRatingCategory, func(ratingWithCategory domain.RatingWithCategory, _ int) float32 {
		return calculateScore(float32(ratingWithCategory.Rating), float32(ratingWithCategory.RatingCategory.Weight))
	})

	overAllScore = lo.Sum(ratingsWithRatingCategoryAndScores) / float32(len(ratingsWithRatingCategoryAndScores))

	return
}
