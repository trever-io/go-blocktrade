package blocktrade

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type APIClient struct {
	apiKey    string
	apiSecret string

	assetCache map[int64]*TradingAsset
	pairCache  map[int64]*TradingPair
}

func NewClient(apiKey, apiSecret string) *APIClient {
	return &APIClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,

		assetCache: make(map[int64]*TradingAsset),
		pairCache:  make(map[int64]*TradingPair),
	}
}

const API_URL = "https://trade.blocktrade.com/api/v1"
const API_KEY_HEADER = "X-Api-Key"
const NONCE_HEADER = "X-Nonce"
const SIGNATURE_HEADER = "X-Signature"
const CONTENT_TYPE_HEADER = "Content-Type"
const CONTENT_TYPE = "application/json"

func (a *APIClient) nonceAndSignature(request interface{}) (int64, string, error) {
	nonce := time.Now().UTC().UnixNano() / 1e3
	message := fmt.Sprintf("%v.%v.", a.apiKey, nonce)
	if request != nil {
		requestb, err := json.Marshal(request)
		if err != nil {
			return 0, "", err
		}

		message += fmt.Sprintf("%v", string(requestb))
	}

	h := hmac.New(sha256.New, []byte(a.apiSecret))
	h.Write([]byte(message))
	sha := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return nonce, sha, nil
}

func (a *APIClient) requestPublicGET(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%v%v", API_URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return a.doRequest(req)
}

func (a *APIClient) requestPOST(endpoint string, request interface{}) ([]byte, error) {
	if a.apiKey == "" || a.apiSecret == "" {
		return nil, errors.New("missing credentials")
	}

	nonce, sig, err := a.nonceAndSignature(request)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)
	url := fmt.Sprintf("%v%v", API_URL, endpoint)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add(API_KEY_HEADER, a.apiKey)
	req.Header.Add(NONCE_HEADER, fmt.Sprint(nonce))
	req.Header.Add(SIGNATURE_HEADER, sig)
	req.Header.Add(CONTENT_TYPE_HEADER, CONTENT_TYPE)

	return a.doRequest(req)
}

func (a *APIClient) requestGET(endpoint string) ([]byte, error) {
	if a.apiKey == "" || a.apiSecret == "" {
		return nil, errors.New("missing credentials")
	}

	nonce, sig, err := a.nonceAndSignature(nil)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v%v", API_URL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add(API_KEY_HEADER, a.apiKey)
	req.Header.Add(NONCE_HEADER, fmt.Sprint(nonce))
	req.Header.Add(SIGNATURE_HEADER, sig)

	return a.doRequest(req)
}

func (a *APIClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
