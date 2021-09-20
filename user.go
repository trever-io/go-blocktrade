package blocktrade

import "encoding/json"

const USER_ENDPOINT = "/user"

type AccountType string

const AccountType_INDIVIDUAL AccountType = "INDIVIDUAL"
const AccountType_COMPANY AccountType = "COMPANY"

type KycStatus string

const KycStatus_NOT_VERIFIED KycStatus = "NOT_VERIFIED"
const KycStatus_PENDING KycStatus = "PENDING"
const KycStatus_RESUBMIT KycStatus = "RESUBMIT"
const KycStatus_OK KycStatus = "OK"
const KycStatus_REJECTED KycStatus = "REJECTED"

type UserResponse struct {
	UserId             int64                `json:"user_id"`
	Email              string               `json:"email"`
	PrimaryCurrency    *PrimaryTradingAsset `json:"primary_currency"`
	KycStatus          KycStatus            `json:"kyc_status"`
	WebsocketAuthToken string               `json:"websocket_auth_token"`
	Is2FAEnabled       bool                 `json:"is_2fa_enabled"`
	FirstName          string               `json:"first_name"`
	LastName           string               `json:"last_name"`
	Gender             string               `json:"gender"`
	DateOfBirth        int64                `json:"date_of_birth"`
	Address            string               `json:"address"`
	PostalCode         string               `json:"postal_code"`
	City               string               `json:"city"`
	State              string               `json:"state"`
	Country            string               `json:"country"`
	PhoneNumber        string               `json:"phone_number"`
	AmlRequired        bool                 `json:"aml_required"`
	AccountType        AccountType          `json:"account_type"`
}

type PrimaryTradingAsset struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (a *APIClient) User() (*UserResponse, error) {
	b, err := a.requestGET(USER_ENDPOINT)
	if err != nil {
		return nil, err
	}

	resp := new(UserResponse)
	err = json.Unmarshal(b, &resp)
	return resp, err
}
