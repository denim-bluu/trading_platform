// api/proto/backtesting_service.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: backtesting_service.proto

package backtesting_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	BacktestingService_RunBacktest_FullMethodName       = "/backtestingservice.BacktestingService/RunBacktest"
	BacktestingService_GetBacktestStatus_FullMethodName = "/backtestingservice.BacktestingService/GetBacktestStatus"
)

// BacktestingServiceClient is the client API for BacktestingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BacktestingServiceClient interface {
	RunBacktest(ctx context.Context, in *BacktestRequest, opts ...grpc.CallOption) (*BacktestResult, error)
	GetBacktestStatus(ctx context.Context, in *BacktestStatusRequest, opts ...grpc.CallOption) (*BacktestStatus, error)
}

type backtestingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBacktestingServiceClient(cc grpc.ClientConnInterface) BacktestingServiceClient {
	return &backtestingServiceClient{cc}
}

func (c *backtestingServiceClient) RunBacktest(ctx context.Context, in *BacktestRequest, opts ...grpc.CallOption) (*BacktestResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BacktestResult)
	err := c.cc.Invoke(ctx, BacktestingService_RunBacktest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtestingServiceClient) GetBacktestStatus(ctx context.Context, in *BacktestStatusRequest, opts ...grpc.CallOption) (*BacktestStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BacktestStatus)
	err := c.cc.Invoke(ctx, BacktestingService_GetBacktestStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BacktestingServiceServer is the server API for BacktestingService service.
// All implementations must embed UnimplementedBacktestingServiceServer
// for forward compatibility
type BacktestingServiceServer interface {
	RunBacktest(context.Context, *BacktestRequest) (*BacktestResult, error)
	GetBacktestStatus(context.Context, *BacktestStatusRequest) (*BacktestStatus, error)
	mustEmbedUnimplementedBacktestingServiceServer()
}

// UnimplementedBacktestingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBacktestingServiceServer struct {
}

func (UnimplementedBacktestingServiceServer) RunBacktest(context.Context, *BacktestRequest) (*BacktestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunBacktest not implemented")
}
func (UnimplementedBacktestingServiceServer) GetBacktestStatus(context.Context, *BacktestStatusRequest) (*BacktestStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBacktestStatus not implemented")
}
func (UnimplementedBacktestingServiceServer) mustEmbedUnimplementedBacktestingServiceServer() {}

// UnsafeBacktestingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BacktestingServiceServer will
// result in compilation errors.
type UnsafeBacktestingServiceServer interface {
	mustEmbedUnimplementedBacktestingServiceServer()
}

func RegisterBacktestingServiceServer(s grpc.ServiceRegistrar, srv BacktestingServiceServer) {
	s.RegisterService(&BacktestingService_ServiceDesc, srv)
}

func _BacktestingService_RunBacktest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BacktestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestingServiceServer).RunBacktest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktestingService_RunBacktest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestingServiceServer).RunBacktest(ctx, req.(*BacktestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktestingService_GetBacktestStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BacktestStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestingServiceServer).GetBacktestStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BacktestingService_GetBacktestStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestingServiceServer).GetBacktestStatus(ctx, req.(*BacktestStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BacktestingService_ServiceDesc is the grpc.ServiceDesc for BacktestingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BacktestingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "backtestingservice.BacktestingService",
	HandlerType: (*BacktestingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunBacktest",
			Handler:    _BacktestingService_RunBacktest_Handler,
		},
		{
			MethodName: "GetBacktestStatus",
			Handler:    _BacktestingService_GetBacktestStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backtesting_service.proto",
}
