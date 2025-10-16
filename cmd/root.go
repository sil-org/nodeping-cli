package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

const NodepingTokenKey = "NODEPING_TOKEN"

var nodepingToken string

var rootCmd = &cobra.Command{
	Use:   "nodeping-cli",
	Short: "Script for making calls to nodeping.com",
	Long:  "Script for making calls to nodeping.com",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	// Get Nodeping Token from env vars
	if nodepingToken = os.Getenv(NodepingTokenKey); nodepingToken == "" {
		log.Fatalf("Error: missing required environment variable: %s", NodepingTokenKey)
	}
}
