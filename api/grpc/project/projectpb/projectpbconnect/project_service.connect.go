// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: grpc/project/projectpb/project_service.proto

package projectpbconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	projectpb "github.com/loomi-labs/star-scope/grpc/project/projectpb"
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
	// ProjectServiceName is the fully-qualified name of the ProjectService service.
	ProjectServiceName = "starscope.grpc.ProjectService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ProjectServiceListProjectsProcedure is the fully-qualified name of the ProjectService's
	// ListProjects RPC.
	ProjectServiceListProjectsProcedure = "/starscope.grpc.ProjectService/ListProjects"
)

// ProjectServiceClient is a client for the starscope.grpc.ProjectService service.
type ProjectServiceClient interface {
	ListProjects(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[projectpb.ListProjectsResponse], error)
}

// NewProjectServiceClient constructs a client for the starscope.grpc.ProjectService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewProjectServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ProjectServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &projectServiceClient{
		listProjects: connect_go.NewClient[emptypb.Empty, projectpb.ListProjectsResponse](
			httpClient,
			baseURL+ProjectServiceListProjectsProcedure,
			opts...,
		),
	}
}

// projectServiceClient implements ProjectServiceClient.
type projectServiceClient struct {
	listProjects *connect_go.Client[emptypb.Empty, projectpb.ListProjectsResponse]
}

// ListProjects calls starscope.grpc.ProjectService.ListProjects.
func (c *projectServiceClient) ListProjects(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[projectpb.ListProjectsResponse], error) {
	return c.listProjects.CallUnary(ctx, req)
}

// ProjectServiceHandler is an implementation of the starscope.grpc.ProjectService service.
type ProjectServiceHandler interface {
	ListProjects(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[projectpb.ListProjectsResponse], error)
}

// NewProjectServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewProjectServiceHandler(svc ProjectServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(ProjectServiceListProjectsProcedure, connect_go.NewUnaryHandler(
		ProjectServiceListProjectsProcedure,
		svc.ListProjects,
		opts...,
	))
	return "/starscope.grpc.ProjectService/", mux
}

// UnimplementedProjectServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedProjectServiceHandler struct{}

func (UnimplementedProjectServiceHandler) ListProjects(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[projectpb.ListProjectsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.ProjectService.ListProjects is not implemented"))
}
