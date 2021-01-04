package gokraken

import (
	"log"
	"net/url"
	"strconv"
)

// OrderDescription - structure of order description
type OrderDescription struct {
	Pair           string  `json:"pair"`
	Side           string  `json:"type"` // buy/sell
	OrderType      string  `json:"ordertype"`
	Price          float64 `json:"price,string"`
	Price2         float64 `json:"price2,string"`
	Leverage       string  `json:"leverage"`
	Info           string  `json:"order"`
	CloseCondition string  `json:"close"`
}

type AddOrderResponse struct {
	Description    OrderDescription `json:"descr"`
	TransactionIds []string         `json:"txid"`
}

type CancelResponse struct {
	Count   int64 `json:"count"`
	Pending bool  `json:"pending,omitempty"`
}

func (kc *KrakenClient) AddOrder(pair string, side string, orderType string, volume float64, args map[string]interface{}) (*AddOrderResponse, error) {
	data := url.Values{
		"pair":      {pair},
		"type":      {side},
		"ordertype": {orderType},
		"volume":    {strconv.FormatFloat(volume, 'f', 8, 64)},
	}

	for key, value := range args {
		switch v := value.(type) {
		case string:
			data.Set(key, v)
		case int64:
			data.Set(key, strconv.FormatInt(v, 10))
		case float64:
			data.Set(key, strconv.FormatFloat(v, 'f', 8, 64))
		case bool:
			data.Set(key, strconv.FormatBool(v))
		default:
			log.Printf("[WARNING] Unknown value type %v for key %s", value, key)
		}
	}

	response := AddOrderResponse{}
	err := kc.request("AddOrder", true, data, &response)
	if err != nil {
		return &response, err
	}
	return &response, nil
}

func (kc *KrakenClient) CancelOrder(orderID string) (*CancelResponse, error) {
	data := url.Values{
		"txid": {orderID},
	}

	response := CancelResponse{}
	err := kc.request("CancelOrder", true, data, &response)
	if err != nil {
		return &response, err
	}

	return &response, nil
}
