// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: namespaces/service.proto

package namespaces

import (
	context "context"
	idl_common "github.com/krafton-hq/red-fox/apis/idl_common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NamespaceServerClient is the client API for NamespaceServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NamespaceServerClient interface {
	GetNamespace(ctx context.Context, in *idl_common.SingleObjectReq, opts ...grpc.CallOption) (*GetNamespaceRes, error)
	ListNamespaces(ctx context.Context, in *idl_common.ListObjectReq, opts ...grpc.CallOption) (*ListNamespacesRes, error)
	CreateNamespace(ctx context.Context, in *CreateNamespaceReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error)
	UpdateNamespace(ctx context.Context, in *UpdateNamespaceReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error)
	DeleteNamespaces(ctx context.Context, in *idl_common.SingleObjectReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error)
}

type namespaceServerClient struct {
	cc grpc.ClientConnInterface
}

func NewNamespaceServerClient(cc grpc.ClientConnInterface) NamespaceServerClient {
	return &namespaceServerClient{cc}
}

func (c *namespaceServerClient) GetNamespace(ctx context.Context, in *idl_common.SingleObjectReq, opts ...grpc.CallOption) (*GetNamespaceRes, error) {
	out := new(GetNamespaceRes)
	err := c.cc.Invoke(ctx, "/redfox.api.namespaces.NamespaceServer/GetNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceServerClient) ListNamespaces(ctx context.Context, in *idl_common.ListObjectReq, opts ...grpc.CallOption) (*ListNamespacesRes, error) {
	out := new(ListNamespacesRes)
	err := c.cc.Invoke(ctx, "/redfox.api.namespaces.NamespaceServer/ListNamespaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceServerClient) CreateNamespace(ctx context.Context, in *CreateNamespaceReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error) {
	out := new(idl_common.CommonRes)
	err := c.cc.Invoke(ctx, "/redfox.api.namespaces.NamespaceServer/CreateNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceServerClient) UpdateNamespace(ctx context.Context, in *UpdateNamespaceReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error) {
	out := new(idl_common.CommonRes)
	err := c.cc.Invoke(ctx, "/redfox.api.namespaces.NamespaceServer/UpdateNamespace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namespaceServerClient) DeleteNamespaces(ctx context.Context, in *idl_common.SingleObjectReq, opts ...grpc.CallOption) (*idl_common.CommonRes, error) {
	out := new(idl_common.CommonRes)
	err := c.cc.Invoke(ctx, "/redfox.api.namespaces.NamespaceServer/DeleteNamespaces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NamespaceServerServer is the server API for NamespaceServer service.
// All implementations must embed UnimplementedNamespaceServerServer
// for forward compatibility
type NamespaceServerServer interface {
	GetNamespace(context.Context, *idl_common.SingleObjectReq) (*GetNamespaceRes, error)
	ListNamespaces(context.Context, *idl_common.ListObjectReq) (*ListNamespacesRes, error)
	CreateNamespace(context.Context, *CreateNamespaceReq) (*idl_common.CommonRes, error)
	UpdateNamespace(context.Context, *UpdateNamespaceReq) (*idl_common.CommonRes, error)
	DeleteNamespaces(context.Context, *idl_common.SingleObjectReq) (*idl_common.CommonRes, error)
	mustEmbedUnimplementedNamespaceServerServer()
}

// UnimplementedNamespaceServerServer must be embedded to have forward compatible implementations.
type UnimplementedNamespaceServerServer struct {
}

func (UnimplementedNamespaceServerServer) GetNamespace(context.Context, *idl_common.SingleObjectReq) (*GetNamespaceRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNamespace not implemented")
}
func (UnimplementedNamespaceServerServer) ListNamespaces(context.Context, *idl_common.ListObjectReq) (*ListNamespacesRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListNamespaces not implemented")
}
func (UnimplementedNamespaceServerServer) CreateNamespace(context.Context, *CreateNamespaceReq) (*idl_common.CommonRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNamespace not implemented")
}
func (UnimplementedNamespaceServerServer) UpdateNamespace(context.Context, *UpdateNamespaceReq) (*idl_common.CommonRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNamespace not implemented")
}
func (UnimplementedNamespaceServerServer) DeleteNamespaces(context.Context, *idl_common.SingleObjectReq) (*idl_common.CommonRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNamespaces not implemented")
}
func (UnimplementedNamespaceServerServer) mustEmbedUnimplementedNamespaceServerServer() {}

// UnsafeNamespaceServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NamespaceServerServer will
// result in compilation errors.
type UnsafeNamespaceServerServer interface {
	mustEmbedUnimplementedNamespaceServerServer()
}

func RegisterNamespaceServerServer(s grpc.ServiceRegistrar, srv NamespaceServerServer) {
	s.RegisterService(&NamespaceServer_ServiceDesc, srv)
}

func _NamespaceServer_GetNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(idl_common.SingleObjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceServerServer).GetNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redfox.api.namespaces.NamespaceServer/GetNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceServerServer).GetNamespace(ctx, req.(*idl_common.SingleObjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamespaceServer_ListNamespaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(idl_common.ListObjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceServerServer).ListNamespaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redfox.api.namespaces.NamespaceServer/ListNamespaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceServerServer).ListNamespaces(ctx, req.(*idl_common.ListObjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamespaceServer_CreateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNamespaceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceServerServer).CreateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redfox.api.namespaces.NamespaceServer/CreateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceServerServer).CreateNamespace(ctx, req.(*CreateNamespaceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamespaceServer_UpdateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNamespaceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceServerServer).UpdateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redfox.api.namespaces.NamespaceServer/UpdateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceServerServer).UpdateNamespace(ctx, req.(*UpdateNamespaceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamespaceServer_DeleteNamespaces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(idl_common.SingleObjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamespaceServerServer).DeleteNamespaces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/redfox.api.namespaces.NamespaceServer/DeleteNamespaces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamespaceServerServer).DeleteNamespaces(ctx, req.(*idl_common.SingleObjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

// NamespaceServer_ServiceDesc is the grpc.ServiceDesc for NamespaceServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NamespaceServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "redfox.api.namespaces.NamespaceServer",
	HandlerType: (*NamespaceServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNamespace",
			Handler:    _NamespaceServer_GetNamespace_Handler,
		},
		{
			MethodName: "ListNamespaces",
			Handler:    _NamespaceServer_ListNamespaces_Handler,
		},
		{
			MethodName: "CreateNamespace",
			Handler:    _NamespaceServer_CreateNamespace_Handler,
		},
		{
			MethodName: "UpdateNamespace",
			Handler:    _NamespaceServer_UpdateNamespace_Handler,
		},
		{
			MethodName: "DeleteNamespaces",
			Handler:    _NamespaceServer_DeleteNamespaces_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "namespaces/service.proto",
}
