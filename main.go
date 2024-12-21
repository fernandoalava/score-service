package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/fernandoalava/softwareengineer-test-task/grpc"
	"github.com/fernandoalava/softwareengineer-test-task/repository"
	"github.com/fernandoalava/softwareengineer-test-task/server"
	"github.com/fernandoalava/softwareengineer-test-task/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port     = flag.Int("port", 5000, "The server port")
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
	scoreRepository := repository.NewScoreRepository(db)

	scoreService := service.NewScoreService(ratingCategoryRepository, scoreRepository)

	server := server.NewScoreServer(scoreService)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterScoresServer(grpcServer, server)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
