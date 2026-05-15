package clients

import (
	parkingv1 "github.com/nalrematvonesah/parking.proto/gen/parking/v1"
	sessionv1 "github.com/nalrematvonesah/session.proto/gen/session/v1"
	userv1 "github.com/nalrematvonesah/user.proto/gen/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	UserConn    *grpc.ClientConn
	ParkingConn *grpc.ClientConn
	SessionConn *grpc.ClientConn

	User    userv1.UserServiceClient
	Parking parkingv1.ParkingServiceClient
	Session sessionv1.SessionServiceClient
}

func Dial(userAddr, parkingAddr, sessionAddr string) (*Clients, error) {
	uConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	pConn, err := grpc.NewClient(parkingAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	sConn, err := grpc.NewClient(sessionAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Clients{
		UserConn: uConn, ParkingConn: pConn, SessionConn: sConn,
		User:    userv1.NewUserServiceClient(uConn),
		Parking: parkingv1.NewParkingServiceClient(pConn),
		Session: sessionv1.NewSessionServiceClient(sConn),
	}, nil
}

func (c *Clients) Close() {
	_ = c.UserConn.Close()
	_ = c.ParkingConn.Close()
	_ = c.SessionConn.Close()
}
