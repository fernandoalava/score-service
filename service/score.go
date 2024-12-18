package service

import (
	"context"
	"fmt"
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

func (scoreService *ScoreService) GetScoreByTicket(ctx context.Context, from time.Time, to time.Time) (res []domain.ScoreByTicket, err error) {
	tickets, err := scoreService.ticketRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	ratingCategories, err := scoreService.ratingCategoryRepository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	ratingCategoriesMap := lo.Associate(ratingCategories, func(ratingCategory domain.RatingCategory) (uint64, domain.RatingCategory) {
		return ratingCategory.ID, ratingCategory
	})
	ratings, err := scoreService.ratingRepository.FindByCreatedAtBetween(ctx, from, to)
	if err != nil {
		return nil, err
	}
	ratingsMap := lo.Map(ratings, func(rating domain.Rating, _ int) domain.RatingWithCategory {
		return domain.RatingWithCategory{
			ID:             rating.ID,
			Rating:         rating.Rating,
			TicketID:       rating.TicketID,
			RatingCategory: ratingCategoriesMap[rating.RatingCategoryID],
			CreatedAt:      rating.CreatedAt,
		}
	})
	fmt.Println(tickets)
	fmt.Println(ratingCategories)
	fmt.Println(ratings)
	fmt.Println(ratingsMap)
	return
}
