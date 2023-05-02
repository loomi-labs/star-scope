package indexer

import (
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func parseDuration(input string) (*duration.Duration, error) {
	d, err := time.ParseDuration(input)
	if err != nil {
		return nil, err
	}
	return durationpb.New(d), nil
}

func parseTime(input string) (*timestamppb.Timestamp, error) {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", input)
	if err != nil {
		return nil, err
	}
	return timestamppb.New(t), nil
}
