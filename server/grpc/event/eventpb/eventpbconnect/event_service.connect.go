// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: grpc/event/eventpb/event_service.proto

package eventpbconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	eventpb "github.com/loomi-labs/star-scope/grpc/event/eventpb"
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
	// EventServiceName is the fully-qualified name of the EventService service.
	EventServiceName = "starscope.grpc.event.EventService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// EventServiceEventStreamProcedure is the fully-qualified name of the EventService's EventStream
	// RPC.
	EventServiceEventStreamProcedure = "/starscope.grpc.event.EventService/EventStream"
	// EventServiceListEventsProcedure is the fully-qualified name of the EventService's ListEvents RPC.
	EventServiceListEventsProcedure = "/starscope.grpc.event.EventService/ListEvents"
	// EventServiceListChainsProcedure is the fully-qualified name of the EventService's ListChains RPC.
	EventServiceListChainsProcedure = "/starscope.grpc.event.EventService/ListChains"
	// EventServiceListEventsCountProcedure is the fully-qualified name of the EventService's
	// ListEventsCount RPC.
	EventServiceListEventsCountProcedure = "/starscope.grpc.event.EventService/ListEventsCount"
	// EventServiceMarkEventReadProcedure is the fully-qualified name of the EventService's
	// MarkEventRead RPC.
	EventServiceMarkEventReadProcedure = "/starscope.grpc.event.EventService/MarkEventRead"
	// EventServiceGetWelcomeMessageProcedure is the fully-qualified name of the EventService's
	// GetWelcomeMessage RPC.
	EventServiceGetWelcomeMessageProcedure = "/starscope.grpc.event.EventService/GetWelcomeMessage"
)

// EventServiceClient is a client for the starscope.grpc.event.EventService service.
type EventServiceClient interface {
	EventStream(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.ServerStreamForClient[eventpb.NewEvent], error)
	ListEvents(context.Context, *connect_go.Request[eventpb.ListEventsRequest]) (*connect_go.Response[eventpb.EventList], error)
	ListChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ChainList], error)
	ListEventsCount(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ListEventsCountResponse], error)
	MarkEventRead(context.Context, *connect_go.Request[eventpb.MarkEventReadRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetWelcomeMessage(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.WelcomeMessageResponse], error)
}

// NewEventServiceClient constructs a client for the starscope.grpc.event.EventService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewEventServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) EventServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &eventServiceClient{
		eventStream: connect_go.NewClient[emptypb.Empty, eventpb.NewEvent](
			httpClient,
			baseURL+EventServiceEventStreamProcedure,
			opts...,
		),
		listEvents: connect_go.NewClient[eventpb.ListEventsRequest, eventpb.EventList](
			httpClient,
			baseURL+EventServiceListEventsProcedure,
			opts...,
		),
		listChains: connect_go.NewClient[emptypb.Empty, eventpb.ChainList](
			httpClient,
			baseURL+EventServiceListChainsProcedure,
			opts...,
		),
		listEventsCount: connect_go.NewClient[emptypb.Empty, eventpb.ListEventsCountResponse](
			httpClient,
			baseURL+EventServiceListEventsCountProcedure,
			opts...,
		),
		markEventRead: connect_go.NewClient[eventpb.MarkEventReadRequest, emptypb.Empty](
			httpClient,
			baseURL+EventServiceMarkEventReadProcedure,
			opts...,
		),
		getWelcomeMessage: connect_go.NewClient[emptypb.Empty, eventpb.WelcomeMessageResponse](
			httpClient,
			baseURL+EventServiceGetWelcomeMessageProcedure,
			opts...,
		),
	}
}

// eventServiceClient implements EventServiceClient.
type eventServiceClient struct {
	eventStream       *connect_go.Client[emptypb.Empty, eventpb.NewEvent]
	listEvents        *connect_go.Client[eventpb.ListEventsRequest, eventpb.EventList]
	listChains        *connect_go.Client[emptypb.Empty, eventpb.ChainList]
	listEventsCount   *connect_go.Client[emptypb.Empty, eventpb.ListEventsCountResponse]
	markEventRead     *connect_go.Client[eventpb.MarkEventReadRequest, emptypb.Empty]
	getWelcomeMessage *connect_go.Client[emptypb.Empty, eventpb.WelcomeMessageResponse]
}

// EventStream calls starscope.grpc.event.EventService.EventStream.
func (c *eventServiceClient) EventStream(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.ServerStreamForClient[eventpb.NewEvent], error) {
	return c.eventStream.CallServerStream(ctx, req)
}

// ListEvents calls starscope.grpc.event.EventService.ListEvents.
func (c *eventServiceClient) ListEvents(ctx context.Context, req *connect_go.Request[eventpb.ListEventsRequest]) (*connect_go.Response[eventpb.EventList], error) {
	return c.listEvents.CallUnary(ctx, req)
}

