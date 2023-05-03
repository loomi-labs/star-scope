package indexer

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/osmosis-labs/osmosis/osmoutils/noapptest"
	"io"
	"net/http"
)

func GetAndDecode(url string, encodingConfig noapptest.TestEncodingConfig, target proto.Message) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 503, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, errors.New(resp.Status)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, encodingConfig.Codec.UnmarshalJSON(body, target)
}
