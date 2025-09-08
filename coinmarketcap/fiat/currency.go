package fiat

type Currency struct {
	symbol string
}

func newCurrency(symbol string) Currency {
	return Currency{symbol: symbol}
}

func (c Currency) String() string {
	return c.symbol
}
