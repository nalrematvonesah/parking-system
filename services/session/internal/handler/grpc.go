package handler

import (
	"context"

	"session-service/internal/service"

	sessionv1 "github.com/nalrematvonesah/session.proto/gen/session/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPC struct {
	sessionv1.UnimplementedSessionServiceServer
	svc *service.Service
}

func New(svc *service.Service) *GRPC {
	return &GRPC{svc: svc}
}

func (h *GRPC) StartSession(ctx context.Context, req *sessionv1.StartRequest) (*sessionv1.SessionResponse, error) {
	if req.GetUserId() == 0 || req.GetVehicleId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id and vehicle_id are required")
	}
	s, err := h.svc.StartSession(ctx, req.GetUserId(), req.GetVehicleId(), req.GetRequestId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sessionv1.SessionResponse{
		SessionId:     s.ID,
		SlotId:        s.SlotID,
		VehicleId:     s.VehicleID,
		StartTimeUnix: s.StartTime.Unix(),
	}, nil
}

func (h *GRPC) EndSession(ctx context.Context, req *sessionv1.EndRequest) (*sessionv1.SessionResponse, error) {
	if req.GetSessionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	s, amount, err := h.svc.EndSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &sessionv1.SessionResponse{
		SessionId:     s.ID,
		SlotId:        s.SlotID,
		VehicleId:     s.VehicleID,
		StartTimeUnix: s.StartTime.Unix(),
		Amount:        amount,
	}
	if s.EndTime != nil {
		resp.EndTimeUnix = s.EndTime.Unix()
	}
	return resp, nil
}

func (h *GRPC) GetSession(ctx context.Context, req *sessionv1.SessionRequest) (*sessionv1.SessionResponse, error) {
	if req.GetSessionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	s, err := h.svc.GetSession(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if s == nil {
		return nil, status.Error(codes.NotFound, "session not found")
	}
	resp := &sessionv1.SessionResponse{
		SessionId:     s.ID,
		SlotId:        s.SlotID,
		VehicleId:     s.VehicleID,
		StartTimeUnix: s.StartTime.Unix(),
	}
	if s.EndTime != nil {
		resp.EndTimeUnix = s.EndTime.Unix()
	}
	return resp, nil
}

func (h *GRPC) CalculatePrice(ctx context.Context, req *sessionv1.SessionRequest) (*sessionv1.PriceResponse, error) {
	if req.GetSessionId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	amount, elapsed, err := h.svc.PriceFor(ctx, req.GetSessionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &sessionv1.PriceResponse{
		Amount:         amount,
		ElapsedSeconds: int64(elapsed.Seconds()),
	}, nil
}

func (h *GRPC) GetActiveSessions(ctx context.Context, req *sessionv1.UserRequest) (*sessionv1.SessionList, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	sessions, err := h.svc.ListActiveByUser(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := make([]*sessionv1.SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		item := &sessionv1.SessionResponse{
			SessionId:     s.ID,
			SlotId:        s.SlotID,
			VehicleId:     s.VehicleID,
			StartTimeUnix: s.StartTime.Unix(),
		}
		if s.EndTime != nil {
			item.EndTimeUnix = s.EndTime.Unix()
		}
		out = append(out, item)
	}
	return &sessionv1.SessionList{Sessions: out}, nil
}
