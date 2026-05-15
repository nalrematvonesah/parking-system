package handler

import (
	"context"

	"parking-service/internal/service"

	parkingpb "github.com/nalrematvonesah/parking.proto/gen/parking/v1"
)

type GRPC struct {
	parkingpb.UnimplementedParkingServiceServer
	svc *service.Service
}

func NewGRPC(svc *service.Service) *GRPC {
	return &GRPC{svc: svc}
}

func (h *GRPC) AssignSlot(ctx context.Context, _ *parkingpb.AssignSlotRequest) (*parkingpb.AssignSlotResponse, error) {
	id, err := h.svc.AssignSlot(ctx)
	if err != nil {
		return nil, err
	}
	return &parkingpb.AssignSlotResponse{SlotId: id}, nil
}

func (h *GRPC) ReleaseSlot(ctx context.Context, req *parkingpb.ReleaseSlotRequest) (*parkingpb.Empty, error) {
	if err := h.svc.ReleaseSlot(ctx, req.SlotId); err != nil {
		return nil, err
	}
	return &parkingpb.Empty{}, nil
}

func (h *GRPC) GetAvailableSlots(ctx context.Context, _ *parkingpb.Empty) (*parkingpb.AvailableSlotsResponse, error) {
	c, err := h.svc.GetAvailable(ctx)
	if err != nil {
		return nil, err
	}
	return &parkingpb.AvailableSlotsResponse{Count: c}, nil
}
