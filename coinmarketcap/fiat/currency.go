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

func init() {
	for _, curr := range All() {
		currencyBySymbol[strings.ToLower(curr.symbol)] = curr
	}
}

var currencyBySymbol = map[string]Currency{}

func BySymbol(s string) (Currency, bool) {
	curr, ok := currencyBySymbol[strings.ToLower(s)]
	return curr, ok
}
