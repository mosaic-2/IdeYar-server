// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package LivenessServicePb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LivenessClient is the client API for Liveness service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LivenessClient interface {
	CheckLiveness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CheckLivenessResponse, error)
}

type livenessClient struct {
	cc grpc.ClientConnInterface
}

func NewLivenessClient(cc grpc.ClientConnInterface) LivenessClient {
	return &livenessClient{cc}
}

func (c *livenessClient) CheckLiveness(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CheckLivenessResponse, error) {
	out := new(CheckLivenessResponse)
	err := c.cc.Invoke(ctx, "/IdeYarAPI.Liveness/CheckLiveness", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LivenessServer is the server API for Liveness service.
// All implementations must embed UnimplementedLivenessServer
// for forward compatibility
type LivenessServer interface {
	CheckLiveness(context.Context, *emptypb.Empty) (*CheckLivenessResponse, error)
	mustEmbedUnimplementedLivenessServer()
}

// UnimplementedLivenessServer must be embedded to have forward compatible implementations.
type UnimplementedLivenessServer struct {
}

func (UnimplementedLivenessServer) CheckLiveness(context.Context, *emptypb.Empty) (*CheckLivenessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckLiveness not implemented")
}
func (UnimplementedLivenessServer) mustEmbedUnimplementedLivenessServer() {}

// UnsafeLivenessServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LivenessServer will
// result in compilation errors.
type UnsafeLivenessServer interface {
	mustEmbedUnimplementedLivenessServer()
}

func RegisterLivenessServer(s grpc.ServiceRegistrar, srv LivenessServer) {
	s.RegisterService(&Liveness_ServiceDesc, srv)
}

func _Liveness_CheckLiveness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivenessServer).CheckLiveness(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IdeYarAPI.Liveness/CheckLiveness",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivenessServer).CheckLiveness(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Liveness_ServiceDesc is the grpc.ServiceDesc for Liveness service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Liveness_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IdeYarAPI.Liveness",
	HandlerType: (*LivenessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckLiveness",
			Handler:    _Liveness_CheckLiveness_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/liveness.proto",
}
