package server

import (
	"context"

	pb "github.com/fernandoalava/softwareengineer-test-task/grpc"

	"github.com/fernandoalava/softwareengineer-test-task/service"
	_ "modernc.org/sqlite"
)

type ScoreServer struct {
	pb.UnimplementedScoresServer
	scoreService *service.ScoreService
}

func (server *ScoreServer) GetScoreByTicket(request *pb.DateRangeRequest, stream pb.Scores_GetScoreByTicketServer) error {
	result, err := server.scoreService.GetScoreByTicket(stream.Context(), request.From.AsTime(), request.To.AsTime())
	if err != nil {
		return err
	}
	for _, r := range result {
		if err := stream.Send(service.ToGrpcScoreByTicket(r)); err != nil {
			return err
		}
	}
	return nil
}

func (server *ScoreServer) GetAggregatedCategoryScoresOverTime(request *pb.DateRangeRequest, stream pb.Scores_GetAggregatedCategoryScoresOverTimeServer) error {
	result, err := server.scoreService.GetAggregatedCategoryScoresOverTime(stream.Context(), request.From.AsTime(), request.To.AsTime())
	if err != nil {
		return err
	}
	for _, r := range result {
		if err := stream.Send(service.ToGrpcCategoryScoreOverTime(r)); err != nil {
			return err
		}
	}
	return nil
}

func (server *ScoreServer) GetOverAllQualityScore(ctx context.Context, request *pb.DateRangeRequest) (*pb.OverAllQualityScoreResponse, error) {
	result, err := server.scoreService.GetOverAllQualityScore(ctx, request.From.AsTime(), request.To.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.OverAllQualityScoreResponse{OverAllScore: float32(result)}, nil
}

func (server *ScoreServer) GetPeriodOverPeriodScoreChange(ctx context.Context, request *pb.DateRangeRequest) (*pb.GetPeriodOverPeriodScoreChangeResponse, error) {
	result, err := server.scoreService.GetPeriodOverPeriodScoreChange(ctx, request.From.AsTime(), request.To.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetPeriodOverPeriodScoreChangeResponse{
		CurrentPeriod:   service.ToGrpcPeriodScore(result.CurrentPeriod),
		PreviousPeriod:  service.ToGrpcPeriodScore(result.PreviousPeriod),
		ScoreDifference: float32(result.ScoreDifference),
	}, nil
}

func NewScoreServer(scoreService *service.ScoreService) *ScoreServer {
	server := &ScoreServer{scoreService: scoreService}
	return server
}
