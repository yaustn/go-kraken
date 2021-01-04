package gokraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// APIUrl - Kraken API Endpoint
	apiURL = "https://api.kraken.com"

	// APIVersion - Kraken API Version Number
	apiVersion = "0"
)

// Response - template of Kraken API response
type Response struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

// Client for interfacing with the Kraken REST API
type KrakenClient struct {
	httpClient *http.Client
	apiKey     string
	apiSecret  string
}

func NewKrakenClient(key, secret string) *KrakenClient {
	return &KrakenClient{
		httpClient: http.DefaultClient,
		apiKey:     key,
		apiSecret:  secret,
	}
}

func (kc *KrakenClient) request(endpoint string, isPrivate bool, data url.Values, respType interface{}) error {
	req, err := kc.buildRequest(endpoint, isPrivate, data)
	if err != nil {
		return err
	}

	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return errors.New("failed to make http request")
	}
	defer resp.Body.Close()

	return kc.parseResponse(resp, respType)
}

func (kc *KrakenClient) buildRequest(endpoint string, isPrivate bool, data url.Values) (*http.Request, error) {
	if data == nil {
		data = url.Values{}
	}

	var requestURL string
	if isPrivate {
		requestURL = fmt.Sprintf("%s/%s/private/%s", apiURL, apiVersion, endpoint)
		data.Set("nonce", kc.timestamp())
	} else {
		requestURL = fmt.Sprintf("%s/%s/public/%s", apiURL, apiVersion, endpoint)
	}

	req, err := http.NewRequest("POST", requestURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create api request: %s", err.Error())
	}

	if isPrivate {
		urlPath := fmt.Sprintf("/%s/private/%s", apiVersion, endpoint)
		req.Header.Add("API-Key", kc.apiKey)

		signature, err := kc.getSignature(urlPath, data)
		if err != nil {
			return nil, fmt.Errorf("failed to sign private request: %s", err.Error())
		}
		req.Header.Add("API-Sign", signature)
	}

	return req, nil
}

func (kc *KrakenClient) getSignature(requestURL string, data url.Values) (string, error) {
	sha := sha256.New()

	_, err := sha.Write([]byte(data.Get("nonce") + data.Encode()))
	if err != nil {
		return "", err
	}

	hash := sha.Sum(nil)

	secret, err := base64.StdEncoding.DecodeString(kc.apiSecret)
	if err != nil {
		return "", err
	}

	hmacObj := hmac.New(sha512.New, secret)

	_, err = hmacObj.Write(append([]byte(requestURL), hash...))
	if err != nil {
		return "", err
	}

	hmacData := hmacObj.Sum(nil)

	return base64.StdEncoding.EncodeToString(hmacData), nil
}

func (kc *KrakenClient) parseResponse(response *http.Response, respType interface{}) error {
	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to get a successful response. Status %d", response.StatusCode)
	}

	if response.Body == nil {
		return fmt.Errorf("Failed to get a response body")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %+v", err)
	}

	var resp Response
	if respType != nil {
		resp.Result = respType
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response: %+v", err)
	}

	if len(resp.Error) > 0 {
		return fmt.Errorf("Got server errors: %+v", resp.Error)
	}

	return nil
}

// timestamp returns a string formatted epoch millisecond timestamp
func (kc *KrakenClient) timestamp() string {
	epochMilli := int64(time.Now().UTC().UnixNano()) / int64(time.Millisecond)
	return fmt.Sprintf("%d", epochMilli)
}
