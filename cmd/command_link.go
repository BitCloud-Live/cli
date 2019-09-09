package cmd

import (
	"github.com/spf13/cobra"
)

var (
	linkCmd = &cobra.Command{
		Use:   "link [command]",
		Short: "link [set|unset] between: [Volume <-> Application <-> Service <-> Domain]",
		Long: `link [set|unset] between: 
		[Application <-> Service]
		[Application <-> Volume]
		[Application <-> Domain]
		[Service     <-> Domain]`,
		Example: `
  $: yb link set \
        --service="my-wordpress" \
        --domain="example.com"`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	setLinkCmd = &cobra.Command{
		Use:   "set",
		Run:   setDetect,
		Short: "creates a new link",
		Long:  `creates a new link.`,
		Example: `
  $: yb link set \
        --service="my-wordpress" \
		--domain="example.com"`}

	unsetLinkCmd = &cobra.Command{
		Use:   "unset",
		Run:   unsetDetect,
		Short: "destroy a link",
		Long:  `destroy a link.`,
		Example: `
  $: yb link unset \
        --service="my-wordpress" \
        --domain="example.com"`}
)

func setDetect(cmd *cobra.Command, args []string) {
	app := cmd.Flag("application").Value.String()
	srv := cmd.Flag("service").Value.String()
	vol := cmd.Flag("volume").Value.String()
	dom := cmd.Flag("domain").Value.String()

	if app != "" {
		if dom != "" {
			appAttachDomain(cmd, args)
		} else if vol != "" {
			appAttachVolume(cmd, args)
		} else if srv != "" {
			appSrvBind(cmd, args)
		} else {
			log.Fatalf("Err")
		}

	} else if srv != "" && dom != "" {
		srvAttachDomain(cmd, args)
	} else {
		log.Fatalf("Err")
	}
}

func unsetDetect(cmd *cobra.Command, args []string) {
	app := cmd.Flag("application").Value.String()
	srv := cmd.Flag("service").Value.String()
	vol := cmd.Flag("volume").Value.String()
	dom := cmd.Flag("domain").Value.String()

	if app != "" {
		if dom != "" {
			appDetachDomain(cmd, args)
		} else if vol != "" {
			appDetachVolume(cmd, args)
		} else if srv != "" {
			appSrvUnBind(cmd, args)
		} else {
			log.Fatalf("Err")
		}

	} else if srv != "" && dom != "" {
		srvDetachDomain(cmd, args)
	} else {
		log.Fatalf("Err")
	}
}

func init() {
	// set:
	setLinkCmd.Flags().StringP("application", "a", "", "the uniquely identifiable name for the Application.")
	setLinkCmd.Flags().StringP("service", "s", "", "the uniquely identifiable name for the Service.")
	setLinkCmd.Flags().StringP("volume", "v", "", "name of Volume for attachment")
	setLinkCmd.Flags().StringP("domain", "d", "", "name of Domain for attachment")

	// app/srv </- Domain
	// app     <-  Volume
	setLinkCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")

	// srv/app </- Domain
	setLinkCmd.Flags().StringP("endpoint", "e", "", "name of the service endpoint")
	unsetLinkCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")

	// unset:
	unsetLinkCmd.Flags().StringP("application", "a", "", "the uniquely identifiable name for the Application.")
	unsetLinkCmd.Flags().StringP("service", "s", "", "the uniquely identifiable name for the Service.")
	unsetLinkCmd.Flags().StringP("volume", "v", "", "name of Volume for Detattachment")
	unsetLinkCmd.Flags().StringP("domain", "d", "", "name of Domain for Detattachment")

	// srv/app </- Domain
	unsetLinkCmd.Flags().StringP("endpoint", "e", "", "name of the service endpoint")

	rootCmd.AddCommand(linkCmd)
	linkCmd.AddCommand(
		setLinkCmd,
		unsetLinkCmd)
}
