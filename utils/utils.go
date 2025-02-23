package utils

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// Query an API.
// The URL should be pre-formatted.
// This function panics if any error occurs.
func QueryApi(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// User-agent taken from here:
	// https://stackoverflow.com/questions/79450295/yfinance-429-client-error-too-many-requests-for-url
	// This is to avoid 429 errors
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	)

	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	responseBody, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

// Parse a comma-separated string into a slice of strings.
// This is used to parse a list of tickers, e.g:
// BBLN,BNTX,NVDA
func ParseCommaString(argStr string) []string {
	return strings.Split(argStr, ",")
}

// Sometimes the currency of a symbol/ticker is in pennies.
// This function will convert the price to the real unit.
func ConvertDecimal(price float64, currency string) (float64, string) {
	switch currency {
	case "GBp":
		return price / 100, "GBP"
	default:
		return price, currency
	}
}
