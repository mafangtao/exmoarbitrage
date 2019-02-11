package service

import (
	"context"
	"math"
	"sort"
	"github.com/tusupov/exmoarbitrage/api"
	"github.com/tusupov/exmoarbitrage/model"
)

type ArbitrageService struct {
	api api.Apier
}

func NewArbitrage(api api.Apier) *ArbitrageService {
	return &ArbitrageService{
		api: api,
	}
}

func (s *ArbitrageService) GetCurrencyList(ctx context.Context) (result []model.Currency, err error) {
	return s.api.GetCurrencyList(ctx)
}

// Get Arbitrage list from orders
func (s *ArbitrageService) GetArbitrage(ctx context.Context) (result []model.Arbitrage, err error) {

	currencyList, err := s.api.GetCurrencyList(ctx)
	if err != nil {
		return
	}

	pairList, err := s.api.GetPairList(ctx)
	if err != nil {
		return
	}

	pairOrders, err := s.api.GetOrders(ctx, pairList.GetList()...)
	if err != nil {
		return
	}

	result = s.floydWarshall(currencyList, pairOrders)

	return

}

// Floyd-Worshel Algorithm
// Buy and Sell
func (s ArbitrageService) floydWarshall(currencyList []model.Currency, pairOrders model.PairOrders) (result []model.Arbitrage) {

	const EPS = 1e-6

	n := len(currencyList)
	dist, p := make([][]float64, n), make([][]int, n)

	// Init arrays
	for i := 0; i < n; i++ {

		dist[i], p[i] = make([]float64, n), make([]int, n)

		for j := 0; j < n; j++ {

			pair := model.Pair(currencyList[i] + "_" + currencyList[j])
			if order, ok := pairOrders.GetOrder(pair); ok {
				dist[i][j] = order.Bid.Price
			} else {
				if order, ok := pairOrders.GetOrder(pair.Reverse()); ok {
					dist[i][j] = 1 / order.Ask.Price
				}
			}

			p[i][j] = i

		}
	}

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {

				// Check pair exchange
				mustPair := model.Pair(currencyList[i] + "_" + currencyList[j])
				if !pairOrders.Exists(mustPair) && !pairOrders.Exists(mustPair.Reverse()) {
					continue
				}

				// Find the maximum profitable course
				price := dist[i][k] * dist[k][j]
				if dist[i][j]+EPS < price {
					dist[i][j] = price
					p[i][j] = p[k][j]
				}

			}
		}
	}

	// Full path recovery
	returnRoute := func(i, j int) (route []model.Currency) {

		route = append(route, currencyList[i])

		k := i
		for k != p[k][j] {
			k = p[k][j]
			route = append(route, currencyList[k])
		}
		route = append(route, currencyList[j])
		route = append(route, currencyList[i])

		return route
	}

	// Result
	for i := 0; i < n; i++ {

		for j := 0; j < n; j++ {

			var offer model.Offer

			pair := model.Pair(currencyList[i] + "_" + currencyList[j])
			if order, ok := pairOrders.GetOrder(pair); ok {
				offer = order.Ask
			} else if order, ok := pairOrders.GetOrder(pair.Reverse()); ok {
				offer = order.Bid
				offer.Price = 1 / offer.Price
			} else {
				// If there is no exchange between courses, skip
				continue
			}

			// Add to result
			result = append(result, model.Arbitrage{
				Profit: dist[i][j] / offer.Price,
				Route:  returnRoute(i, j),
			})

		}
	}

	// Sort descending
	// if the value is equal, the sorting will be according to the number of paths to increase
	sort.Slice(result, func(i, j int) bool {
		if math.Abs(result[i].Profit-result[j].Profit) <= EPS {
			return len(result[i].Route) < len(result[j].Route)
		}
		return result[i].Profit > result[j].Profit
	})

	return

}
