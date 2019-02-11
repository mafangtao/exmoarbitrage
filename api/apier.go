package api

import (
	"context"
	"github.com/tusupov/exmoarbitrage/model"
)

type Apier interface {
	GetCurrencyList(ctx context.Context) ([]model.Currency, error)
	GetPairList(ctx context.Context) (model.PairSettings, error)
	GetOrders(ctx context.Context, pairs ...model.Pair) (model.PairOrders, error)
}