// ListChains calls starscope.grpc.event.EventService.ListChains.
func (c *eventServiceClient) ListChains(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ChainList], error) {
	return c.listChains.CallUnary(ctx, req)
}

// ListEventsCount calls starscope.grpc.event.EventService.ListEventsCount.
func (c *eventServiceClient) ListEventsCount(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ListEventsCountResponse], error) {
	return c.listEventsCount.CallUnary(ctx, req)
}

// MarkEventRead calls starscope.grpc.event.EventService.MarkEventRead.
func (c *eventServiceClient) MarkEventRead(ctx context.Context, req *connect_go.Request[eventpb.MarkEventReadRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.markEventRead.CallUnary(ctx, req)
}

// GetWelcomeMessage calls starscope.grpc.event.EventService.GetWelcomeMessage.
func (c *eventServiceClient) GetWelcomeMessage(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.WelcomeMessageResponse], error) {
	return c.getWelcomeMessage.CallUnary(ctx, req)
}

// EventServiceHandler is an implementation of the starscope.grpc.event.EventService service.
type EventServiceHandler interface {
	EventStream(context.Context, *connect_go.Request[emptypb.Empty], *connect_go.ServerStream[eventpb.NewEvent]) error
	ListEvents(context.Context, *connect_go.Request[eventpb.ListEventsRequest]) (*connect_go.Response[eventpb.EventList], error)
	ListChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ChainList], error)
	ListEventsCount(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ListEventsCountResponse], error)
	MarkEventRead(context.Context, *connect_go.Request[eventpb.MarkEventReadRequest]) (*connect_go.Response[emptypb.Empty], error)
	GetWelcomeMessage(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.WelcomeMessageResponse], error)
}

// NewEventServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewEventServiceHandler(svc EventServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	eventServiceEventStreamHandler := connect_go.NewServerStreamHandler(
		EventServiceEventStreamProcedure,
		svc.EventStream,
		opts...,
	)
	eventServiceListEventsHandler := connect_go.NewUnaryHandler(
		EventServiceListEventsProcedure,
		svc.ListEvents,
		opts...,
	)
	eventServiceListChainsHandler := connect_go.NewUnaryHandler(
		EventServiceListChainsProcedure,
		svc.ListChains,
		opts...,
	)
	eventServiceListEventsCountHandler := connect_go.NewUnaryHandler(
		EventServiceListEventsCountProcedure,
		svc.ListEventsCount,
		opts...,
	)
	eventServiceMarkEventReadHandler := connect_go.NewUnaryHandler(
		EventServiceMarkEventReadProcedure,
		svc.MarkEventRead,
		opts...,
	)
	eventServiceGetWelcomeMessageHandler := connect_go.NewUnaryHandler(
		EventServiceGetWelcomeMessageProcedure,
		svc.GetWelcomeMessage,
		opts...,
	)
	return "/starscope.grpc.event.EventService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case EventServiceEventStreamProcedure:
			eventServiceEventStreamHandler.ServeHTTP(w, r)
		case EventServiceListEventsProcedure:
			eventServiceListEventsHandler.ServeHTTP(w, r)
		case EventServiceListChainsProcedure:
			eventServiceListChainsHandler.ServeHTTP(w, r)
		case EventServiceListEventsCountProcedure:
			eventServiceListEventsCountHandler.ServeHTTP(w, r)
		case EventServiceMarkEventReadProcedure:
			eventServiceMarkEventReadHandler.ServeHTTP(w, r)
		case EventServiceGetWelcomeMessageProcedure:
			eventServiceGetWelcomeMessageHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedEventServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedEventServiceHandler struct{}

func (UnimplementedEventServiceHandler) EventStream(context.Context, *connect_go.Request[emptypb.Empty], *connect_go.ServerStream[eventpb.NewEvent]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.EventStream is not implemented"))
}

func (UnimplementedEventServiceHandler) ListEvents(context.Context, *connect_go.Request[eventpb.ListEventsRequest]) (*connect_go.Response[eventpb.EventList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.ListEvents is not implemented"))
}

func (UnimplementedEventServiceHandler) ListChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ChainList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.ListChains is not implemented"))
}

func (UnimplementedEventServiceHandler) ListEventsCount(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.ListEventsCountResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.ListEventsCount is not implemented"))
}

func (UnimplementedEventServiceHandler) MarkEventRead(context.Context, *connect_go.Request[eventpb.MarkEventReadRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.MarkEventRead is not implemented"))
}

func (UnimplementedEventServiceHandler) GetWelcomeMessage(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[eventpb.WelcomeMessageResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.event.EventService.GetWelcomeMessage is not implemented"))
}
