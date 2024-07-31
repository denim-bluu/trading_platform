// api/proto/portfolio_state_service.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.12.4
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
	PortfolioStateService_GetPortfolioState_FullMethodName    = "/portfoliostateservice.PortfolioStateService/GetPortfolioState"
	PortfolioStateService_UpdatePortfolioState_FullMethodName = "/portfoliostateservice.PortfolioStateService/UpdatePortfolioState"
)

// PortfolioStateServiceClient is the client API for PortfolioStateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortfolioStateServiceClient interface {
	GetPortfolioState(ctx context.Context, in *GetPortfolioStateRequest, opts ...grpc.CallOption) (*PortfolioState, error)
	UpdatePortfolioState(ctx context.Context, in *UpdatePortfolioStateRequest, opts ...grpc.CallOption) (*UpdatePortfolioStateResponse, error)
}

type portfolioStateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPortfolioStateServiceClient(cc grpc.ClientConnInterface) PortfolioStateServiceClient {
	return &portfolioStateServiceClient{cc}
}

func (c *portfolioStateServiceClient) GetPortfolioState(ctx context.Context, in *GetPortfolioStateRequest, opts ...grpc.CallOption) (*PortfolioState, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PortfolioState)
	err := c.cc.Invoke(ctx, PortfolioStateService_GetPortfolioState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portfolioStateServiceClient) UpdatePortfolioState(ctx context.Context, in *UpdatePortfolioStateRequest, opts ...grpc.CallOption) (*UpdatePortfolioStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePortfolioStateResponse)
	err := c.cc.Invoke(ctx, PortfolioStateService_UpdatePortfolioState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PortfolioStateServiceServer is the server API for PortfolioStateService service.
// All implementations must embed UnimplementedPortfolioStateServiceServer
// for forward compatibility
type PortfolioStateServiceServer interface {
	GetPortfolioState(context.Context, *GetPortfolioStateRequest) (*PortfolioState, error)
	UpdatePortfolioState(context.Context, *UpdatePortfolioStateRequest) (*UpdatePortfolioStateResponse, error)
	mustEmbedUnimplementedPortfolioStateServiceServer()
}

// UnimplementedPortfolioStateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPortfolioStateServiceServer struct {
}

func (UnimplementedPortfolioStateServiceServer) GetPortfolioState(context.Context, *GetPortfolioStateRequest) (*PortfolioState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPortfolioState not implemented")
}
func (UnimplementedPortfolioStateServiceServer) UpdatePortfolioState(context.Context, *UpdatePortfolioStateRequest) (*UpdatePortfolioStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePortfolioState not implemented")
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

func _PortfolioStateService_GetPortfolioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPortfolioStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioStateServiceServer).GetPortfolioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioStateService_GetPortfolioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioStateServiceServer).GetPortfolioState(ctx, req.(*GetPortfolioStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PortfolioStateService_UpdatePortfolioState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePortfolioStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortfolioStateServiceServer).UpdatePortfolioState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PortfolioStateService_UpdatePortfolioState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortfolioStateServiceServer).UpdatePortfolioState(ctx, req.(*UpdatePortfolioStateRequest))
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
			MethodName: "GetPortfolioState",
			Handler:    _PortfolioStateService_GetPortfolioState_Handler,
		},
		{
			MethodName: "UpdatePortfolioState",
			Handler:    _PortfolioStateService_UpdatePortfolioState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "portfolio_state_service.proto",
}
