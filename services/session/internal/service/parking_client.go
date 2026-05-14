package service

import (
	"context"

	pb "github.com/nalrematvonesah/parking-proto/gen/parking/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ParkingClient struct {
	client pb.ParkingServiceClient
}

func NewParkingClient(addr string) (*ParkingClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	client := pb.NewParkingServiceClient(conn)

	return &ParkingClient{
		client: client,
	}, nil
}

func (p *ParkingClient) AssignSlot(
	ctx context.Context,
) (int64, error) {
	resp, err := p.client.AssignSlot(
		ctx,
		&pb.AssignSlotRequest{},
	)

	if err != nil {
		return 0, err
	}

	return resp.SlotId, nil
}
