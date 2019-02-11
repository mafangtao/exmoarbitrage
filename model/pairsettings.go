package model

type PairSettings map[Pair]Setting

func (p PairSettings) GetSetting(pair Pair) (Setting, bool) {
	v, ok := p[pair]
	return v, ok
}

func (p PairSettings) GetList() (result []Pair) {
	result = make([]Pair, 0, len(p))
	for key, _ := range p {
		result = append(result, key)
	}
	return
}

type Setting struct {
	MinQuantity float64 `json:"min_quantity,string"`
	MaxQuantity float64 `json:"max_quantity,string"`
	MinPrice    float64 `json:"min_price,string"`
	MaxPrice    float64 `json:"max_price,string"`
	MinAmount   float64 `json:"min_amount,string"`
	MaxAmount   float64 `json:"max_amount,string"`
}
