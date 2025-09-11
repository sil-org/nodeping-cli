package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/sil-org/nodeping-cli/lib"
	"time"
)

var contactGroupName string
var period string

var uptimeCmd = &cobra.Command{
	Use: "uptime",
	Short: "Get the uptime for checks",
	Long: "Get the uptime for all the checks for a certain Contact Group.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if contactGroupName == "" {
			fmt.Println("Error: The 'contact-group' flag is required (e.g. -g AppsDev).")
			os.Exit(1)
		}

		if period != "" {
			if ! lib.IsPeriodValid(period) {
				fmt.Printf(
					"Error: The period value is not valid. It must be one of these ...\n%v\n",
					lib.GetValidPeriods())
				os.Exit(1)
			}
		}
		runUptime()
	},
}

func init() {
	periods := lib.GetValidPeriods()

	rootCmd.AddCommand(uptimeCmd)
	uptimeCmd.Flags().StringVarP(
		&contactGroupName,
		"contact-group",
		"g",
		"",
		`Name of the Contact Group to retrieve uptime data for.`,
	)
	uptimeCmd.Flags().StringVarP(
		&period,
		"period",
		"p",
		"",
		fmt.Sprintf(`Name of the time period to get uptime values for ... %v`, periods),
	)
}

func runUptime() {
	results, err := lib.GetUptimesForContactGroup(nodepingToken, contactGroupName, period)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf(
		"\nPeriod: %s. From: %v      To: %v\n\n",
		period,
		//time.Unix(start/1000, 0).Format(time.RFC822Z),
		time.Unix(results.StartTime, 0).UTC(),
		time.Unix(results.EndTime, 0).UTC(),
	)


	for _, label := range results.CheckLabels {
		fmt.Printf("%s, %v\n", label, results.Uptimes[label])
	}
}
