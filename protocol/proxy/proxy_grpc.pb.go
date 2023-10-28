// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.1
// source: proxy/proxy.proto

package proxypb

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

const (
	Proxy_Connect_FullMethodName      = "/loopify.proxy.Proxy/Connect"
	Proxy_Subscribe_FullMethodName    = "/loopify.proxy.Proxy/Subscribe"
	Proxy_Subscribed_FullMethodName   = "/loopify.proxy.Proxy/Subscribed"
	Proxy_Unsubscribed_FullMethodName = "/loopify.proxy.Proxy/Unsubscribed"
	Proxy_RPC_FullMethodName          = "/loopify.proxy.Proxy/RPC"
	Proxy_Publish_FullMethodName      = "/loopify.proxy.Proxy/Publish"
	Proxy_Refresh_FullMethodName      = "/loopify.proxy.Proxy/Refresh"
)

// ProxyClient is the client API for Proxy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProxyClient interface {
	// 连接认证
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectReply, error)
	// 检查是否可以订阅该频道
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeReply, error)
	// 订阅频道后的回调
	Subscribed(ctx context.Context, in *SubscribedRequest, opts ...grpc.CallOption) (*SubscribedReply, error)
	// 取消订阅频道后的回调
	Unsubscribed(ctx context.Context, in *UnsubscribedRequest, opts ...grpc.CallOption) (*UnsubscribedReply, error)
	// RPC 调用
	RPC(ctx context.Context, in *RPCRequest, opts ...grpc.CallOption) (*RPCReply, error)
	// 发布消息的回调
	Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishReply, error)
	// 刷新 Session
	Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshReply, error)
}

type proxyClient struct {
	cc grpc.ClientConnInterface
}

func NewProxyClient(cc grpc.ClientConnInterface) ProxyClient {
	return &proxyClient{cc}
}

func (c *proxyClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectReply, error) {
	out := new(ConnectReply)
	err := c.cc.Invoke(ctx, Proxy_Connect_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (*SubscribeReply, error) {
	out := new(SubscribeReply)
	err := c.cc.Invoke(ctx, Proxy_Subscribe_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Subscribed(ctx context.Context, in *SubscribedRequest, opts ...grpc.CallOption) (*SubscribedReply, error) {
	out := new(SubscribedReply)
	err := c.cc.Invoke(ctx, Proxy_Subscribed_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Unsubscribed(ctx context.Context, in *UnsubscribedRequest, opts ...grpc.CallOption) (*UnsubscribedReply, error) {
	out := new(UnsubscribedReply)
	err := c.cc.Invoke(ctx, Proxy_Unsubscribed_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) RPC(ctx context.Context, in *RPCRequest, opts ...grpc.CallOption) (*RPCReply, error) {
	out := new(RPCReply)
	err := c.cc.Invoke(ctx, Proxy_RPC_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Publish(ctx context.Context, in *PublishRequest, opts ...grpc.CallOption) (*PublishReply, error) {
	out := new(PublishReply)
	err := c.cc.Invoke(ctx, Proxy_Publish_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *proxyClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshReply, error) {
	out := new(RefreshReply)
	err := c.cc.Invoke(ctx, Proxy_Refresh_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProxyServer is the server API for Proxy service.
// All implementations must embed UnimplementedProxyServer
// for forward compatibility
type ProxyServer interface {
	// 连接认证
	Connect(context.Context, *ConnectRequest) (*ConnectReply, error)
	// 检查是否可以订阅该频道
	Subscribe(context.Context, *SubscribeRequest) (*SubscribeReply, error)
	// 订阅频道后的回调
	Subscribed(context.Context, *SubscribedRequest) (*SubscribedReply, error)
	// 取消订阅频道后的回调
	Unsubscribed(context.Context, *UnsubscribedRequest) (*UnsubscribedReply, error)
	// RPC 调用
	RPC(context.Context, *RPCRequest) (*RPCReply, error)
	// 发布消息的回调
	Publish(context.Context, *PublishRequest) (*PublishReply, error)
	// 刷新 Session
	Refresh(context.Context, *RefreshRequest) (*RefreshReply, error)
	mustEmbedUnimplementedProxyServer()
}

// UnimplementedProxyServer must be embedded to have forward compatible implementations.
type UnimplementedProxyServer struct {
}

func (UnimplementedProxyServer) Connect(context.Context, *ConnectRequest) (*ConnectReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedProxyServer) Subscribe(context.Context, *SubscribeRequest) (*SubscribeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedProxyServer) Subscribed(context.Context, *SubscribedRequest) (*SubscribedReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Subscribed not implemented")
}
func (UnimplementedProxyServer) Unsubscribed(context.Context, *UnsubscribedRequest) (*UnsubscribedReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribed not implemented")
}
func (UnimplementedProxyServer) RPC(context.Context, *RPCRequest) (*RPCReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPC not implemented")
}
func (UnimplementedProxyServer) Publish(context.Context, *PublishRequest) (*PublishReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedProxyServer) Refresh(context.Context, *RefreshRequest) (*RefreshReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}
func (UnimplementedProxyServer) mustEmbedUnimplementedProxyServer() {}

// UnsafeProxyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProxyServer will
// result in compilation errors.
type UnsafeProxyServer interface {
	mustEmbedUnimplementedProxyServer()
}

func RegisterProxyServer(s grpc.ServiceRegistrar, srv ProxyServer) {
	s.RegisterService(&Proxy_ServiceDesc, srv)
}

func _Proxy_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Connect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Subscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Subscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Subscribe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Subscribe(ctx, req.(*SubscribeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Subscribed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Subscribed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Subscribed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Subscribed(ctx, req.(*SubscribedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Unsubscribed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Unsubscribed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Unsubscribed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Unsubscribed(ctx, req.(*UnsubscribedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_RPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RPCRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).RPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_RPC_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).RPC(ctx, req.(*RPCRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Publish(ctx, req.(*PublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Proxy_Refresh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProxyServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Proxy_Refresh_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProxyServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Proxy_ServiceDesc is the grpc.ServiceDesc for Proxy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Proxy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "loopify.proxy.Proxy",
	HandlerType: (*ProxyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _Proxy_Connect_Handler,
		},
		{
			MethodName: "Subscribe",
			Handler:    _Proxy_Subscribe_Handler,
		},
		{
			MethodName: "Subscribed",
			Handler:    _Proxy_Subscribed_Handler,
		},
		{
			MethodName: "Unsubscribed",
			Handler:    _Proxy_Unsubscribed_Handler,
		},
		{
			MethodName: "RPC",
			Handler:    _Proxy_RPC_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _Proxy_Publish_Handler,
		},
		{
			MethodName: "Refresh",
			Handler:    _Proxy_Refresh_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proxy/proxy.proto",
}
