// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: grpc/user/userpb/user_service.proto

package userpbconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	userpb "github.com/loomi-labs/star-scope/grpc/user/userpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// UserServiceName is the fully-qualified name of the UserService service.
	UserServiceName = "starscope.grpc.UserService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// UserServiceGetUserProcedure is the fully-qualified name of the UserService's GetUser RPC.
	UserServiceGetUserProcedure = "/starscope.grpc.UserService/GetUser"
)

// UserServiceClient is a client for the starscope.grpc.UserService service.
type UserServiceClient interface {
	GetUser(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[userpb.User], error)
}

// NewUserServiceClient constructs a client for the starscope.grpc.UserService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewUserServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) UserServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &userServiceClient{
		getUser: connect_go.NewClient[emptypb.Empty, userpb.User](
			httpClient,
			baseURL+UserServiceGetUserProcedure,
			opts...,
		),
	}
}

// userServiceClient implements UserServiceClient.
type userServiceClient struct {
	getUser *connect_go.Client[emptypb.Empty, userpb.User]
}

// GetUser calls starscope.grpc.UserService.GetUser.
func (c *userServiceClient) GetUser(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[userpb.User], error) {
	return c.getUser.CallUnary(ctx, req)
}

// UserServiceHandler is an implementation of the starscope.grpc.UserService service.
type UserServiceHandler interface {
	GetUser(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[userpb.User], error)
}

// NewUserServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewUserServiceHandler(svc UserServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(UserServiceGetUserProcedure, connect_go.NewUnaryHandler(
		UserServiceGetUserProcedure,
		svc.GetUser,
		opts...,
	))
	return "/starscope.grpc.UserService/", mux
}

// UnimplementedUserServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedUserServiceHandler struct{}

func (UnimplementedUserServiceHandler) GetUser(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[userpb.User], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.UserService.GetUser is not implemented"))
}
