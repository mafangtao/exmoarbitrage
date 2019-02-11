package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tusupov/exmoarbitrage/model"
)

type exmo struct {
	mock.Mock
}

func NewExmo() *exmo {
	return &exmo{}
}

func (m *exmo) GetCurrencyList(ctx context.Context) ([]model.Currency, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Currency), args.Error(1)
}

func (m *exmo) GetPairList(ctx context.Context) (model.PairSettings, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.PairSettings), args.Error(1)
}

func (m *exmo) GetOrders(ctx context.Context, pairs ...model.Pair) (model.PairOrders, error) {
	args := m.Called(ctx, pairs)
	return args.Get(0).(model.PairOrders), args.Error(1)
}
