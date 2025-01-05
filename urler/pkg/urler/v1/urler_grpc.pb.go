// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: urler.proto

package urler

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UrlerService_AddUser_FullMethodName = "/github.com.dbeleon.urler.UrlerService/AddUser"
	UrlerService_MakeUrl_FullMethodName = "/github.com.dbeleon.urler.UrlerService/MakeUrl"
	UrlerService_GetUrl_FullMethodName  = "/github.com.dbeleon.urler.UrlerService/GetUrl"
)

// UrlerServiceClient is the client API for UrlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlerServiceClient interface {
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error)
	MakeUrl(ctx context.Context, in *MakeUrlRequest, opts ...grpc.CallOption) (*MakeUrlResponse, error)
	GetUrl(ctx context.Context, in *GetUrlRequest, opts ...grpc.CallOption) (*GetUrlResponse, error)
}

type urlerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlerServiceClient(cc grpc.ClientConnInterface) UrlerServiceClient {
	return &urlerServiceClient{cc}
}

func (c *urlerServiceClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddUserResponse)
	err := c.cc.Invoke(ctx, UrlerService_AddUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *urlerServiceClient) MakeUrl(ctx context.Context, in *MakeUrlRequest, opts ...grpc.CallOption) (*MakeUrlResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MakeUrlResponse)
	err := c.cc.Invoke(ctx, UrlerService_MakeUrl_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *urlerServiceClient) GetUrl(ctx context.Context, in *GetUrlRequest, opts ...grpc.CallOption) (*GetUrlResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUrlResponse)
	err := c.cc.Invoke(ctx, UrlerService_GetUrl_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlerServiceServer is the server API for UrlerService service.
// All implementations must embed UnimplementedUrlerServiceServer
// for forward compatibility.
type UrlerServiceServer interface {
	AddUser(context.Context, *AddUserRequest) (*AddUserResponse, error)
	MakeUrl(context.Context, *MakeUrlRequest) (*MakeUrlResponse, error)
	GetUrl(context.Context, *GetUrlRequest) (*GetUrlResponse, error)
	mustEmbedUnimplementedUrlerServiceServer()
}

// UnimplementedUrlerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUrlerServiceServer struct{}

func (UnimplementedUrlerServiceServer) AddUser(context.Context, *AddUserRequest) (*AddUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedUrlerServiceServer) MakeUrl(context.Context, *MakeUrlRequest) (*MakeUrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeUrl not implemented")
}
func (UnimplementedUrlerServiceServer) GetUrl(context.Context, *GetUrlRequest) (*GetUrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUrl not implemented")
}
func (UnimplementedUrlerServiceServer) mustEmbedUnimplementedUrlerServiceServer() {}
func (UnimplementedUrlerServiceServer) testEmbeddedByValue()                      {}

// UnsafeUrlerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlerServiceServer will
// result in compilation errors.
type UnsafeUrlerServiceServer interface {
	mustEmbedUnimplementedUrlerServiceServer()
}

func RegisterUrlerServiceServer(s grpc.ServiceRegistrar, srv UrlerServiceServer) {
	// If the following call pancis, it indicates UnimplementedUrlerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UrlerService_ServiceDesc, srv)
}

func _UrlerService_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlerServiceServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UrlerService_AddUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlerServiceServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UrlerService_MakeUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MakeUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlerServiceServer).MakeUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UrlerService_MakeUrl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlerServiceServer).MakeUrl(ctx, req.(*MakeUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UrlerService_GetUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlerServiceServer).GetUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UrlerService_GetUrl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlerServiceServer).GetUrl(ctx, req.(*GetUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UrlerService_ServiceDesc is the grpc.ServiceDesc for UrlerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UrlerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.dbeleon.urler.UrlerService",
	HandlerType: (*UrlerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddUser",
			Handler:    _UrlerService_AddUser_Handler,
		},
		{
			MethodName: "MakeUrl",
			Handler:    _UrlerService_MakeUrl_Handler,
		},
		{
			MethodName: "GetUrl",
			Handler:    _UrlerService_GetUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "urler.proto",
}