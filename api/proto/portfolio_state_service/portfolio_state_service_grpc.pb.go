// api/proto/portfolio_state_service.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: portfolio_state_service.proto

package portfolio_state_service

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
	PortfolioStateService_SavePortfolioState_FullMethodName  = "/portfoliostateservice.PortfolioStateService/SavePortfolioState"
	PortfolioStateService_LoadPortfolioState_FullMethodName  = "/portfoliostateservice.PortfolioStateService/LoadPortfolioState"
	PortfolioStateService_GetPortfolioHistory_FullMethodName = "/portfoliostateservice.PortfolioStateService/GetPortfolioHistory"
)

// PortfolioStateServiceClient is the client API for PortfolioStateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortfolioStateServiceClient interface {
	SavePortfolioState(ctx context.Context, in *PortfolioState, opts ...grpc.CallOption) (*SaveResponse, error)
	LoadPortfolioState(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*PortfolioState, error)
	GetPortfolioHistory(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*PortfolioHistory, error)
}

type portfolioStateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPortfolioStateServiceClient(cc grpc.ClientConnInterface) PortfolioStateServiceClient {
	return &portfolioStateServiceClient{cc}
}

func (c *portfolioStateServiceClient) SavePortfolioState(ctx context.Context, in *PortfolioState, opts ...grpc.CallOption) (*SaveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SaveResponse)
	err := c.cc.Invoke(ctx, PortfolioStateService_SavePortfolioState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioStateServiceClient) LoadPortfolioState(ctx context.Context, in *LoadRequest, opts ...grpc.CallOption) (*PortfolioState, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortfolioState)
	err := c.cc.Invoke(ctx, PortfolioStateService_LoadPortfolioState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioStateServiceClient) GetPortfolioHistory(ctx context.Context, in *HistoryRequest, opts ...grpc.CallOption) (*PortfolioHistory, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortfolioHistory)
	err := c.cc.Invoke(ctx, PortfolioStateService_GetPortfolioHistory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PortfolioStateServiceServer is the server API for PortfolioStateService service.
// All implementations must embed UnimplementedPortfolioStateServiceServer
// for forward compatibility
type PortfolioStateServiceServer interface {
	SavePortfolioState(context.Context, *PortfolioState) (*SaveResponse, error)
	LoadPortfolioState(context.Context, *LoadRequest) (*PortfolioState, error)
	GetPortfolioHistory(context.Context, *HistoryRequest) (*PortfolioHistory, error)
	mustEmbedUnimplementedPortfolioStateServiceServer()
}

// UnimplementedPortfolioStateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPortfolioStateServiceServer struct {
}

func (UnimplementedPortfolioStateServiceServer) SavePortfolioState(context.Context, *PortfolioState) (*SaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SavePortfolioState not implemented")
}
func (UnimplementedPortfolioStateServiceServer) LoadPortfolioState(context.Context, *LoadRequest) (*PortfolioState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPortfolioState not implemented")
}
func (UnimplementedPortfolioStateServiceServer) GetPortfolioHistory(context.Context, *HistoryRequest) (*PortfolioHistory, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolioHistory not implemented")
}
func (UnimplementedPortfolioStateServiceServer) mustEmbedUnimplementedPortfolioStateServiceServer() {}

// UnsafePortfolioStateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PortfolioStateServiceServer will
// result in compilation errors.
type UnsafePortfolioStateServiceServer interface {
	mustEmbedUnimplementedPortfolioStateServiceServer()
}

func RegisterPortfolioStateServiceServer(s grpc.ServiceRegistrar, srv PortfolioStateServiceServer) {
	s.RegisterService(&PortfolioStateService_ServiceDesc, srv)
}

func _PortfolioStateService_SavePortfolioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PortfolioState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioStateServiceServer).SavePortfolioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioStateService_SavePortfolioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioStateServiceServer).SavePortfolioState(ctx, req.(*PortfolioState))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioStateService_LoadPortfolioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioStateServiceServer).LoadPortfolioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioStateService_LoadPortfolioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioStateServiceServer).LoadPortfolioState(ctx, req.(*LoadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioStateService_GetPortfolioHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioStateServiceServer).GetPortfolioHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioStateService_GetPortfolioHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioStateServiceServer).GetPortfolioHistory(ctx, req.(*HistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PortfolioStateService_ServiceDesc is the grpc.ServiceDesc for PortfolioStateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PortfolioStateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "portfoliostateservice.PortfolioStateService",
	HandlerType: (*PortfolioStateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SavePortfolioState",
			Handler:    _PortfolioStateService_SavePortfolioState_Handler,
		},
		{
			MethodName: "LoadPortfolioState",
			Handler:    _PortfolioStateService_LoadPortfolioState_Handler,
		},
		{
			MethodName: "GetPortfolioHistory",
			Handler:    _PortfolioStateService_GetPortfolioHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "portfolio_state_service.proto",
}
