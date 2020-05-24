package cmd

import (
	"github.com/spf13/cobra"
)

var (
	openCmd = &cobra.Command{
		Use:   "open [Application.name]",
		Short: "Open an Application in the default browser",
		Long:  `This subcommand will open an application in the default browser in every os.`,
		Run:   appOpen}

	portforwardCmd = &cobra.Command{
		Use:   "portforward [srvice.name]",
		Short: "portforward to connect to an application running in a cluster",
		Long: `Portforward to connect to an application running in a cluster.
		This type of connection can be useful for database debugging`,
		Run: srvPortforward}
)

func init() {
	rootCmd.AddCommand(
		portforwardCmd,
		openCmd)
}
