// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gdean

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

// IssueServiceClient is the client API for IssueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IssueServiceClient interface {
	CreateIssue(ctx context.Context, in *CreateIssueRequest, opts ...grpc.CallOption) (*CreateIssueResponse, error)
	PutIssue(ctx context.Context, in *PutIssueRequest, opts ...grpc.CallOption) (*PutIssueResponse, error)
	GetIssue(ctx context.Context, in *GetIssueRequest, opts ...grpc.CallOption) (*GetIssueResponse, error)
	DeleteIssue(ctx context.Context, in *DeleteIssueRequest, opts ...grpc.CallOption) (*DeleteIssueResponse, error)
	CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
	PutComment(ctx context.Context, in *PutCommentRequest, opts ...grpc.CallOption) (*PutCommentResponse, error)
	GetComment(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error)
	DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*DeleteCommentResponse, error)
	CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error)
	PutData(ctx context.Context, in *PutDataRequest, opts ...grpc.CallOption) (*PutDataResponse, error)
	GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResponse, error)
	DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error)
	DecideBranch(ctx context.Context, in *DecideBranchRequest, opts ...grpc.CallOption) (*DecideBranchResponse, error)
}

type issueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIssueServiceClient(cc grpc.ClientConnInterface) IssueServiceClient {
	return &issueServiceClient{cc}
}

