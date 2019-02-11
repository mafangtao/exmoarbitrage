package service

import (
	"context"
	"github.com/tusupov/exmoarbitrage/model"
)

type Servicer interface {
	GetCurrencyList(context.Context) ([]model.Currency, error)
	GetArbitrage(context.Context) ([]model.Arbitrage, error)
}
