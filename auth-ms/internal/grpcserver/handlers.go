package grpcserver

import (
	"context"
	"strconv"

	pb "github.com/neglarken/clickhead/auth-ms/protobuf"
)

func (s *Server) SignIn(ctx context.Context, req *pb.AuthRequest) (*pb.TokenResponse, error) {
	res, err := s.service.User.SignIn(req)
	if err != nil {
		return &pb.TokenResponse{}, err
	}
	return res, nil
}

func (s *Server) SignUp(ctx context.Context, req *pb.AuthRequest) (*pb.SignUpResponse, error) {
	res := &pb.SignUpResponse{}

	if err := s.service.User.SignUp(req); err != nil {
		return res, err
	}
	res.Status = "success"
	return res, nil
}

func (s *Server) WhoAmI(ctx context.Context, req *pb.AccessTokenRequest) (*pb.UserResponse, error) {
	value := ctx.Value("id").(string)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return &pb.UserResponse{}, err
	}
	return &pb.UserResponse{Id: int32(id)}, nil
}
