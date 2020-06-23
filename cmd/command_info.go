package cmd

import (
	"github.com/spf13/cobra"
)

// sftp Command
var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get cloud user account info",
		Long:  `Get cloud user account info about credentials, plans, and subscriptions`,
		Run:   info,
	}
)

func init() {
	rootCmd.AddCommand(
		infoCmd)
}
