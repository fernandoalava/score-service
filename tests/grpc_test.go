package tests

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net"
	"testing"

	"github.com/fernandoalava/softwareengineer-test-task/repository"
	"github.com/fernandoalava/softwareengineer-test-task/server"
	"github.com/fernandoalava/softwareengineer-test-task/service"
	"github.com/fernandoalava/softwareengineer-test-task/util"
	"github.com/stretchr/testify/assert"

	pb "github.com/fernandoalava/softwareengineer-test-task/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
	_ "modernc.org/sqlite"
)

func grpcServer(ctx context.Context) (pb.ScoresClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)
	baseServer := grpc.NewServer()

	db, err := sql.Open("sqlite", "../database.db")
	if err != nil {
		log.Fatal(err)
	}

	ratingCategoryRepository := repository.NewRatingCategoryRepository(db)
	ratingRepository := repository.NewRatingRepository(db)
	ticketRepository := repository.NewTicketRepository(db)

	scoreService := service.NewScoreService(ticketRepository, ratingCategoryRepository, ratingRepository)

	server := server.NewScoreServer(scoreService)

	pb.RegisterScoresServer(baseServer, server)

	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := db.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
		err = lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewScoresClient(conn)

	return client, closer
}

func TestGrpcGetScoreByTicket(t *testing.T) {
	ctx := context.Background()
	client, closer := grpcServer(ctx)
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	out, err := client.GetScoreByTicket(ctx, &pb.DateRangeRequest{From: timestamppb.New(from), To: timestamppb.New(to)})
	var outs []*pb.ScoreByTicket

	for {
		o, err := out.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		outs = append(outs, o)
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, outs)
}

func TestGrpcGetOverAllQualityScore(t *testing.T) {
	ctx := context.Background()
	client, closer := grpcServer(ctx)
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	out, err := client.GetOverAllQualityScore(ctx, &pb.DateRangeRequest{From: timestamppb.New(from), To: timestamppb.New(to)})

	assert.Nil(t, err)
	assert.Equal(t, float32(35.8), out.OverAllScore)
}

func TestGrpcGetAggregatedCategoryScoresOverTime(t *testing.T) {
	ctx := context.Background()
	client, closer := grpcServer(ctx)
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	out, err := client.GetAggregatedCategoryScoresOverTime(ctx, &pb.DateRangeRequest{From: timestamppb.New(from), To: timestamppb.New(to)})

	var outs []*pb.CategoryScoreOverTime

	for {
		o, err := out.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		outs = append(outs, o)
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, outs)
}

func TestGrpcGetPeriodOverPeriodScoreChange(t *testing.T) {
	ctx := context.Background()
	client, closer := grpcServer(ctx)
	defer closer()

	from, _ := util.StringToTime("2019-07-17T00:00:00")
	to, _ := util.StringToTime("2019-07-17T23:59:00")

	out, err := client.GetPeriodOverPeriodScoreChange(ctx, &pb.DateRangeRequest{From: timestamppb.New(from), To: timestamppb.New(to)})

	assert.Nil(t, err)
	assert.Equal(t, float32(1.52), out.ScoreDifference)
}
