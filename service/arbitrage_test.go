package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/tusupov/exmoarbitrage/api/mock"
	"github.com/tusupov/exmoarbitrage/model"
)

func TestArbitrageService_GetCurrencyList(t *testing.T) {

	ctx := context.Background()
	currencyList := []model.Currency{"BTC", "RUB", "EUR", "USD"}

	exmoApiMock := mock.NewExmo()
	exmoApiMock.On("GetCurrencyList", ctx).Return(currencyList, nil)

	arbitrageService := NewArbitrage(exmoApiMock)
	resultCurrencyList, err := arbitrageService.GetCurrencyList(ctx)

	assert.Nil(t, err)
	assert.ElementsMatch(t, resultCurrencyList, currencyList)

}

func TestArbitrageService_floydWarshall(t *testing.T) {

	testCase := []struct {
		CurrencyList []model.Currency
		PairOrders   model.PairOrders
		Result       []model.Arbitrage
	}{
		{
			CurrencyList: []model.Currency{"BTC", "USD"},
			PairOrders: model.PairOrders{
				"BTC_USD": model.Order{
					Bid: model.Offer{Price: 3700},
					Ask: model.Offer{Price: 3700},
				},
			},
			Result: []model.Arbitrage{
				{
					Profit: 1,
					Route:  []model.Currency{"BTC", "USD", "BTC"},
				}, {
					Profit: 1,
					Route:  []model.Currency{"USD", "BTC", "USD"},
				},
			},
		},
		{
			CurrencyList: []model.Currency{"BTC", "USD"},
			PairOrders: model.PairOrders{
				"BTC_USD": model.Order{
					Bid: model.Offer{Price: 3600},
					Ask: model.Offer{Price: 3700},
				},
			},
			Result: []model.Arbitrage{
				{
					Profit: 3600 / 3700.0,
					Route:  []model.Currency{"BTC", "USD", "BTC"},
				}, {
					Profit: 3600 / 3700.0,
					Route:  []model.Currency{"USD", "BTC", "USD"},
				},
			},
		},
	}

	for _, test := range testCase {

		arbitrageService := &ArbitrageService{}
		result := arbitrageService.floydWarshall(test.CurrencyList, test.PairOrders)

		if assert.Equal(t, len(test.Result), len(result)) {

			for i := 0; i < len(result); i++ {
				assert.Equal(t, test.Result[i].Profit, result[i].Profit)
				assert.ElementsMatch(t, test.Result[i].Route, result[i].Route)
			}

		}

	}

}

func TestArbitrageService_GetArbitrage(t *testing.T) {

	testCase := []struct {
		CurrencyList []model.Currency
		PairOrders   model.PairOrders
		Result       []model.Arbitrage
	}{
		{
			CurrencyList: []model.Currency{"BTC", "USD"},
			PairOrders: model.PairOrders{
				"BTC_USD": model.Order{
					Bid: model.Offer{Price: 3700},
					Ask: model.Offer{Price: 3700},
				},
			},
			Result: []model.Arbitrage{
				{
					Profit: 1,
					Route:  []model.Currency{"BTC", "USD", "BTC"},
				}, {
					Profit: 1,
					Route:  []model.Currency{"USD", "BTC", "USD"},
				},
			},
		},
		{
			CurrencyList: []model.Currency{"BTC", "USD"},
			PairOrders: model.PairOrders{
				"BTC_USD": model.Order{
					Bid: model.Offer{Price: 3600},
					Ask: model.Offer{Price: 3700},
				},
			},
			Result: []model.Arbitrage{
				{
					Profit: 3600 / 3700.0,
					Route:  []model.Currency{"BTC", "USD", "BTC"},
				}, {
					Profit: 3600 / 3700.0,
					Route:  []model.Currency{"USD", "BTC", "USD"},
				},
			},
		},
	}

	for _, test := range testCase {

		ctx := context.Background()

		exmoApiMock := mock.NewExmo()
		exmoApiMock.On("GetCurrencyList", ctx).Return(test.CurrencyList, nil)
		exmoApiMock.On("GetPairList", ctx).Return(model.PairSettings{}, nil)
		exmoApiMock.On("GetOrders", ctx, []model.Pair{}).Return(test.PairOrders, nil)

		arbitrageService := NewArbitrage(exmoApiMock)

		result, err := arbitrageService.GetArbitrage(ctx)

		if assert.Nil(t, err) && assert.Equal(t, len(test.Result), len(result)) {
			for i := 0; i < len(result); i++ {
				assert.Equal(t, test.Result[i].Profit, result[i].Profit)
				assert.ElementsMatch(t, test.Result[i].Route, result[i].Route)
			}
		}

	}

}
