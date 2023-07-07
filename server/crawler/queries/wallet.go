package queries

import (
	"fmt"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/types"
	"github.com/shifty11/go-logger/log"
)

const urlAccounts = "%v/cosmos/auth/v1beta1/accounts/%v"

func doesWalletExist(restEndpoint string, address string, retries int) (bool, error) {
	url := fmt.Sprintf(urlAccounts, restEndpoint, address)
	var response types.AccountResponse
	status, err := common.GetJson(url, 0, &response)
	if err != nil {
		if status == 404 {
			return false, nil
		}
		if retries > 0 {
			log.Sugar.Errorf("Error getting wallet %v: %v", address, err)
			return doesWalletExist(restEndpoint, address, retries-1)
		}
	}
	return true, nil
}

func DoesWalletExist(restEndpoint string, address string) (bool, error) {
	return doesWalletExist(restEndpoint, address, 2)
}
