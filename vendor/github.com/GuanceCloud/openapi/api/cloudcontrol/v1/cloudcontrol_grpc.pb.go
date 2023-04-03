// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: api/cloudcontrol/v1/cloudcontrol.proto

package v1

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
	CloudControlService_CreateResource_FullMethodName           = "/api.cloudcontrol.v1.CloudControlService/CreateResource"
	CloudControlService_GetResource_FullMethodName              = "/api.cloudcontrol.v1.CloudControlService/GetResource"
	CloudControlService_DeleteResource_FullMethodName           = "/api.cloudcontrol.v1.CloudControlService/DeleteResource"
	CloudControlService_ListResources_FullMethodName            = "/api.cloudcontrol.v1.CloudControlService/ListResources"
	CloudControlService_UpdateResource_FullMethodName           = "/api.cloudcontrol.v1.CloudControlService/UpdateResource"
	CloudControlService_CancelResourceRequest_FullMethodName    = "/api.cloudcontrol.v1.CloudControlService/CancelResourceRequest"
	CloudControlService_GetResourceRequestStatus_FullMethodName = "/api.cloudcontrol.v1.CloudControlService/GetResourceRequestStatus"
	CloudControlService_ListResourceRequests_FullMethodName     = "/api.cloudcontrol.v1.CloudControlService/ListResourceRequests"
)

// CloudControlServiceClient is the client API for CloudControlService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudControlServiceClient interface {
	// Creates a new resource.
	CreateResource(ctx context.Context, in *CreateResourceRequest, opts ...grpc.CallOption) (*CreateResourceResponse, error)
	// Gets the resource by name.
	GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error)
	// Deletes the resource by name.
	DeleteResource(ctx context.Context, in *DeleteResourceRequest, opts ...grpc.CallOption) (*DeleteResourceResponse, error)
	// Lists all resources.
	ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesResponse, error)
	// Updates the resource.
	UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*UpdateResourceResponse, error)
	// Cancels the specified resource request.
	CancelResourceRequest(ctx context.Context, in *CancelResourceRequestRequest, opts ...grpc.CallOption) (*CancelResourceRequestResponse, error)
	// Gets the specified resource request status.
	GetResourceRequestStatus(ctx context.Context, in *GetResourceRequestStatusRequest, opts ...grpc.CallOption) (*GetResourceRequestStatusResponse, error)
	// List resource requests.
	ListResourceRequests(ctx context.Context, in *ListResourceRequestsRequest, opts ...grpc.CallOption) (*ListResourceRequestsResponse, error)
}

type cloudControlServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudControlServiceClient(cc grpc.ClientConnInterface) CloudControlServiceClient {
	return &cloudControlServiceClient{cc}
}

