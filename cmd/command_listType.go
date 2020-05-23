package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Applications
var (
	appCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "list or informaition of accessible applications",
		Long:  `This subcommand can pageing the applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				appInfo(cmd, args)
			} else {
				appList(cmd, args)
			}
		}}
)

// Service
var (
	srvCmd = &cobra.Command{
		Use:   "service [service_name]",
		Short: "list or informaition of accessible services",
		Long: `This subcommand can pageing the Image.
  $: yb service                                 # list of all Services
  $: yb service [image_name]                    # detail of one Service
  $: yb service application [application_name]  # list of accessible Services for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				srvInfo(cmd, args)
			} else {
				srvList(cmd, args)
			}
		}}

	srvOneAppCmd = &cobra.Command{
		Use:   "application [application_name]",
		Short: "list of accessible Services for only one application",
		Long: `This subcommand can pageing the Image.
  $: yb service                                 # list of all Services
  $: yb service [image_name]                    # detail of one Service
  $: yb service application [application_name]  # list of accessible Services for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				srvList(cmd, args)
			} else {
				log.Println("$: yb service application [application_name]")
			}
		}}
)

// Domain
var (
	domCmd = &cobra.Command{
		Use:   "domain [domain_name]",
		Short: "list or informaition of accessible Domains",
		Long: `This subcommand can pageing the Domains.
  $: yb domain                                 # list of all domain
  $: yb domain [domain_name]                   # detail of one domain
  $: yb domain application [application_name]  # list of accessible domain for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				domainInfo(cmd, args)
			} else {
				domainList(cmd, args)
			}
		}}

	domOneAppCmd = &cobra.Command{
		Use:   "application [application_name]",
		Short: "list of accessible Domains for only one application",
		Long: `This subcommand can pageing the Domains.
  $: yb domain                                 # list of all Domains
  $: yb domain [domain_name]                   # detail of one Domain
  $: yb domain application [application_name]  # list of accessible Domains for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				domainList(cmd, args)
			} else {
				log.Println("$: yb domain application [application_name]")
			}
		}}
)

// Volume
var (
	volCmd = &cobra.Command{
		Use:   "volume [volume_name]",
		Short: "list or informaition of accessible Volume",
		Long: `This subcommand can pageing the Volume.
  $: yb volume                                 # list of all Volume
  $: yb volume [volume_name]                   # detail of one Volume
  $: yb volume application [application_name]  # list of accessible Volume for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				volumeInfo(cmd, args)
			} else {
				volumeList(cmd, args)
			}
		}}

	volOneAppCmd = &cobra.Command{
		Use:   "application [volume_name]",
		Short: "list of accessible Volumes for only one applicatio",
		Long: `This subcommand can pageing the Volume.
  $: yb volume                                 # list of all Volume
  $: yb volume [volume_name]                   # detail of one Volume
  $: yb volume application [application_name]  # list of accessible Volume for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				volumeList(cmd, args)
			} else {
				log.Println("$: yb volume application [application_name]")
			}
		}}
)

// Image
var (
	imgCmd = &cobra.Command{
		Use:   "image [image_name]",
		Short: "list or informaition of accessible Image",
		Long: `This subcommand can pageing the Image.
  $: yb image                                 # list of all image
  $: yb image [image_name]                    # detail of one image
  $: yb image application [application_name]  # list of accessible Image for only one application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				// detail of one image
				imgInfo(cmd, args)
			} else {
				// list all image for one user
				imgList(cmd, args)
			}
		}}

	imgOneAppCmd = &cobra.Command{
		Use:   "application [name]",
		Short: "list of accessible Image for only one application",
		Long:  `This subcommand can pageing the Image.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				imgList(cmd, args)
			} else {
				log.Println("$: yb image application [application_name]")
			}
		}}
)

// Product
var (
	prdCmd = &cobra.Command{
		Use:   "product [type]",
		Short: "list or informaition of accessible Products",
		Long: `This subcommand can pageing the Products.
		This subcommand show the information, praice and ... of a product.`}

	prdServiceCmd = &cobra.Command{
		Use:   "service [name]",
		Short: "list or informaition of accessible Service",
		Long: `This subcommand can pageing the Products.
			This subcommand show the information, praice and ... of a product.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				prdInfo(cmd, args)
			} else {
				prdList(cmd, args)
			}
		}}

	prdVolumeSpecListCmd = &cobra.Command{
		Use:   "volume [name]",
		Short: "list or informaition of accessible Volume",
		Long: `This subcommand can pageing the Volume.
		This subcommand show the information, Size and ... of a Volume.`,
		Run: func(cmd *cobra.Command, args []string) {
			//Deprecated
			// if len(args) == 1 {
			// 	volumeSpecInfo(cmd, args)
			// } else {
			volumeSpecList(cmd, args)
			// }
		}}
)

// Activities
var (
	activityCmd = &cobra.Command{
		Use:   "activity [tag] [name]",
		Run:   actList,
		Short: "show the USER.activities and available filter by [tag] and [name]",
		Long: `
  This subcommand shows list of all USER.activities.
  $: yb activity

  available filter by [tag] and [name]
  [tag]
      filter the list of USER.activities by "tag.id" in first arg: 
	  $: yb activity [tag_id]
	  
	  you can list all of 'TAG' and 'ID' by command <tags>:
	  $: yb activity tags

  [name]
      filter the list of USER.activities by "application.name" in secent arg: 
	  $: yb activity [tag] [name]`}

	activityTagListCmd = &cobra.Command{
		Use:     "tags",
		Run:     actTags,
		Short:   "show all available tags",
		Long:    `This subcommand shows list of available activity tags.`,
		Example: "$: yb activity tags"}

	logCmd = &cobra.Command{
		Use:   "log [APP.name]",
		Run:   appLog,
		Short: "tail application log",
		Long:  `This subcommand tails the current application logs`}
)

func init() {
	// Applications
	appCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Applications list")

	// Service
	srvCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Services list")
	srvOneAppCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Services list")

	// Domain
	domCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Domains list")
	domOneAppCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Domains list")

	// Volumes
	volCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Volumes list")

	// Images
	imgOneAppCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Volumes list")
	imgCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Images list")

	// Products
	prdServiceCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of products list")
	prdVolumeSpecListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of volume list")

	// Activites
	activityCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number of Activites list")

	rootCmd.AddCommand(
		appCmd,
		srvCmd,
		domCmd,
		volCmd,
		imgCmd,
		prdCmd,
		activityCmd,
		logCmd)

	prdCmd.AddCommand(
		prdServiceCmd,
		prdVolumeSpecListCmd)

	srvCmd.AddCommand(srvOneAppCmd)
	domCmd.AddCommand(domOneAppCmd)
	volCmd.AddCommand(volOneAppCmd)
	imgCmd.AddCommand(imgOneAppCmd)
	activityCmd.AddCommand(activityTagListCmd)
}
