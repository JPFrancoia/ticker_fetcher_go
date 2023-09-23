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
	r, err := http.Get(url)

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
