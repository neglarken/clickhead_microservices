package grpcserver

import (
	"context"

	pb "github.com/neglarken/clickhead/some-ms/protobuf"
)

func (s *Server) Create(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	res := &pb.CreateItemResponse{}
	id, err := s.Repo.Item.Create(req)
	if err != nil {
		return res, err
	}
	return &pb.CreateItemResponse{Id: id}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	res := &pb.UpdateItemResponse{}
	item, err := s.Repo.Item.Edit(req)
	if err != nil {
		return res, err
	}
	res.Iten = item
	return res, nil
}

func (s *Server) Get(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	res := &pb.GetItemResponse{}
	items, err := s.Repo.Item.Get()
	if err != nil {
		return res, err
	}
	res.Items = items
	return res, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	res := &pb.DeleteItemResponse{}
	id, err := s.Repo.Item.Delete(req)
	if err != nil {
		return res, err
	}
	res.Id = id
	return res, nil
}
