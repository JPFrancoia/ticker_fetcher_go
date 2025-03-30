package yahoo

import (
	"cmp"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"

	"local/ticker_fetcher/utils"

	"github.com/tidwall/gjson"
)

const URL_API = "https://query2.finance.yahoo.com/v8/finance/chart"

// Fetch the market information for one symbol/ticker.
// Returns the info + pushes the info to a channel.
func FetchInfoFromYahoo(symbol string, c chan<- YahooInfo, wg *sync.WaitGroup) YahooInfo {
	defer wg.Done()

	responseBody := utils.QueryApi(fmt.Sprintf("%s/%s?range=%s&interval=%s", URL_API, symbol, "1d", "1d"))

	value := gjson.GetBytes(responseBody, "chart.result.0.meta")

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
	PreviousClose    float64 `json:"chartPreviousClose"`
	Currency         string  `json:"currency"`
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
	tickersAndAliases := utils.ParseCommaString(symbols)

	// Build a map of aliases (if any) and a list of tickers
	tickerToAlias := make(map[string]string)
	var tickers []string

	for _, tickAlias := range tickersAndAliases {
		// We use the ":" separator to define an alias for a ticker
		tickAliasSplit := strings.Split(tickAlias, ":")

		if len(tickAliasSplit) == 2 {
			// Map the ticker to its alias, if alias was provided
			tickerToAlias[tickAliasSplit[0]] = tickAliasSplit[1]
		}
		tickers = append(tickers, tickAliasSplit[0])
	}

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
		// If an alias was provided, use it instead of the ticker
		if alias, ok := tickerToAlias[data.Symbol]; ok {
			data.Symbol = alias
		}
		results = append(results, data)
	}

	// Ascending sort by symbol name
	slices.SortFunc(results, func(a, b YahooInfo) int {
		return cmp.Compare(a.Symbol, b.Symbol)
	})

	for _, data := range results {
		display(data)
	}
}
