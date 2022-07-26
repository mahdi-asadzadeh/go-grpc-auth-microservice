package main

import (
	"fmt"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/config"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/db"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/pb"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/services"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting auth grpc server ...")
	// Load config
	config.LoadSettings(true)

	// Database config
	DB_URL := os.Getenv("DB_URL")
	h := db.Init(DB_URL)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	jwtWrapper := utils.JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET"),
		Issuer:          "go-grpc-auth-microservice",
		ExpirationHours: 24 * 365,
	}
	// Register product service
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &services.Server{H: h, Jwt: jwtWrapper})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
