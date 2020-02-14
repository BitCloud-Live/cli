package cmd

import (
	"github.com/spf13/cobra"
)

// sftp Command
var (
	flagCommand = []string{}
	runCmd      = &cobra.Command{
		Use:   "run [app-name]",
		Short: "Run a remote command in a app",
		Long:  `Run a remote command in a app running in a cluster (no tty)`,
		Run:   appRun,
	}
)

func init() {
	rootCmd.AddCommand(
		runCmd)
}
