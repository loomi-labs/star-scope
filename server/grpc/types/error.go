package types

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UnknownErr         = status.Error(codes.Unknown, "Unknown error occurred")
	InvalidArgumentErr = status.Error(codes.InvalidArgument, "Invalid argument")
	UserNotFoundErr    = status.Errorf(codes.NotFound, "user not found")
)