func (c *issueServiceClient) CreateIssue(ctx context.Context, in *CreateIssueRequest, opts ...grpc.CallOption) (*CreateIssueResponse, error) {
	out := new(CreateIssueResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/CreateIssue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) PutIssue(ctx context.Context, in *PutIssueRequest, opts ...grpc.CallOption) (*PutIssueResponse, error) {
	out := new(PutIssueResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/PutIssue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) GetIssue(ctx context.Context, in *GetIssueRequest, opts ...grpc.CallOption) (*GetIssueResponse, error) {
	out := new(GetIssueResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/GetIssue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) DeleteIssue(ctx context.Context, in *DeleteIssueRequest, opts ...grpc.CallOption) (*DeleteIssueResponse, error) {
	out := new(DeleteIssueResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/DeleteIssue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/CreateComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) PutComment(ctx context.Context, in *PutCommentRequest, opts ...grpc.CallOption) (*PutCommentResponse, error) {
	out := new(PutCommentResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/PutComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) GetComment(ctx context.Context, in *GetCommentRequest, opts ...grpc.CallOption) (*GetCommentResponse, error) {
	out := new(GetCommentResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/GetComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*DeleteCommentResponse, error) {
	out := new(DeleteCommentResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/DeleteComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error) {
	out := new(CreateDataResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/CreateData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) PutData(ctx context.Context, in *PutDataRequest, opts ...grpc.CallOption) (*PutDataResponse, error) {
	out := new(PutDataResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/PutData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) GetData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*GetDataResponse, error) {
	out := new(GetDataResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error) {
	out := new(DeleteDataResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/DeleteData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *issueServiceClient) DecideBranch(ctx context.Context, in *DecideBranchRequest, opts ...grpc.CallOption) (*DecideBranchResponse, error) {
	out := new(DecideBranchResponse)
	err := c.cc.Invoke(ctx, "/gdean.IssueService/DecideBranch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IssueServiceServer is the server API for IssueService service.
// All implementations must embed UnimplementedIssueServiceServer
// for forward compatibility
type IssueServiceServer interface {
	CreateIssue(context.Context, *CreateIssueRequest) (*CreateIssueResponse, error)
	PutIssue(context.Context, *PutIssueRequest) (*PutIssueResponse, error)
	GetIssue(context.Context, *GetIssueRequest) (*GetIssueResponse, error)
	DeleteIssue(context.Context, *DeleteIssueRequest) (*DeleteIssueResponse, error)
	CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
	PutComment(context.Context, *PutCommentRequest) (*PutCommentResponse, error)
	GetComment(context.Context, *GetCommentRequest) (*GetCommentResponse, error)
	DeleteComment(context.Context, *DeleteCommentRequest) (*DeleteCommentResponse, error)
	CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error)
	PutData(context.Context, *PutDataRequest) (*PutDataResponse, error)
	GetData(context.Context, *GetDataRequest) (*GetDataResponse, error)
	DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error)
	DecideBranch(context.Context, *DecideBranchRequest) (*DecideBranchResponse, error)
	mustEmbedUnimplementedIssueServiceServer()
}

// UnimplementedIssueServiceServer must be embedded to have forward compatible implementations.
type UnimplementedIssueServiceServer struct {
}

func (UnimplementedIssueServiceServer) CreateIssue(context.Context, *CreateIssueRequest) (*CreateIssueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIssue not implemented")
}
func (UnimplementedIssueServiceServer) PutIssue(context.Context, *PutIssueRequest) (*PutIssueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutIssue not implemented")
}
func (UnimplementedIssueServiceServer) GetIssue(context.Context, *GetIssueRequest) (*GetIssueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIssue not implemented")
}
func (UnimplementedIssueServiceServer) DeleteIssue(context.Context, *DeleteIssueRequest) (*DeleteIssueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIssue not implemented")
}
func (UnimplementedIssueServiceServer) CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedIssueServiceServer) PutComment(context.Context, *PutCommentRequest) (*PutCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutComment not implemented")
}
func (UnimplementedIssueServiceServer) GetComment(context.Context, *GetCommentRequest) (*GetCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComment not implemented")
}
func (UnimplementedIssueServiceServer) DeleteComment(context.Context, *DeleteCommentRequest) (*DeleteCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteComment not implemented")
}
func (UnimplementedIssueServiceServer) CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateData not implemented")
}
func (UnimplementedIssueServiceServer) PutData(context.Context, *PutDataRequest) (*PutDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutData not implemented")
}
func (UnimplementedIssueServiceServer) GetData(context.Context, *GetDataRequest) (*GetDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetData not implemented")
}
func (UnimplementedIssueServiceServer) DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteData not implemented")
}
func (UnimplementedIssueServiceServer) DecideBranch(context.Context, *DecideBranchRequest) (*DecideBranchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecideBranch not implemented")
}
func (UnimplementedIssueServiceServer) mustEmbedUnimplementedIssueServiceServer() {}

// UnsafeIssueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IssueServiceServer will
// result in compilation errors.
type UnsafeIssueServiceServer interface {
	mustEmbedUnimplementedIssueServiceServer()
}

func RegisterIssueServiceServer(s grpc.ServiceRegistrar, srv IssueServiceServer) {
	s.RegisterService(&IssueService_ServiceDesc, srv)
}

func _IssueService_CreateIssue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateIssueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).CreateIssue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/CreateIssue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).CreateIssue(ctx, req.(*CreateIssueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_PutIssue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutIssueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).PutIssue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/PutIssue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).PutIssue(ctx, req.(*PutIssueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_GetIssue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIssueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).GetIssue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/GetIssue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).GetIssue(ctx, req.(*GetIssueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_DeleteIssue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteIssueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).DeleteIssue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/DeleteIssue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).DeleteIssue(ctx, req.(*DeleteIssueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/CreateComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).CreateComment(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_PutComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).PutComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/PutComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).PutComment(ctx, req.(*PutCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_GetComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).GetComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/GetComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).GetComment(ctx, req.(*GetCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_DeleteComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).DeleteComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/DeleteComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).DeleteComment(ctx, req.(*DeleteCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_CreateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).CreateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/CreateData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).CreateData(ctx, req.(*CreateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_PutData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).PutData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/PutData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).PutData(ctx, req.(*PutDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).GetData(ctx, req.(*GetDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_DeleteData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).DeleteData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/DeleteData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).DeleteData(ctx, req.(*DeleteDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IssueService_DecideBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecideBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IssueServiceServer).DecideBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gdean.IssueService/DecideBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IssueServiceServer).DecideBranch(ctx, req.(*DecideBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IssueService_ServiceDesc is the grpc.ServiceDesc for IssueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IssueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gdean.IssueService",
	HandlerType: (*IssueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIssue",
			Handler:    _IssueService_CreateIssue_Handler,
		},
		{
			MethodName: "PutIssue",
			Handler:    _IssueService_PutIssue_Handler,
		},
		{
			MethodName: "GetIssue",
			Handler:    _IssueService_GetIssue_Handler,
		},
		{
			MethodName: "DeleteIssue",
			Handler:    _IssueService_DeleteIssue_Handler,
		},
		{
			MethodName: "CreateComment",
			Handler:    _IssueService_CreateComment_Handler,
		},
		{
			MethodName: "PutComment",
			Handler:    _IssueService_PutComment_Handler,
		},
		{
			MethodName: "GetComment",
			Handler:    _IssueService_GetComment_Handler,
		},
		{
			MethodName: "DeleteComment",
			Handler:    _IssueService_DeleteComment_Handler,
		},
		{
			MethodName: "CreateData",
			Handler:    _IssueService_CreateData_Handler,
		},
		{
			MethodName: "PutData",
			Handler:    _IssueService_PutData_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _IssueService_GetData_Handler,
		},
		{
			MethodName: "DeleteData",
			Handler:    _IssueService_DeleteData_Handler,
		},
		{
			MethodName: "DecideBranch",
			Handler:    _IssueService_DecideBranch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gdeanissue.proto",
}