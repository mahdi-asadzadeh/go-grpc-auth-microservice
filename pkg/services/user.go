package services

import (
	"context"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/db"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/models"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/pb"
	"github.com/mahdi-asadzadeh/go-grpc-auth-microservice/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	H   db.Handler
	Jwt utils.JwtWrapper
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	dataUser := req.GetUser()
	err := s.H.DB.Create(&models.User{Email: dataUser.GetEmail(), Password: utils.HashPassword(dataUser.GetPassword())}).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "Invalid user data.")
	}
	return &pb.RegisterResponse{Id: 0, Email: dataUser.GetEmail(), CreateAt: "", UpdateAt: ""}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User
	dataUser := req.GetUser()
	result := s.H.DB.Where(&models.User{Email: dataUser.GetEmail()}).First(&user)
	if result == nil {
		return nil, status.Error(codes.NotFound, "Not found user.")
	}
	match := utils.CheckPasswordHash(dataUser.GetPassword(), user.Password)
	if !match {
		return nil, status.Error(codes.NotFound, "Not found user.")
	}
	token, err := s.Jwt.GenerateToken(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.Jwt.ValidateToken(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.NotFound, "Not found user.")
	}
	var user models.User
	err = s.H.DB.Where("email = ?", claims.Email).First(&user).Error
	if err != nil {
		return nil, status.Error(codes.NotFound, "Not found user.")
	}
	return &pb.ValidateResponse{Email: user.Email, Id: int64(user.ID)}, nil
}
