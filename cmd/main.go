package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sil-org/nodeping-cli"
	"github.com/spf13/cobra"
)

func main() {
	const NodepingTokenKey = "NODEPING_TOKEN"

	var contactGroupName string
	nodepingToken := getRequiredEnv(NodepingTokenKey)
	period := nodeping.GetTodayPeriod(time.Now().UTC())

	rootCmd := &cobra.Command{
		Use:   "nodeping-cli",
		Short: "Script for making calls to nodeping.com",
		Long:  "Script for making calls to nodeping.com",
	}
	uptimeCmd := &cobra.Command{
		Use:   "uptime",
		Short: "Get the uptime for checks",
		Long:  "Get the uptime for all the checks for a certain Contact Group.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runUptime(nodepingToken, contactGroupName, period)
		},
	}
	rootCmd.AddCommand(uptimeCmd)

	uptimeCmd.Flags().StringVarP(
		&contactGroupName,
		"contact-group",
		"g",
		"",
		`Name of the Contact Group to retrieve uptime data for.`,
	)
	_ = uptimeCmd.MarkFlagRequired("contact-group")

	uptimeCmd.Flags().VarP(
		&period,
		"period",
		"p",
		fmt.Sprintf(`Name of the time period to get uptime values for ... %v`, nodeping.GetValidPeriods()),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func runUptime(nodepingToken, contactGroupName string, period nodeping.Period) {
	results, err := nodeping.GetUptimesForContactGroup(nodepingToken, contactGroupName, period)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s\n\n", period)
	for _, label := range results.CheckLabels {
		fmt.Printf("%s, %v\n", label, results.Uptimes[label])
	}
}

func getRequiredEnv(name string) (value string) {
	if value = os.Getenv(name); value == "" {
		log.Fatalf("Error: missing required environment variable: %s", name)
	}
	return value
}
