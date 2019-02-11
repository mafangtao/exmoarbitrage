package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"github.com/tusupov/exmoarbitrage/model"
)

func newTestClient() (*httptest.Server, *http.Client, string) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.URL.Path {

		case "/currency/":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`["USD","EUR","RUB"]`))
			return

		case "/pair_settings/":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"BTC_USD":{"min_quantity":"0.001","max_quantity":"1000","min_price":"1","max_price":"30000","max_amount":"500000","min_amount":"1"},"BTC_EUR":{"min_quantity":"0.001","max_quantity":"1000","min_price":"1","max_price":"30000","max_amount":"500000","min_amount":"1"},"BTC_RUB":{"min_quantity":"0.001","max_quantity":"1000","min_price":"1","max_price":"2000000","max_amount":"50000000","min_amount":"10"}}`))
			return

		case "/order_book/":
			query := r.URL.Query()
			pairList := strings.Split(query.Get("pair"), ",")
			sort.Slice(pairList, func(i, j int) bool {
				return pairList[i] > pairList[j]
			})

			if len(pairList) == 3 && assert.ObjectsAreEqual(pairList, []string{"BTC_USD", "BTC_RUB", "BTC_EUR"}) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"BTC_USD":{"ask_quantity":"322.85129185","ask_amount":"2506250.14134118","ask_top":"3681.40738965","bid_quantity":"11732.67600564","bid_amount":"579061.75794034","bid_top":"3670.00211426","ask":[["3681.40738965","0.0230784","84.9609923"],["3764.996235","0.0013876","5.22430877"]],"bid":[["3670.00211426","0.03715396","136.35511175"],["3618.201094","0.01","36.18201094"]]},"BTC_EUR":{"ask_quantity":"97.51416943","ask_amount":"511716.00372549","ask_top":"3286.27812155","bid_quantity":"123.61959283","bid_amount":"281835.26698855","bid_top":"3267.85080247","ask":[["3286.27812155","0.01787873","58.75447924"],["3749.59155","0.0011","4.1245507"]],"bid":[["3267.85080247","1.1883","3883.18710857"],["3224.91052006","0.01","32.2491052"]]},"BTC_RUB":{"ask_quantity":"157.70795936","ask_amount":"66710379.06824972","ask_top":"246638.52016029","bid_quantity":"5113086.66286215","bid_amount":"24941065.74845995","bid_top":"245953","ask":[["246638.52016029","0.27315","67369.31178178"],["255130","0.34685","88491.8405"]],"bid":[["245953","0.0034964","859.9500692"],["237700","0.50783","120711.191"]]}}`))
			}
			return

		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not found`))

	}))

	return server, server.Client(), server.URL

}

func TestExmo_GetCurrencyList(t *testing.T) {

	server, client, baseUrl := newTestClient()
	defer server.Close()

	api := NewExmo(baseUrl, client)

	ctx := context.Background()
	list, err := api.GetCurrencyList(ctx)

	assert.Nil(t, err)
	assert.ElementsMatch(t, list, []model.Currency{"USD", "EUR", "RUB"})

}

func TestExmo_GetPairList(t *testing.T) {

	server, client, baseUrl := newTestClient()
	defer server.Close()

	api := NewExmo(baseUrl, client)
	ctx := context.Background()

	pairList, err := api.GetPairList(ctx)
	assert.Nil(t, err)
	assert.ElementsMatch(t, pairList.GetList(), []model.Pair{"BTC_EUR", "BTC_USD", "BTC_RUB"})

	pairSetting, ok := pairList.GetSetting("BTC_USD")
	assert.True(t, ok)
	assert.Equal(t, pairSetting, model.Setting{MinQuantity: 1e-3, MaxQuantity: 1e3, MinPrice: 1, MaxPrice: 3e4, MinAmount: 1, MaxAmount: 5e5})

	_, ok = pairList.GetSetting("BTC_ETH")
	assert.False(t, ok)

}

func TestExmo_GetOrders(t *testing.T) {

	server, client, baseUrl := newTestClient()
	defer server.Close()

	api := NewExmo(baseUrl, client)
	ctx := context.Background()

	pairList := []model.Pair{"BTC_EUR", "BTC_USD", "BTC_RUB"}

	pairOrders, err := api.GetOrders(ctx, pairList...)
	assert.Nil(t, err)
	for _, pair := range pairList {
		assert.True(t, pairOrders.Exists(pair))
	}

	pairList = append(pairList, model.Pair("FAKE_PAIR"))
	_, err = api.GetOrders(ctx, pairList...)
	assert.NotNil(t, err)

}
