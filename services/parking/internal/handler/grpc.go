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

func (h *GRPC) GetSlot(ctx context.Context, req *parkingpb.GetSlotRequest) (*parkingpb.SlotResponse, error) {
	s, err := h.svc.GetSlot(ctx, req.GetSlotId())
	if err != nil {
		return nil, err
	}

	return &parkingpb.SlotResponse{
		SlotId:     s.ID,
		IsOccupied: s.IsOccupied,
	}, nil
}

func (h *GRPC) ListAllSlots(ctx context.Context, _ *parkingpb.Empty) (*parkingpb.ListSlotsResponse, error) {
	slots, err := h.svc.ListAllSlots(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*parkingpb.SlotResponse, 0, len(slots))
	for _, s := range slots {
		out = append(out, &parkingpb.SlotResponse{
			SlotId:     s.ID,
			IsOccupied: s.IsOccupied,
		})
	}

	return &parkingpb.ListSlotsResponse{Slots: out}, nil
}

func (h *GRPC) AddSlots(ctx context.Context, req *parkingpb.AddSlotsRequest) (*parkingpb.Empty, error) {
	if err := h.svc.AddSlots(ctx, req.GetCount()); err != nil {
		return nil, err
	}

	return &parkingpb.Empty{}, nil
}
