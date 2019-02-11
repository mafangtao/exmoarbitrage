package model

import (
	"strings"
)

type Pair string

func (p Pair) Split() (Currency, Currency, bool) {

	tmp := strings.Split(string(p), "_")
	if len(tmp) != 2 {
		return "", "", false
	}

	return Currency(tmp[0]), Currency(tmp[1]), true
}

func (p Pair) Reverse() Pair {

	tmp := strings.Split(string(p), "_")
	for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}

	return Pair(strings.Join(tmp, "_"))
}
