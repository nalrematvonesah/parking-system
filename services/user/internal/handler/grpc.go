package handler

import (
	"context"

	"user-service/internal/service"

	userpb "github.com/nalrematvonesah/parking-proto/gen/user/v1"
)

type GRPC struct {
	userpb.UnimplementedUserServiceServer
	svc *service.Service
}

func NewGRPC(svc *service.Service) *GRPC {
	return &GRPC{svc: svc}
}

func (h *GRPC) mustEmbedUnimplementedUserServiceServer() {}

func (h *GRPC) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	id, err := h.svc.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{UserId: id}, nil
}

func (h *GRPC) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.UserResponse, error) {
	id, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{UserId: id}, nil
}

func (h *GRPC) Logout(ctx context.Context, req *userpb.UserRequest) (*userpb.Empty, error) {
	if err := h.svc.Logout(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	return &userpb.Empty{}, nil
}

func (h *GRPC) AddVehicle(ctx context.Context, req *userpb.AddVehicleRequest) (*userpb.Empty, error) {
	if err := h.svc.AddVehicle(ctx, req.GetUserId(), req.GetPlateNumber()); err != nil {
		return nil, err
	}
	return &userpb.Empty{}, nil
}

func (h *GRPC) DeleteVehicle(ctx context.Context, req *userpb.AddVehicleRequest) (*userpb.Empty, error) {
	if err := h.svc.DeleteVehicle(ctx, req.GetUserId(), req.GetPlateNumber()); err != nil {
		return nil, err
	}
	return &userpb.Empty{}, nil
}

func (h *GRPC) GetVehicles(ctx context.Context, req *userpb.UserRequest) (*userpb.VehicleList, error) {
	plates, err := h.svc.GetVehicles(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &userpb.VehicleList{Vehicles: plates}, nil
}
