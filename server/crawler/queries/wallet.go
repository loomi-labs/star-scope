package queries

import (
	"fmt"
	"github.com/loomi-labs/star-scope/common"
	"github.com/loomi-labs/star-scope/types"
	"time"
)

const urlAccounts = "%v/cosmos/auth/v1beta1/accounts/%v"

func doesWalletExist(restEndpoint string, address string, retries int) (bool, error) {
	url := fmt.Sprintf(urlAccounts, restEndpoint, address)
	var response types.AccountResponse
	status, err := common.GetJsonWithCustomTimeout(url, 0, &response, time.Second*5)
	if err != nil {
		if status == 404 {
			return false, nil
		}
		if retries > 0 {
			return doesWalletExist(restEndpoint, address, retries-1)
		}
		return false, err
	}
	return true, err
}

func DoesWalletExist(restEndpoint string, address string) (bool, error) {
	return doesWalletExist(restEndpoint, address, 1)
}
