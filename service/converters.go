package service

import (
	"github.com/fernandoalava/softwareengineer-test-task/grpc"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToGrpcRatingCategoryScore(ratingCategoryScore RatingCategoryScore) *grpc.RatingCategoryScore {
	return &grpc.RatingCategoryScore{
		RatingCategoryID:   int64(ratingCategoryScore.RatingCategoryID),
		RatingCategoryName: ratingCategoryScore.RatingCategoryName,
		Score:              float32(ratingCategoryScore.Score),
	}
}

func ToGrpcScoreByTicket(ticketScoreByCategory TicketScoreByCategory) *grpc.ScoreByTicket {
	return &grpc.ScoreByTicket{
		TicketId: int64(ticketScoreByCategory.TicketID),
		RatingCategoryScore: lo.Map(ticketScoreByCategory.RatingCategoryScores, func(ratingCategoryScore RatingCategoryScore, _ int) *grpc.RatingCategoryScore {
			return ToGrpcRatingCategoryScore(ratingCategoryScore)
		}),
	}
}

func ToGrpPeriodScoreWithRatings(periodScoreWithRatings PeriodScoreWithRatings) *grpc.PeriodScoreWithRatings {

	return &grpc.PeriodScoreWithRatings{
		From:    timestamppb.New(periodScoreWithRatings.From),
		To:      timestamppb.New(periodScoreWithRatings.To),
		Score:   float32(periodScoreWithRatings.Score),
		Ratings: int32(periodScoreWithRatings.Ratings),
	}
}

func ToGrpcCategoryScoreOverTime(categoryScoreOverTime CategoryScoreOverTime) *grpc.CategoryScoreOverTime {
	return &grpc.CategoryScoreOverTime{
		CategoryName: categoryScoreOverTime.CategoryName,
		PeriodScoreWithRatings: lo.Map(categoryScoreOverTime.PeriodScoresWithRatings, func(periodScoreWithRatings PeriodScoreWithRatings, _ int) *grpc.PeriodScoreWithRatings {
			return ToGrpPeriodScoreWithRatings(periodScoreWithRatings)
		}),
		TotalScore:  float32(categoryScoreOverTime.TotalScore),
		TotalRating: int32(categoryScoreOverTime.TotalRating),
	}
}

func ToGrpcPeriodScore(periodScore PeriodScore)*grpc.PeriodScore{
	return &grpc.PeriodScore{
		From: timestamppb.New(periodScore.From),
		To: timestamppb.New(periodScore.To),
		Score: float32(periodScore.Score),
	}
}