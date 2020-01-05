package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:   "rm [command] [name]",
		Short: "Remove by name for [service|application|domain|volume|worker|config]",
		Long:  `This subcommand can Remove the [services|app|domain|...] by name`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	rmAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "destroy an application",
		Long:  `This subcommand destroy an application.`,
		Run:   appDestroy}

	rmSrvCmd = &cobra.Command{
		Use:   "service [name]",
		Short: "destroy an service",
		Long:  `This subcommand destroy an service.`,
		Run:   srvDestroy}

	rmDomCmd = &cobra.Command{
		Use:   "domain [name]",
		Short: "destroy an domain",
		Long:  `This subcommand destroy an domain.`,
		Run:   domainDelete}

	rmVolCmd = &cobra.Command{
		Use:   "volume [name]",
		Short: "destroy an volume",
		Long:  `This subcommand destroy an volume.`,
		Run:   volumeDelete}

	rmImgCmd = &cobra.Command{
		Use:   "image [name]",
		Short: "destroy an image",
		Long:  `This subcommand destroy an image.`,
		Run:   imgDelete}
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.AddCommand(
		rmAppCmd,
		rmSrvCmd,
		rmDomCmd,
		rmVolCmd,
		rmImgCmd)
}
