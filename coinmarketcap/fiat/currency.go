package fiat

import "strings"

type Currency struct {
	symbol string
}

func newCurrency(symbol string) Currency {
	return Currency{symbol: symbol}
}

func (c Currency) String() string {
	return c.symbol
}

func BySymbol(s string) (Currency, bool) {
	for _, curr := range All() {
		if strings.EqualFold(s, curr.symbol) {
			return curr, true
		}
	}
	return Currency{}, false
}
