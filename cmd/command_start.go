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
		Args:  cobra.MinimumNArgs(1)}

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
		Args:  cobra.MinimumNArgs(1)}

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
		Args:  cobra.MinimumNArgs(1)}

	resetAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "reset the running application",
		Long:  `Destroys an application instance and create a new instance of application.`,
		Run:   appReset}

	/*
		resetSrvCmd = &cobra.Command{
			Use:   "service [name]",
			Short: "reset the running service",
			Long:  `reset the running Service.`,
			Run:   srvReset}
	*/
)

func init() {
	rootCmd.AddCommand(
		startCmd,
		stopCmd,
		resetCmd)

	startCmd.AddCommand(
		startAppCmd,
		startSrvCmd)

	stopCmd.AddCommand(
		stopAppCmd,
		stopSrvCmd)

	resetCmd.AddCommand(
		resetAppCmd)
}
