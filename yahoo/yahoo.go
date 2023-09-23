package yahoo

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"sync"

	"local/ticker_fetcher/utils"

	"github.com/tidwall/gjson"
)

const URL_API = "https://query2.finance.yahoo.com/v7/finance/options"

// Fetch the market information for one symbol/ticker.
// Returns the info + pushes the info to a channel.
func FetchInfoFromYahoo(symbol string, c chan<- YahooInfo, wg *sync.WaitGroup) YahooInfo {
	defer wg.Done()

	responseBody := utils.QueryApi(fmt.Sprintf("%s/%s", URL_API, symbol))

	value := gjson.GetBytes(responseBody, "optionChain.result.0.quote")

	var target YahooInfo
	if err := json.Unmarshal([]byte(value.Raw), &target); err != nil {
		log.Fatal(err)
	}

	c <- target

	return target
}

type YahooInfo struct {
	Symbol           string  `json:"symbol"`
	FullExchangeName string  `json:"fullExchangeName"`
	ExchangeName     string  `json:"exchange"`
	Price            float64 `json:"regularMarketPrice"`
	PreviousClose    float64 `json:"regularMarketPreviousClose"`
	Currency         string  `json:"currency"`
	FromCurrency     string  `json:"fromCurrency"`
	ShortName        string  `json:"shortName"`
}

// Compute the diff since the previous close (most of the time, this is a one-day diff)
func (fi *YahooInfo) Diff() float64 {
	return (fi.Price - fi.PreviousClose) / fi.PreviousClose * 100
}

// Main function, triggers the parsing of the input data + the info fetching.
// The info for the symbols will be fetched at the same time, asynchronously.
func ProcessFromYahoo(symbols string, display func(YahooInfo)) {

	// Extract all the symbols from a comma-separated list
	tickers := utils.ParseCommaString(symbols)

	// Make a wait group to wait until the info for all the symbols is fetched
	c := make(chan YahooInfo)
	var wg sync.WaitGroup

	for _, tick := range tickers {
		wg.Add(1)
		go FetchInfoFromYahoo(tick, c, &wg)
	}

	// Necessary to avoid deadlock
	// https://stackoverflow.com/a/70877210/1585507
	go func() {
		wg.Wait()
		close(c)
	}()

	// Read from the channel and build a list of results
	var results []YahooInfo
	for data := range c {
		results = append(results, data)
	}

	// Ascending sort by symbol name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Symbol < results[j].Symbol
	})

	for _, data := range results {
		display(data)
	}
}
