package services

import (
	"fmt"
	"testing"

	"github.com/cosmostation/cosmostation-cosmos/wallet/api/models"
)

// Sample accounts
var sampleAccounts = []models.Account{
	{
		ChainID:     1,
		DeviceType:  "android",
		Address:     "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "dIhWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuJtAajGjHuqwNGxtOH0LL1sRi3l4ubgk0KJB7ZIBvoQUal-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
	{
		ChainID:     2,
		DeviceType:  "ios",
		Address:     "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "eIzWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuBtAajMjPuqwAGxtOH0LL1sRi2l1ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
	{
		ChainID:     1,
		DeviceType:  "ANDROID",
		Address:     "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "eIzWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuBtAajMjPuqwAGxtOH0LL1sRi2l1ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
	{
		ChainID:     1,
		DeviceType:  "GOOGLE",
		Address:     "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "eIzWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuBtAajMjPuqwAGxtOH0LL1sRi2l1ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
	{
		ChainID:     1,
		DeviceType:  "GOOGLE",
		Address:     "kava1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "eIzWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuBtAajMjPuqwAGxtOH0LL1sRi2l1ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
	{
		ChainID:     4,
		DeviceType:  "ios",
		Address:     "cosmos1dlpqr4nrk4khgkxj78fcwwh6dl3uw4ep4tgu9c",
		AlarmToken:  "eIzWxooEtBY:APA91bGcX_rhNSi4GfXEdWCa1yle_p7QmZl8CbU5KwFUMkDaKBPi--mBZNwQi3eGUA8KwBJXp9rcd0NuBtAajMjPuqwAGxtOH0LL1sRi2l1ubgk0KJB7ZIBvpLUty-_7C0FriGztaPEn",
		AlarmStatus: true,
	},
}

func TestAccountValidation(t *testing.T) {
	// t.Log("Account[1]", sampleAccounts[0])
	// t.Log("Account[2]", sampleAccounts[1])

	fmt.Println("Account[1]", sampleAccounts[0])
	fmt.Println("Account[2]", sampleAccounts[1])

	t.Error() // to indicate test failed
}
