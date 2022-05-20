package main

import (
	"log"
	"net"

	"github.com/himanshuk42/product/pkg/config"
	"github.com/himanshuk42/product/pkg/db"
	"github.com/himanshuk42/product/pkg/pb"
	"github.com/himanshuk42/product/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v\n", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("couldn't listen: %v\n", err)
	}

	log.Println("Product Microservice listening on Port", c.Port)

	s := services.Server{
		H: h,
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	pb.RegisterProductServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
