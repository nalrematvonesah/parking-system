package pb

import (
	"context"

	userpb "github.com/nalrematvonesah/parking-proto/gen/user/v1"
	"google.golang.org/grpc"
)

type ExtendedUserServiceServer interface {
	Register(context.Context, *userpb.RegisterRequest) (*userpb.UserResponse, error)
	Login(context.Context, *userpb.LoginRequest) (*userpb.UserResponse, error)
	Logout(context.Context, *userpb.UserRequest) (*userpb.Empty, error)
	AddVehicle(context.Context, *userpb.AddVehicleRequest) (*userpb.Empty, error)
	DeleteVehicle(context.Context, *userpb.AddVehicleRequest) (*userpb.Empty, error)
	GetVehicles(context.Context, *userpb.UserRequest) (*userpb.VehicleList, error)
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterExtendedUserServiceServer(s grpc.ServiceRegistrar, srv ExtendedUserServiceServer) {
	s.RegisterService(&userServiceExtendedDesc, srv)
}

var userServiceExtendedDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.UserService",
	HandlerType: (*ExtendedUserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: registerHandler},
		{MethodName: "Login", Handler: loginHandler},
		{MethodName: "Logout", Handler: logoutHandler},
		{MethodName: "AddVehicle", Handler: addVehicleHandler},
		{MethodName: "DeleteVehicle", Handler: deleteVehicleHandler},
		{MethodName: "GetVehicles", Handler: getVehiclesHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/v1/user.proto",
}

func registerHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/Register"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).Register(ctx, req.(*userpb.RegisterRequest))
	})
}

func loginHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/Login"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).Login(ctx, req.(*userpb.LoginRequest))
	})
}

func logoutHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/Logout"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).Logout(ctx, req.(*userpb.UserRequest))
	})
}

func addVehicleHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.AddVehicleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).AddVehicle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/AddVehicle"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).AddVehicle(ctx, req.(*userpb.AddVehicleRequest))
	})
}

func deleteVehicleHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.AddVehicleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).DeleteVehicle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/DeleteVehicle"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).DeleteVehicle(ctx, req.(*userpb.AddVehicleRequest))
	})
}

func getVehiclesHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(userpb.UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExtendedUserServiceServer).GetVehicles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/user.v1.UserService/GetVehicles"}
	return interceptor(ctx, in, info, func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExtendedUserServiceServer).GetVehicles(ctx, req.(*userpb.UserRequest))
	})
}
