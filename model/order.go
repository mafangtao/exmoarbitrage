package model

type PairOrders map[Pair]Order

func (po PairOrders) GetOrder(pair Pair) (Order, bool) {
	order, ok := po[pair]
	return order, ok
}

func (po PairOrders) Exists(pair Pair) bool {

	if _, ok := po.GetOrder(pair); ok {
		return true
	}
	return false

}

type Order struct {
	Ask Offer
	Bid Offer
}

type Offer struct {
	Price float64
	//Quantity float64 `json:",string"`
	//Amount   float64 `json:",string"`
}
