// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: post.proto

package postsrvicepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PostClient is the client API for Post service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
}

type postClient struct {
	cc grpc.ClientConnInterface
}

func NewPostClient(cc grpc.ClientConnInterface) PostClient {
	return &postClient{cc}
}

func (c *postClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/IdeYarAPI.Post/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServer is the server API for Post service.
// All implementations must embed UnimplementedPostServer
// for forward compatibility
type PostServer interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	mustEmbedUnimplementedPostServer()
}

// UnimplementedPostServer must be embedded to have forward compatible implementations.
type UnimplementedPostServer struct {
}

func (UnimplementedPostServer) Create(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPostServer) mustEmbedUnimplementedPostServer() {}

// UnsafePostServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServer will
// result in compilation errors.
type UnsafePostServer interface {
	mustEmbedUnimplementedPostServer()
}

func RegisterPostServer(s grpc.ServiceRegistrar, srv PostServer) {
	s.RegisterService(&Post_ServiceDesc, srv)
}

func _Post_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/IdeYarAPI.Post/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Post_ServiceDesc is the grpc.ServiceDesc for Post service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Post_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "IdeYarAPI.Post",
	HandlerType: (*PostServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Post_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post.proto",
}