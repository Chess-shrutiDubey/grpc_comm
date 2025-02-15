// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: performance.proto

package generated

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

// PerformanceTestClient is the client API for PerformanceTest service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PerformanceTestClient interface {
	// Simple RPC for RTT testing
	SimpleRPC(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	// Streaming for bandwidth testing
	StreamData(ctx context.Context, opts ...grpc.CallOption) (PerformanceTest_StreamDataClient, error)
	// Server streaming
	ServerStream(ctx context.Context, in *Message, opts ...grpc.CallOption) (PerformanceTest_ServerStreamClient, error)
	// Client streaming
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (PerformanceTest_ClientStreamClient, error)
}

type performanceTestClient struct {
	cc grpc.ClientConnInterface
}

func NewPerformanceTestClient(cc grpc.ClientConnInterface) PerformanceTestClient {
	return &performanceTestClient{cc}
}

func (c *performanceTestClient) SimpleRPC(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/performance.PerformanceTest/SimpleRPC", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *performanceTestClient) StreamData(ctx context.Context, opts ...grpc.CallOption) (PerformanceTest_StreamDataClient, error) {
	stream, err := c.cc.NewStream(ctx, &PerformanceTest_ServiceDesc.Streams[0], "/performance.PerformanceTest/StreamData", opts...)
	if err != nil {
		return nil, err
	}
	x := &performanceTestStreamDataClient{stream}
	return x, nil
}

type PerformanceTest_StreamDataClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type performanceTestStreamDataClient struct {
	grpc.ClientStream
}

func (x *performanceTestStreamDataClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *performanceTestStreamDataClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *performanceTestClient) ServerStream(ctx context.Context, in *Message, opts ...grpc.CallOption) (PerformanceTest_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PerformanceTest_ServiceDesc.Streams[1], "/performance.PerformanceTest/ServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &performanceTestServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PerformanceTest_ServerStreamClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type performanceTestServerStreamClient struct {
	grpc.ClientStream
}

func (x *performanceTestServerStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *performanceTestClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (PerformanceTest_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PerformanceTest_ServiceDesc.Streams[2], "/performance.PerformanceTest/ClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &performanceTestClientStreamClient{stream}
	return x, nil
}

type PerformanceTest_ClientStreamClient interface {
	Send(*Message) error
	CloseAndRecv() (*Message, error)
	grpc.ClientStream
}

type performanceTestClientStreamClient struct {
	grpc.ClientStream
}

func (x *performanceTestClientStreamClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *performanceTestClientStreamClient) CloseAndRecv() (*Message, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PerformanceTestServer is the server API for PerformanceTest service.
// All implementations must embed UnimplementedPerformanceTestServer
// for forward compatibility
type PerformanceTestServer interface {
	// Simple RPC for RTT testing
	SimpleRPC(context.Context, *Message) (*Message, error)
	// Streaming for bandwidth testing
	StreamData(PerformanceTest_StreamDataServer) error
	// Server streaming
	ServerStream(*Message, PerformanceTest_ServerStreamServer) error
	// Client streaming
	ClientStream(PerformanceTest_ClientStreamServer) error
	mustEmbedUnimplementedPerformanceTestServer()
}

// UnimplementedPerformanceTestServer must be embedded to have forward compatible implementations.
type UnimplementedPerformanceTestServer struct {
}

func (UnimplementedPerformanceTestServer) SimpleRPC(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SimpleRPC not implemented")
}
func (UnimplementedPerformanceTestServer) StreamData(PerformanceTest_StreamDataServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamData not implemented")
}
func (UnimplementedPerformanceTestServer) ServerStream(*Message, PerformanceTest_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (UnimplementedPerformanceTestServer) ClientStream(PerformanceTest_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedPerformanceTestServer) mustEmbedUnimplementedPerformanceTestServer() {}

// UnsafePerformanceTestServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PerformanceTestServer will
// result in compilation errors.
type UnsafePerformanceTestServer interface {
	mustEmbedUnimplementedPerformanceTestServer()
}

func RegisterPerformanceTestServer(s grpc.ServiceRegistrar, srv PerformanceTestServer) {
	s.RegisterService(&PerformanceTest_ServiceDesc, srv)
}

func _PerformanceTest_SimpleRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformanceTestServer).SimpleRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/performance.PerformanceTest/SimpleRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformanceTestServer).SimpleRPC(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerformanceTest_StreamData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PerformanceTestServer).StreamData(&performanceTestStreamDataServer{stream})
}

type PerformanceTest_StreamDataServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type performanceTestStreamDataServer struct {
	grpc.ServerStream
}

func (x *performanceTestStreamDataServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *performanceTestStreamDataServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PerformanceTest_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Message)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PerformanceTestServer).ServerStream(m, &performanceTestServerStreamServer{stream})
}

type PerformanceTest_ServerStreamServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type performanceTestServerStreamServer struct {
	grpc.ServerStream
}

func (x *performanceTestServerStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _PerformanceTest_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PerformanceTestServer).ClientStream(&performanceTestClientStreamServer{stream})
}

type PerformanceTest_ClientStreamServer interface {
	SendAndClose(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type performanceTestClientStreamServer struct {
	grpc.ServerStream
}

func (x *performanceTestClientStreamServer) SendAndClose(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *performanceTestClientStreamServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PerformanceTest_ServiceDesc is the grpc.ServiceDesc for PerformanceTest service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PerformanceTest_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "performance.PerformanceTest",
	HandlerType: (*PerformanceTestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SimpleRPC",
			Handler:    _PerformanceTest_SimpleRPC_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamData",
			Handler:       _PerformanceTest_StreamData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ServerStream",
			Handler:       _PerformanceTest_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _PerformanceTest_ClientStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "performance.proto",
}
