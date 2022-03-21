package entities

// Currency is a data structure model for currency.
type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CurrencyHistory is a data structure model for currency history.
type CurrencyHistory struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}
