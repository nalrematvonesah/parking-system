package handler

import (
	"context"

	"user-service/internal/service"

	userv1 "github.com/nalrematvonesah/user.proto/gen/user/v1"
)

type GRPC struct {
	userv1.UnimplementedUserServiceServer
	svc *service.Service
}

func NewGRPC(svc *service.Service) *GRPC {
	return &GRPC{svc: svc}
}

func (h *GRPC) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.UserResponse, error) {
	id, err := h.svc.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userv1.UserResponse{UserId: id}, nil
}

func (h *GRPC) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.UserResponse, error) {
	id, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userv1.UserResponse{UserId: id}, nil
}

func (h *GRPC) Logout(ctx context.Context, req *userv1.UserRequest) (*userv1.Empty, error) {
	if err := h.svc.Logout(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	return &userv1.Empty{}, nil
}

func (h *GRPC) AddVehicle(ctx context.Context, req *userv1.AddVehicleRequest) (*userv1.Empty, error) {
	if err := h.svc.AddVehicle(ctx, req.GetUserId(), req.GetPlateNumber()); err != nil {
		return nil, err
	}
	return &userv1.Empty{}, nil
}

func (h *GRPC) DeleteVehicle(ctx context.Context, req *userv1.DeleteVehicleRequest) (*userv1.Empty, error) {
	if err := h.svc.DeleteVehicle(ctx, req.GetUserId(), req.GetPlateNumber()); err != nil {
		return nil, err
	}
	return &userv1.Empty{}, nil
}

func (h *GRPC) GetVehicles(ctx context.Context, req *userv1.UserRequest) (*userv1.VehicleList, error) {
	plates, err := h.svc.GetVehicles(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &userv1.VehicleList{Vehicles: plates}, nil
}
