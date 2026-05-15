package service

import (
	"context"

	parkingv1 "github.com/nalrematvonesah/parking.proto/gen/parking/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ParkingClient struct {
	conn   *grpc.ClientConn
	client parkingv1.ParkingServiceClient
}

func NewParkingClient(addr string) (*ParkingClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return &ParkingClient{
		conn:   conn,
		client: parkingv1.NewParkingServiceClient(conn),
	}, nil
}

func (p *ParkingClient) AssignSlot(ctx context.Context, requestID string) (int64, error) {
	resp, err := p.client.AssignSlot(ctx, &parkingv1.AssignSlotRequest{RequestId: requestID})
	if err != nil {
		return 0, err
	}
	return resp.GetSlotId(), nil
}

func (p *ParkingClient) ReleaseSlot(ctx context.Context, slotID int64) error {
	_, err := p.client.ReleaseSlot(ctx, &parkingv1.ReleaseSlotRequest{SlotId: slotID})
	return err
}

func (p *ParkingClient) Close() error {
	return p.conn.Close()
}
