package handler

import (
	"context"

	"session-service/internal/service"
)

type GRPCHandler struct {
	svc *service.Service
}

func New(svc *service.Service) *GRPCHandler {
	return &GRPCHandler{
		svc: svc,
	}
}

func (h *GRPCHandler) StartSession(
	ctx context.Context,
	vehicleID string,
	slotID int64,
) error {
	return h.svc.StartSession(
		ctx,
		vehicleID,
	)
}
