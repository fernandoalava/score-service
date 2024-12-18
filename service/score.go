package service

import (
	"context"
	"time"

	"github.com/fernandoalava/softwareengineer-test-task/domain"
)

type ScoreRepository interface {
	FetchAggregateScoreByTicketInRange(ctx context.Context, from time.Time, to time.Time) (response []domain.ScoreByTicket, err error)
}

type RatingCategory interface {
	FetchAllInRange(ctx context.Context, from time.Time, to time.Time) (result []domain.Rating, err error)
}

type ScoreService struct {
	scoreRepository ScoreRepository
}

func NewScoreService(scoreRepository ScoreRepository) *ScoreService {
	return &ScoreService{
		scoreRepository: scoreRepository,
	}
}

func (scoreService *ScoreService) GetScoreByTicketInRange(ctx context.Context, from time.Time, to time.Time) (res []domain.ScoreByTicket, err error) {
	res, err = scoreService.scoreRepository.FetchAggregateScoreByTicketInRange(ctx, from, to)
	if err != nil {
		return nil, err
	}
	return
}
