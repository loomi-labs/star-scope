package common

import (
	"errors"
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ParseDuration(input string) (*duration.Duration, error) {
	d, err := time.ParseDuration(input)
	if err != nil {
		return nil, err
	}
	return durationpb.New(d), nil
}

var layouts = []string{
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-02T15:04:05Z",
}

func ParseTime(input string) (*timestamppb.Timestamp, error) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, input)
		if err == nil {
			return timestamppb.New(t), nil
		}
	}
	return nil, errors.New("could not parse time")
}
