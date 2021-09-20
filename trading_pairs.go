package blocktrade

import (
	"encoding/json"
	"fmt"
)

const TRADING_PAIRS_ENDPOINT = "/trading_pairs"

type TradingPair struct {
	Id               int64  `json:"id"`
	BaseAssetId      int64  `json:"base_asset_id"`
	QuoteAssetId     int64  `json:"quote_asset_id"`
	DecimalPrecision int    `json:"decimal_precision"`
	LotSize          string `json:"lot_size"`
	TickSize         string `json:"tick_size"`
}

func (a *APIClient) TradingPairs() ([]*TradingPair, error) {
	b, err := a.requestPublicGET(TRADING_PAIRS_ENDPOINT)
	if err != nil {
		return nil, err
	}

	resp := make([]*TradingPair, 0)
	err = json.Unmarshal(b, &resp)
	return resp, err
}

func (a *APIClient) TradingPairFromId(id int64) (*TradingPair, error) {
	if val, ok := a.pairCache[id]; ok {
		return val, nil
	}

	// not in cache. refetching
	pairs, err := a.TradingPairs()
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		a.pairCache[pair.Id] = pair
	}

	if val, ok := a.pairCache[id]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("pair not found for id %d", id)
	}
}

func (a *APIClient) TradingPairFromBaseQuote(baseId int64, quoteId int64) (*TradingPair, error) {
	for _, pair := range a.pairCache {
		if pair.BaseAssetId == baseId && pair.QuoteAssetId == quoteId {
			return pair, nil
		}
	}

	// not in cache. refetching
	pairs, err := a.TradingPairs()
	if err != nil {
		return nil, err
	}

	for _, pair := range pairs {
		a.pairCache[pair.Id] = pair
	}

	for _, pair := range a.pairCache {
		if pair.BaseAssetId == baseId && pair.QuoteAssetId == quoteId {
			return pair, nil
		}
	}

	return nil, fmt.Errorf("pair not founf for base %d and quote %d", baseId, quoteId)
}
