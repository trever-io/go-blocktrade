package blocktrade

import "encoding/json"

const TRADING_ASSETS_ENDPOINT = "/trading_assets"

type CurrencyType string

const CurrencyType_FIAT CurrencyType = "FIAT"
const CurrencyType_CRYPTOCURRENCY CurrencyType = "CRYPTOCURRENCY"

type DepositMethod string

const DepositMethod_WALLET_ADDRESS = "WALLET_ADDRESS"
const DepositMethod_CLEAR_JUNCTION_SEPA = "CLEAR_JUNCTION_SEPA"
const DepositMethod_COINIFY = "COINIFY"

type TradingAsset struct {
	Id                      int64           `json:"id"`
	FullName                string          `json:"full_name"`
	IsoCode                 string          `json:"iso_code"`
	IconPath                string          `json:"icon_path"`
	IconPathPng             string          `json:"icon_path_png"`
	Color                   string          `json:"color"`
	Sign                    string          `json:"sign"`
	CurrencyType            CurrencyType    `json:"currency_type"`
	MinimalWithdrawalAmount string          `json:"minimal_withdrawal_amount"`
	MinimalOrderValue       string          `json:"minimal_order_value"`
	MaximumOrderValue       string          `json:"maximum_order_value"`
	DecimalPrecision        int64           `json:"decimal_precision"`
	LotSize                 string          `json:"lot_size"`
	DepositMethods          []DepositMethod `json:"deposit_methods"`
}

func (a *APIClient) TradingAssets() ([]*TradingAsset, error) {
	b, err := a.requestPublicGET(TRADING_ASSETS_ENDPOINT)
	if err != nil {
		return nil, err
	}

	resp := make([]*TradingAsset, 0)
	err = json.Unmarshal(b, &resp)
	return resp, err
}