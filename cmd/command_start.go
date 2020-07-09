package cmd

import (
	"github.com/spf13/cobra"
)

// START command
var (
	startCmd = &cobra.Command{
		Use:   "start [command] [name]",
		Short: "start the stopped container",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	startAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "start the running application",
		Long:  `This subcommand start the running application.`,
		Run:   appStart}

	startSrvCmd = &cobra.Command{
		Use:   "service [name]",
		Short: "start the stopped service",
		Long:  `This subcommand start the stopped service.`,
		Run:   srvStart}
)

// STOP  command
var (
	stopCmd = &cobra.Command{
		Use:   "stop [command] [name]",
		Short: "stop the running container",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	stopAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "stop the running application",
		Long:  `This subcommand stop the running application.`,
		Run:   appStop}

	stopSrvCmd = &cobra.Command{
		Use:   "service [name]",
		Short: "stop the running service",
		Long:  `This subcommand stop the running service.`,
		Run:   srvStop}
)

// RESET command
var (
	resetCmd = &cobra.Command{
		Use:   "reset [command] [name]",
		Short: "reset the running container",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	resetAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "reset the running application",
		Long:  `Destroys an application instance and create a new instance of application.`,
		Run:   appReset}
)

func init() {
	rootCmd.AddCommand(
		startCmd,
		stopCmd,
		resetCmd)

	startCmd.AddCommand(
		composeCreateCmd,
		startAppCmd,
		startSrvCmd)

	stopCmd.AddCommand(
		stopAppCmd,
		stopSrvCmd)

	resetCmd.AddCommand(
		resetAppCmd)
}
