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

func (a *APIClient) GetPortfolioId() (int64, error) {
	if a.portfolio != nil {
		return a.portfolio.Id, nil
	}

	portfolio, err := a.Portfolios()
	if err != nil {
		return 0, err
	}

	a.portfolio = portfolio[0]

	return a.portfolio.Id, nil
}
