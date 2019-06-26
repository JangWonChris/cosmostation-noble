package services

import (
	"fmt"
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/api/wallet/api/models"
)

// Sample accounts
var sampleAccounts = []models.Account{
	{
		IdfAccount: 1,
		AlarmToken: "AAAAAAAAAAAAAAA",
		DeviceType: "android",
		Address:    "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		CoinType:   "ATOM",
		Status:     true,
	},
	{
		IdfAccount: 2,
		AlarmToken: "BBBBBBBBBBBBBBB",
		DeviceType: "android",
		Address:    "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9d",
		CoinType:   "ATOM",
		Status:     true,
	},
}

func TestAccountValidation(t *testing.T) {
	// t.Log("Account[1]", sampleAccounts[0])
	// t.Log("Account[2]", sampleAccounts[1])

	fmt.Println("Account[1]", sampleAccounts[0])
	fmt.Println("Account[2]", sampleAccounts[1])

	t.Error() // to indicate test failed
}
