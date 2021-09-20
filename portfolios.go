package blocktrade

import "encoding/json"

const PORTFOLIOS_ENDPOINT = "/portfolios"

type Portfolio struct {
	Id     int64             `json:"id"`
	Assets []*PortfolioAsset `json:"assets"`
}

type PortfolioAsset struct {
	TradingAssetId  int64  `json:"trading_asset_id"`
	AvailableAmount string `json:"available_amount"`
	ReservedAmount  string `json:"reserved_amount"`
	WalletAddress   string `json:"wallet_address"`
}

func (a *APIClient) Portfolios() ([]*Portfolio, error) {
	b, err := a.requestGET(PORTFOLIOS_ENDPOINT)
	if err != nil {
		return nil, err
	}

	resp := make([]*Portfolio, 0)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
