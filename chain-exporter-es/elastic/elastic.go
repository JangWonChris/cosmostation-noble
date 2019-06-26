package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmostation-cosmos/chain-exporter-es/model"
	"gopkg.in/olivere/elastic.v5"

	"github.com/aws/aws-sdk-go/aws/credentials"
	aws "github.com/olivere/elastic/aws/v4"
	"github.com/pkg/errors"

	"github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	txIndex   = "tx_index"
	txDocType = "tx"
)

const txMappings = `{
				"mappings":{
			"tx":{
				"properties":{
					"hash":{
						"type":"text"
					},
					"height": {
						"type": "long"
					},
					"time": {
						"type": "date"
					},
					"tx": {
						"properties":{
							"value":{
								"properties":{
									"signatures" : {
										"enabled": false						
									},
									"msg": {
										"properties":{
											"value": {
												"properties":{
													"description": {
														"enabled": false
													}
												}
											}
										}
									}
								}
							}
						}
					}
			}
	}}
}`

type ElasticSearch struct {
	ElasticSearch *elastic.Client
	Cdc *codec.Codec
	TxDecoder sdk.TxDecoder
	Height int64 `json:"height"`
}

func New(config *Config) (*ElasticSearch, error) {
	signingClient := aws.NewV4SigningClient(
		credentials.NewStaticCredentials(
			config.AccessKey,
			config.SecretKey,
			"",
		), config.Region)

	elasticSearch, err := elastic.NewClient(
		elastic.SetURL(config.ElasticHost),
		elastic.SetScheme("https"),
		elastic.SetHealthcheck(config.Sniff),
		elastic.SetSniff(config.Sniff),
		elastic.SetHttpClient(signingClient))

	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to elasticsearch")
	}

	var cdc = amino.NewCodec()

	ctypes.RegisterAmino(cdc)
	sdk.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)


	es := &ElasticSearch{
		ElasticSearch:elasticSearch,
		Cdc:cdc,
		TxDecoder:auth.DefaultTxDecoder(cdc),
		Height:0,
	}
	return es, nil
}

func (es *ElasticSearch) GetNextHeight(ctx context.Context) (int64, error) {
	return es.Height + 1, nil
}

func (es *ElasticSearch) GetCurrHeight(ctx context.Context) (int64, error) {
	searchResult, err := es.ElasticSearch.Search().
		Index(txIndex). // search in index
		Sort("height", false).
		Size(1).
		Do(ctx) // execute

	if err != nil {
		fmt.Printf("Error during execution GetCurrHeightFromES : %s", err.Error())
	}
	return convertSearchResultToEsTxInfo(searchResult), err
}

// convertSearchResultToUsers ...
func convertSearchResultToEsTxInfo(searchResult *elastic.SearchResult) int64 {
	var resHeight int64
	for _, hit := range searchResult.Hits.Hits {
		var txObj *model.ElasticsearchTxInfo
		err := json.Unmarshal(*hit.Source, &txObj)
		if err != nil {
			fmt.Printf("Can't deserialize 'ElasticsearchTxInfo' object : %s", err.Error())
			continue
		}
		resHeight = txObj.Height
	}
	return resHeight
}

func (es *ElasticSearch) CreateIndex(ctx context.Context) {
	err := createIndexIfDoesNotExist(ctx, es.ElasticSearch, txIndex, txMappings)
	if err != nil {
		panic(err)
	}
}

// 최초실행시 index 만들기
func createIndexIfDoesNotExist(ctx context.Context, client *elastic.Client, indexName string, mappings string) error {
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		panic(err)
		return err
	}

	if exists {
		return nil
	}

	res, err := client.CreateIndex(indexName).BodyString(mappings).Do(ctx)
	if err != nil {
		panic(err)
		return err
	}

	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	return nil
}

func (es *ElasticSearch) SetCurrHeight(height int64) error {
	es.Height = height
	return nil
}

func (es *ElasticSearch) InsertTx(ctx context.Context, tx *model.ElasticsearchTxInfo) error {
	_, err := es.ElasticSearch.Index().Index(txIndex).Type(txDocType).BodyJson(tx).Do(ctx)
	if err != nil {
		fmt.Printf("tx insert Error : %s \n", err.Error())
		return err
	}
	return nil
}

func (es *ElasticSearch) UnmarshalBinaryLengthPrefixed(tx []byte, o interface{}) error {
	return es.Cdc.UnmarshalBinaryLengthPrefixed(tx, o)
}

func (es *ElasticSearch) MarsharJson(o interface{}) (json.RawMessage, error) {
	return es.Cdc.MarshalJSON(o)
}

func (es *ElasticSearch) UnMarsharJson(json []byte, o interface{}) error {
	return es.Cdc.UnmarshalJSON(json, o)
}