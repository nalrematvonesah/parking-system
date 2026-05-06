package handler

import (
	"context"

	"user-service/internal/service"

	parkingv1 "github.com/nalrematvonesah/parking-proto/gen/user/v1"
)

type GRPC struct {
	parkingv1.UnimplementedUserServiceServer
	svc *service.Service
}

func NewGRPC(svc *service.Service) *GRPC {
	return &GRPC{svc: svc}
}

func (h *GRPC) Register(ctx context.Context, req *parkingv1.RegisterRequest) (*parkingv1.UserResponse, error) {
	id, err := h.svc.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &parkingv1.UserResponse{UserId: id}, nil
}

func (h *GRPC) Login(ctx context.Context, req *parkingv1.LoginRequest) (*parkingv1.UserResponse, error) {
	id, err := h.svc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &parkingv1.UserResponse{UserId: id}, nil
}

func (h *GRPC) AddVehicle(ctx context.Context, req *parkingv1.AddVehicleRequest) (*parkingv1.Empty, error) {
	if err := h.svc.AddVehicle(ctx, req.GetUserId(), req.GetPlateNumber()); err != nil {
		return nil, err
	}
	return &parkingv1.Empty{}, nil
}

func (h *GRPC) GetVehicles(ctx context.Context, req *parkingv1.UserRequest) (*parkingv1.VehicleList, error) {
	plates, err := h.svc.GetVehicles(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &parkingv1.VehicleList{Vehicles: plates}, nil
}
