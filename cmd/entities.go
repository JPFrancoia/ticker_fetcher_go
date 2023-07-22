// This file contains the various entities used by the ticker_fetcher CLI.
package cmd

type fundInfo struct {
	Symbol        string  `json:"symbol"`
	ExchangeName  string  `json:"fullExchangeName"`
	Price         float64 `json:"regularMarketPrice"`
	PreviousClose float64 `json:"regularMarketPreviousClose"`
	Currency      string  `json:"currency"`
}

func (fi *fundInfo) diff() float64 {
	return (fi.Price - fi.PreviousClose) / fi.PreviousClose * 100
}
