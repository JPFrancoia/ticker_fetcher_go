// Fetches information for stocks.
// Uses the Yahoo API.

package cmd

import (
	"fmt"
	"local/ticker_fetcher/yahoo"
	"strings"

	"github.com/spf13/cobra"
)

var fiatCmd = &cobra.Command{
	Use:   "fiat",
	Short: "Fetch exchange rates.",
	Long:  "Fetch exchange rates. Use Yahoo symbols. E.g GBP to EUR is GBPEUR=X.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		display := func(data yahoo.YahooInfo) {
			fromTo := strings.Split(data.ShortName, "/")
			fmt.Printf(
				"${alignc}%s: %g %s (%.2f %%)\n",
				fromTo[0],
				data.Price,
				fromTo[1],
				data.Diff(),
			)
		}

		yahoo.ProcessFromYahoo(args[0], display)
	},
}
