// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: grpc/indexer/indexerpb/indexer_service.proto

package indexerpbconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	indexerpb "github.com/loomi-labs/star-scope/grpc/indexer/indexerpb"
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
	// IndexerServiceName is the fully-qualified name of the IndexerService service.
	IndexerServiceName = "starscope.grpc.IndexerService"
	// TxHandlerServiceName is the fully-qualified name of the TxHandlerService service.
	TxHandlerServiceName = "starscope.grpc.TxHandlerService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// IndexerServiceGetIndexingChainsProcedure is the fully-qualified name of the IndexerService's
	// GetIndexingChains RPC.
	IndexerServiceGetIndexingChainsProcedure = "/starscope.grpc.IndexerService/GetIndexingChains"
	// IndexerServiceUpdateIndexingChainsProcedure is the fully-qualified name of the IndexerService's
	// UpdateIndexingChains RPC.
	IndexerServiceUpdateIndexingChainsProcedure = "/starscope.grpc.IndexerService/UpdateIndexingChains"
	// TxHandlerServiceHandleTxsProcedure is the fully-qualified name of the TxHandlerService's
	// HandleTxs RPC.
	TxHandlerServiceHandleTxsProcedure = "/starscope.grpc.TxHandlerService/HandleTxs"
)

// IndexerServiceClient is a client for the starscope.grpc.IndexerService service.
type IndexerServiceClient interface {
	GetIndexingChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[indexerpb.GetIndexingChainsResponse], error)
	UpdateIndexingChains(context.Context, *connect_go.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect_go.Response[indexerpb.UpdateIndexingChainsResponse], error)
}

// NewIndexerServiceClient constructs a client for the starscope.grpc.IndexerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewIndexerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) IndexerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &indexerServiceClient{
		getIndexingChains: connect_go.NewClient[emptypb.Empty, indexerpb.GetIndexingChainsResponse](
			httpClient,
			baseURL+IndexerServiceGetIndexingChainsProcedure,
			opts...,
		),
		updateIndexingChains: connect_go.NewClient[indexerpb.UpdateIndexingChainsRequest, indexerpb.UpdateIndexingChainsResponse](
			httpClient,
			baseURL+IndexerServiceUpdateIndexingChainsProcedure,
			opts...,
		),
	}
}

// indexerServiceClient implements IndexerServiceClient.
type indexerServiceClient struct {
	getIndexingChains    *connect_go.Client[emptypb.Empty, indexerpb.GetIndexingChainsResponse]
	updateIndexingChains *connect_go.Client[indexerpb.UpdateIndexingChainsRequest, indexerpb.UpdateIndexingChainsResponse]
}

// GetIndexingChains calls starscope.grpc.IndexerService.GetIndexingChains.
func (c *indexerServiceClient) GetIndexingChains(ctx context.Context, req *connect_go.Request[emptypb.Empty]) (*connect_go.Response[indexerpb.GetIndexingChainsResponse], error) {
	return c.getIndexingChains.CallUnary(ctx, req)
}

// UpdateIndexingChains calls starscope.grpc.IndexerService.UpdateIndexingChains.
func (c *indexerServiceClient) UpdateIndexingChains(ctx context.Context, req *connect_go.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect_go.Response[indexerpb.UpdateIndexingChainsResponse], error) {
	return c.updateIndexingChains.CallUnary(ctx, req)
}

// IndexerServiceHandler is an implementation of the starscope.grpc.IndexerService service.
type IndexerServiceHandler interface {
	GetIndexingChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[indexerpb.GetIndexingChainsResponse], error)
	UpdateIndexingChains(context.Context, *connect_go.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect_go.Response[indexerpb.UpdateIndexingChainsResponse], error)
}

// NewIndexerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewIndexerServiceHandler(svc IndexerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	indexerServiceGetIndexingChainsHandler := connect_go.NewUnaryHandler(
		IndexerServiceGetIndexingChainsProcedure,
		svc.GetIndexingChains,
		opts...,
	)
	indexerServiceUpdateIndexingChainsHandler := connect_go.NewUnaryHandler(
		IndexerServiceUpdateIndexingChainsProcedure,
		svc.UpdateIndexingChains,
		opts...,
	)
	return "/starscope.grpc.IndexerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case IndexerServiceGetIndexingChainsProcedure:
			indexerServiceGetIndexingChainsHandler.ServeHTTP(w, r)
		case IndexerServiceUpdateIndexingChainsProcedure:
			indexerServiceUpdateIndexingChainsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedIndexerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedIndexerServiceHandler struct{}

func (UnimplementedIndexerServiceHandler) GetIndexingChains(context.Context, *connect_go.Request[emptypb.Empty]) (*connect_go.Response[indexerpb.GetIndexingChainsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.IndexerService.GetIndexingChains is not implemented"))
}

func (UnimplementedIndexerServiceHandler) UpdateIndexingChains(context.Context, *connect_go.Request[indexerpb.UpdateIndexingChainsRequest]) (*connect_go.Response[indexerpb.UpdateIndexingChainsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.IndexerService.UpdateIndexingChains is not implemented"))
}

// TxHandlerServiceClient is a client for the starscope.grpc.TxHandlerService service.
type TxHandlerServiceClient interface {
	HandleTxs(context.Context, *connect_go.Request[indexerpb.HandleTxsRequest]) (*connect_go.Response[indexerpb.HandleTxsResponse], error)
}

// NewTxHandlerServiceClient constructs a client for the starscope.grpc.TxHandlerService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewTxHandlerServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) TxHandlerServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &txHandlerServiceClient{
		handleTxs: connect_go.NewClient[indexerpb.HandleTxsRequest, indexerpb.HandleTxsResponse](
			httpClient,
			baseURL+TxHandlerServiceHandleTxsProcedure,
			opts...,
		),
	}
}

// txHandlerServiceClient implements TxHandlerServiceClient.
type txHandlerServiceClient struct {
	handleTxs *connect_go.Client[indexerpb.HandleTxsRequest, indexerpb.HandleTxsResponse]
}

// HandleTxs calls starscope.grpc.TxHandlerService.HandleTxs.
func (c *txHandlerServiceClient) HandleTxs(ctx context.Context, req *connect_go.Request[indexerpb.HandleTxsRequest]) (*connect_go.Response[indexerpb.HandleTxsResponse], error) {
	return c.handleTxs.CallUnary(ctx, req)
}

// TxHandlerServiceHandler is an implementation of the starscope.grpc.TxHandlerService service.
type TxHandlerServiceHandler interface {
	HandleTxs(context.Context, *connect_go.Request[indexerpb.HandleTxsRequest]) (*connect_go.Response[indexerpb.HandleTxsResponse], error)
}

// NewTxHandlerServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewTxHandlerServiceHandler(svc TxHandlerServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	txHandlerServiceHandleTxsHandler := connect_go.NewUnaryHandler(
		TxHandlerServiceHandleTxsProcedure,
		svc.HandleTxs,
		opts...,
	)
	return "/starscope.grpc.TxHandlerService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case TxHandlerServiceHandleTxsProcedure:
			txHandlerServiceHandleTxsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedTxHandlerServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedTxHandlerServiceHandler struct{}

func (UnimplementedTxHandlerServiceHandler) HandleTxs(context.Context, *connect_go.Request[indexerpb.HandleTxsRequest]) (*connect_go.Response[indexerpb.HandleTxsResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("starscope.grpc.TxHandlerService.HandleTxs is not implemented"))
}
