package coinmarketcap

var (
	CurrencyBTC = newCurrency("BTC", 1)
	CurrencyLTC = newCurrency("LTC", 2)
	CurrencyKAS = newCurrency("KAS", 20396)
	CurrencyXMR = newCurrency("KAS", 328)
)

type Currency struct {
	symbol string
	id     int
}

func newCurrency(symbol string, id int) Currency {
	return Currency{symbol: symbol, id: id}
}

func (c Currency) String() string {
	return c.symbol
}

type CurrencyUpdate struct {
	Data      CurrencyUpdateData `json:"d"`
	Timestamp string             `json:"t"` // timestamp, appears to be a string representation of an int
	Context   string             `json:"c"` // context or source tag
}

type CurrencyUpdateData struct {
	ID                      int     `json:"id"`
	Price                   float64 `json:"p"`        // price
	Volume24H               float64 `json:"v"`        // 24h volume
	Price1H                 float64 `json:"p1h"`      // 1 hour %
	Price24H                float64 `json:"p24h"`     // 24 hour %
	Price7D                 float64 `json:"p7d"`      // 7 day %
	Price30D                float64 `json:"p30d"`     // 30 day %
	Price3M                 float64 `json:"p3m"`      // 3 month %
	Price1Y                 float64 `json:"p1y"`      // 1 year %
	PriceYearToDate         float64 `json:"pytd"`     // Year to date %
	PriceAllTime            float64 `json:"pall"`     // All time %
	TotalSupply             float64 `json:"ts"`       // total supply
	AvailableSupply         float64 `json:"as"`       // available supply
	FullyDilutedMarketCap   float64 `json:"fmc"`      // fully diluted market cap
	MarketCap               float64 `json:"mc"`       // market cap
	MarketCap24HPercent     float64 `json:"mc24hpc"`  // market cap 24h %
	Volume24HPercent        float64 `json:"vol24hpc"` // volume 24h %
	FullMarketCap24HPercent float64 `json:"fmc24hpc"` // full market cap 24h %
	Dominance               float64 `json:"d"`        // dominance
	VolumeDominance         float64 `json:"vd"`       // volume dominance
}
