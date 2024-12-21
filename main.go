package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "github.com/fernandoalava/softwareengineer-test-task/grpc"
	"github.com/fernandoalava/softwareengineer-test-task/repository"
	"github.com/fernandoalava/softwareengineer-test-task/server"
	"github.com/fernandoalava/softwareengineer-test-task/service"
	"github.com/fernandoalava/softwareengineer-test-task/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port     = util.GetEnv("PORT", "9000")
	database = util.GetEnv("DB_PATH", "")
)

func main() {
	log.Printf("trying to listen on port %s", port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if len(database) == 0 {
		log.Fatalf("DB_PATH is empty or undefined")
	}
	log.Printf("loading database %s", database)
	db, err := sql.Open("sqlite", database)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("database %s loaded successfully", database)
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
