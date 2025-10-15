package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/sil-org/nodeping-cli/lib"
	"github.com/spf13/cobra"
)

var (
	contactGroupName string
	period           lib.PeriodValue
)

var uptimeCmd = &cobra.Command{
	Use:   "uptime",
	Short: "Get the uptime for checks",
	Long:  "Get the uptime for all the checks for a certain Contact Group.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
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
	uptimeCmd.Flags().VarP(
		&period,
		"period",
		"p",
		fmt.Sprintf(`Name of the time period to get uptime values for ... %v`, periods),
	)
	uptimeCmd.MarkFlagRequired("contact-group")
	uptimeCmd.RegisterFlagCompletionFunc("period", periodAutoCompletion)
}

func periodAutoCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	periods := lib.GetValidPeriods()
	valid := make([]string, len(periods))
	for i, v := range periods {
		valid[i] = string(v)
	}

	return valid, cobra.ShellCompDirectiveDefault
}

func runUptime() {
	results, err := lib.GetUptimesForContactGroup(nodepingToken, contactGroupName, period)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"\nPeriod: %s. From: %v      To: %v\n\n",
		period,
		// time.Unix(start/1000, 0).Format(time.RFC822Z),
		time.Unix(results.StartTime, 0).UTC(),
		time.Unix(results.EndTime, 0).UTC(),
	)

	for _, label := range results.CheckLabels {
		fmt.Printf("%s, %v\n", label, results.Uptimes[label])
	}
}
