package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/tusupov/exmoarbitrage/model"
)

const (
	ExmoBaseUrl = "https://api.exmo.com/v1" // api url
	CacheTime   = 24 * 60 * 60              // 24 hour
)

var (
	ErrPairEmpty        = errors.New("Валютные пары пустые")
	ErrPairMustNotEmpty = errors.New("`pairs` не должен быть пустым")
)

type exmo struct {
	baseUrl string
	client  *http.Client

	pairsList model.PairSettings
	pairsTime time.Time

	currencyList []model.Currency
	currencyTime time.Time
}

func NewExmo(baseUrl string, client *http.Client) *exmo {
	if client == nil {
		client = http.DefaultClient
	}

	return &exmo{
		baseUrl: baseUrl,
		client:  client,
	}
}

// Get currency list
func (e *exmo) GetCurrencyList(ctx context.Context) (list []model.Currency, err error) {

	// Load from cache
	if time.Since(e.currencyTime).Seconds() < CacheTime {
		return e.currencyList, nil
	}

	resp, err := e.doRequest(ctx, http.MethodGet, e.baseUrl+"/currency/", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return
	}

	// Caching
	e.currencyList = list
	e.currencyTime = time.Now()

	return
}

// Get Pair list with settings
func (e *exmo) GetPairList(ctx context.Context) (pairs model.PairSettings, err error) {

	// Load from cache
	if time.Since(e.pairsTime).Seconds() < CacheTime {
		return e.pairsList, nil
	}

	resp, err := e.doRequest(ctx, http.MethodGet, e.baseUrl+"/pair_settings/", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&pairs)
	if err != nil {
		return
	}

	if len(pairs) == 0 {
		err = ErrPairEmpty
		return
	}

	// Caching
	e.pairsList = pairs
	e.pairsTime = time.Now()

	return
}

// Get Orders for pairs
func (e *exmo) GetOrders(ctx context.Context, pairs ...model.Pair) (pairOrders model.PairOrders, err error) {

	if len(pairs) == 0 {
		err = errors.New("`pairs` must not be empty")
		return
	}

	pairParam := ""
	for _, pair := range pairs {
		pairParam += string(pair) + ","
	}
	pairParam = strings.TrimRight(pairParam, ",")

	resp, err := e.doRequest(ctx, http.MethodGet, e.baseUrl+"/order_book/?limit=1&pair="+pairParam, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var bodyStruct map[model.Pair]struct {
		Ask [][]string
		Bid [][]string
	}
	err = json.NewDecoder(resp.Body).Decode(&bodyStruct)
	if err != nil {
		return
	}

	pairOrders = model.PairOrders{}

	for bodyPair, bodyOrder := range bodyStruct {

		if len(bodyOrder.Ask) == 0 || len(bodyOrder.Bid) == 0 {
			continue
		}

		if len(bodyOrder.Ask[0]) != 3 || len(bodyOrder.Bid[0]) != 3 {
			continue
		}

		askPrice, errParse := strconv.ParseFloat(bodyOrder.Ask[0][0], 64)
		if errParse != nil {
			continue
		}

		bidPrice, errParse := strconv.ParseFloat(bodyOrder.Bid[0][0], 64)
		if errParse != nil {
			continue
		}

		pairOrders[bodyPair] = model.Order{
			Ask: model.Offer{
				Price: askPrice,
			},
			Bid: model.Offer{
				Price: bidPrice,
			},
		}

	}

	return

}

func (e *exmo) doRequest(ctx context.Context, method string, url string, body io.Reader) (resp *http.Response, err error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		e.pairsTime = time.Time{}
		return
	}

	req = req.WithContext(ctx)

	return e.client.Do(req)

}