func (c *cloudControlServiceClient) CreateResource(ctx context.Context, in *CreateResourceRequest, opts ...grpc.CallOption) (*CreateResourceResponse, error) {
	out := new(CreateResourceResponse)
	err := c.cc.Invoke(ctx, CloudControlService_CreateResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) GetResource(ctx context.Context, in *GetResourceRequest, opts ...grpc.CallOption) (*GetResourceResponse, error) {
	out := new(GetResourceResponse)
	err := c.cc.Invoke(ctx, CloudControlService_GetResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) DeleteResource(ctx context.Context, in *DeleteResourceRequest, opts ...grpc.CallOption) (*DeleteResourceResponse, error) {
	out := new(DeleteResourceResponse)
	err := c.cc.Invoke(ctx, CloudControlService_DeleteResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) ListResources(ctx context.Context, in *ListResourcesRequest, opts ...grpc.CallOption) (*ListResourcesResponse, error) {
	out := new(ListResourcesResponse)
	err := c.cc.Invoke(ctx, CloudControlService_ListResources_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) UpdateResource(ctx context.Context, in *UpdateResourceRequest, opts ...grpc.CallOption) (*UpdateResourceResponse, error) {
	out := new(UpdateResourceResponse)
	err := c.cc.Invoke(ctx, CloudControlService_UpdateResource_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) CancelResourceRequest(ctx context.Context, in *CancelResourceRequestRequest, opts ...grpc.CallOption) (*CancelResourceRequestResponse, error) {
	out := new(CancelResourceRequestResponse)
	err := c.cc.Invoke(ctx, CloudControlService_CancelResourceRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) GetResourceRequestStatus(ctx context.Context, in *GetResourceRequestStatusRequest, opts ...grpc.CallOption) (*GetResourceRequestStatusResponse, error) {
	out := new(GetResourceRequestStatusResponse)
	err := c.cc.Invoke(ctx, CloudControlService_GetResourceRequestStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudControlServiceClient) ListResourceRequests(ctx context.Context, in *ListResourceRequestsRequest, opts ...grpc.CallOption) (*ListResourceRequestsResponse, error) {
	out := new(ListResourceRequestsResponse)
	err := c.cc.Invoke(ctx, CloudControlService_ListResourceRequests_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudControlServiceServer is the server API for CloudControlService service.
// All implementations must embed UnimplementedCloudControlServiceServer
// for forward compatibility
type CloudControlServiceServer interface {
	// Creates a new resource.
	CreateResource(context.Context, *CreateResourceRequest) (*CreateResourceResponse, error)
	// Gets the resource by name.
	GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error)
	// Deletes the resource by name.
	DeleteResource(context.Context, *DeleteResourceRequest) (*DeleteResourceResponse, error)
	// Lists all resources.
	ListResources(context.Context, *ListResourcesRequest) (*ListResourcesResponse, error)
	// Updates the resource.
	UpdateResource(context.Context, *UpdateResourceRequest) (*UpdateResourceResponse, error)
	// Cancels the specified resource request.
	CancelResourceRequest(context.Context, *CancelResourceRequestRequest) (*CancelResourceRequestResponse, error)
	// Gets the specified resource request status.
	GetResourceRequestStatus(context.Context, *GetResourceRequestStatusRequest) (*GetResourceRequestStatusResponse, error)
	// List resource requests.
	ListResourceRequests(context.Context, *ListResourceRequestsRequest) (*ListResourceRequestsResponse, error)
	mustEmbedUnimplementedCloudControlServiceServer()
}

// UnimplementedCloudControlServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCloudControlServiceServer struct {
}

func (UnimplementedCloudControlServiceServer) CreateResource(context.Context, *CreateResourceRequest) (*CreateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateResource not implemented")
}
func (UnimplementedCloudControlServiceServer) GetResource(context.Context, *GetResourceRequest) (*GetResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResource not implemented")
}
func (UnimplementedCloudControlServiceServer) DeleteResource(context.Context, *DeleteResourceRequest) (*DeleteResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteResource not implemented")
}
func (UnimplementedCloudControlServiceServer) ListResources(context.Context, *ListResourcesRequest) (*ListResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResources not implemented")
}
func (UnimplementedCloudControlServiceServer) UpdateResource(context.Context, *UpdateResourceRequest) (*UpdateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResource not implemented")
}
func (UnimplementedCloudControlServiceServer) CancelResourceRequest(context.Context, *CancelResourceRequestRequest) (*CancelResourceRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelResourceRequest not implemented")
}
func (UnimplementedCloudControlServiceServer) GetResourceRequestStatus(context.Context, *GetResourceRequestStatusRequest) (*GetResourceRequestStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResourceRequestStatus not implemented")
}
func (UnimplementedCloudControlServiceServer) ListResourceRequests(context.Context, *ListResourceRequestsRequest) (*ListResourceRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListResourceRequests not implemented")
}
func (UnimplementedCloudControlServiceServer) mustEmbedUnimplementedCloudControlServiceServer() {}

// UnsafeCloudControlServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudControlServiceServer will
// result in compilation errors.
type UnsafeCloudControlServiceServer interface {
	mustEmbedUnimplementedCloudControlServiceServer()
}

func RegisterCloudControlServiceServer(s grpc.ServiceRegistrar, srv CloudControlServiceServer) {
	s.RegisterService(&CloudControlService_ServiceDesc, srv)
}

func _CloudControlService_CreateResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).CreateResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_CreateResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).CreateResource(ctx, req.(*CreateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_GetResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).GetResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_GetResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).GetResource(ctx, req.(*GetResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_DeleteResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).DeleteResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_DeleteResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).DeleteResource(ctx, req.(*DeleteResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_ListResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).ListResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_ListResources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).ListResources(ctx, req.(*ListResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_UpdateResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).UpdateResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_UpdateResource_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).UpdateResource(ctx, req.(*UpdateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_CancelResourceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelResourceRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).CancelResourceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_CancelResourceRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).CancelResourceRequest(ctx, req.(*CancelResourceRequestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_GetResourceRequestStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResourceRequestStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).GetResourceRequestStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_GetResourceRequestStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).GetResourceRequestStatus(ctx, req.(*GetResourceRequestStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudControlService_ListResourceRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListResourceRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudControlServiceServer).ListResourceRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudControlService_ListResourceRequests_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudControlServiceServer).ListResourceRequests(ctx, req.(*ListResourceRequestsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CloudControlService_ServiceDesc is the grpc.ServiceDesc for CloudControlService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloudControlService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.cloudcontrol.v1.CloudControlService",
	HandlerType: (*CloudControlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateResource",
			Handler:    _CloudControlService_CreateResource_Handler,
		},
		{
			MethodName: "GetResource",
			Handler:    _CloudControlService_GetResource_Handler,
		},
		{
			MethodName: "DeleteResource",
			Handler:    _CloudControlService_DeleteResource_Handler,
		},
		{
			MethodName: "ListResources",
			Handler:    _CloudControlService_ListResources_Handler,
		},
		{
			MethodName: "UpdateResource",
			Handler:    _CloudControlService_UpdateResource_Handler,
		},
		{
			MethodName: "CancelResourceRequest",
			Handler:    _CloudControlService_CancelResourceRequest_Handler,
		},
		{
			MethodName: "GetResourceRequestStatus",
			Handler:    _CloudControlService_GetResourceRequestStatus_Handler,
		},
		{
			MethodName: "ListResourceRequests",
			Handler:    _CloudControlService_ListResourceRequests_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/cloudcontrol/v1/cloudcontrol.proto",
}
