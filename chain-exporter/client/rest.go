package client

// import (
// 	"encoding/json"
// 	"fmt"

// 	sdktypes "github.com/cosmos/cosmos-sdk/types"
// 	lschema "github.com/cosmostation/mintscan-backend-library/db/schema"
// 	"github.com/cosmostation/mintscan-backend-library/schema"
// 	"github.com/cosmostation/mintscan-backend-library/types"
// )

// // --------------------
// // REST SERVER APIs
// // --------------------

// // GetTxAPIClient queries for a transaction from the REST client and decodes it into a sdktypes.Tx [Another way to query a transaction.]
// // if the transaction exists. An error is returned if the tx doesn't exist or
// // decoding fails.
// func (c *Client) GetTxAPIClient(hash string) (sdktypes.TxResponse, error) {
// 	resp, err := c.apiClient.R().Get("/txs/" + hash)
// 	if err != nil {
// 		return sdktypes.TxResponse{}, fmt.Errorf("failed to request tx hash: %s", err)
// 	}

// 	var txResponse sdktypes.TxResponse
// 	if err := c.cliCtx.LegacyAmino.UnmarshalJSON(resp.Body(), &txResponse); err != nil {
// 		return sdktypes.TxResponse{}, fmt.Errorf("failed to unmarshal tx hash: %s", err)
// 	}

// 	return txResponse, nil
// }

// // GetValidatorsIdentities returns identities of all validators in the active chain.
// func (c *Client) GetValidatorsIdentities(vals []lschema.Validator) (result []lschema.Validator, err error) {
// 	for _, val := range vals {
// 		if val.Identity != "" {
// 			resp, err := c.keyBaseClient.R().Get("_/api/1.0/user/lookup.json?fields=pictures&key_suffix=" + val.Identity)
// 			if err != nil {
// 				return []lschema.Validator{}, fmt.Errorf("failed to request identity: %s", err)
// 			}

// 			var keyBase types.KeyBase
// 			err = json.Unmarshal(resp.Body(), &keyBase)
// 			if err != nil {
// 				return []lschema.Validator{}, fmt.Errorf("failed to unmarshal keybase: %s", err)
// 			}

// 			var url string
// 			if len(keyBase.Them) > 0 {
// 				for _, k := range keyBase.Them {
// 					url = k.Pictures.Primary.URL
// 				}
// 			}

// 			v := lschema.NewValidator(schema.Validator{
// 				ID:         val.ID,
// 				KeybaseURL: url,
// 			})

// 			result = append(result, *v)
// 		}
// 	}

// 	return result, nil
// }
