package main

import (
	"context"

	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/fernandoalava/softwareengineer-test-task/grpc"
	"github.com/fernandoalava/softwareengineer-test-task/repository"
	"github.com/fernandoalava/softwareengineer-test-task/service"
	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

type ScoreServer struct {
	pb.UnimplementedScoresServer
	scoreService *service.ScoreService
	ctx          context.Context
}

func (server *ScoreServer) GetScoreByTicket(request *pb.DateRangeRequest, stream pb.Scores_GetScoreByTicketServer) error {
	result, err := server.scoreService.GetScoreByTicket(server.ctx, request.From.AsTime(), request.To.AsTime())
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
	result, err := server.scoreService.GetAggregatedCategoryScoresOverTime(server.ctx, request.From.AsTime(), request.To.AsTime())
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

func newScoreServer(scoreService *service.ScoreService) *ScoreServer {
	server := &ScoreServer{scoreService: scoreService, ctx: context.Background()}
	return server
}

var (
	port     = flag.Int("port", 50051, "The server port")
	database = flag.String("database", "database.db", "The database path")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := sql.Open("sqlite", *database)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	ratingCategoryRepository := repository.NewRatingCategoryRepository(db)
	ratingRepository := repository.NewRatingRepository(db)
	ticketRepository := repository.NewTicketRepository(db)
	scoreService := service.NewScoreService(ticketRepository, ratingCategoryRepository, ratingRepository)
	server := newScoreServer(scoreService)
	grpcServer := grpc.NewServer()
	pb.RegisterScoresServer(grpcServer, server)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
