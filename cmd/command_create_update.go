package cmd

import (
	"github.com/spf13/cobra"
)

// Create Command
var (
	createCmd = &cobra.Command{
		Use:   "create [command]",
		Short: "creates new [service|application|domain|volume|worker|config]",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flag("compose-file").Changed {
				//NOOP
				// composeCreate(cmd, args)
			} else {
				cmd.Help()
			}
		}}

	appCreateCmd = &cobra.Command{
		Use:   "application",
		Short: "creates a new application",
		Long: `Creates a new application.
			if no <name> is provided, one will be generated automatically.`,
		Run: appCreate}

	srvCreateCmd = &cobra.Command{
		Use:   "service [PRODUCT.name]",
		Run:   srvCreate,
		Short: "creates a new servive from product list",
		Long: `Creates a new servive.
  if no 'name' is provided, one will be generated automatically.`,
		Example: `
  $: yb product
  $: yb create service mysql \
        --name=db \
        --plan=starter \
        --variable="Database.password=DaRHEm@DaX" \
	    --variable="Database.user=root"`}

	domainCreateCmd = &cobra.Command{
		Use:   "domain [domain]",
		Run:   domainCreate,
		Short: "creates a new Domian",
		Long:  `Creates a new Domian`,
		Example: `
  $: yb create domain example.com \
		--TLS=true`}

	volumeCreateCmd = &cobra.Command{
		Use:   "volume",
		Short: "create new volume",
		Long:  `create new volume`,
		Run:   volumeCreate}

	bucketCreateCmd = &cobra.Command{
		Use:   "bucket [name]",
		Short: "create an Object Storage bucket over the YOTTAb",
		Long:  `create an Object Storage bucket over the YOTTAb`,
		Run:   bucketCreate}
)

// Update Command
var (
	updateCmd = &cobra.Command{
		Use:   "update [APP.name]",
		Short: "update the config of existing application",
		Long:  `This subcommand Update an existing application.`,
		Run:   appUpdate}
)

func init() {
	createCmd.Flags().StringP("compose-file", "f", "", "the path of compose file")

	// service Create flag:
	srvCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvCreateCmd.Flags().StringP("plan", "", "starter", "the plan of sell")
	srvCreateCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "variable of service")
	srvCreateCmd.MarkFlagRequired("name")

	// application Create flag:
	appCreateCmd.Flags().StringP("plan", "s", "", "name of plan")
	appCreateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	appCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	appCreateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	appCreateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc, default to http")
	appCreateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 1, "min scale of application")
	appCreateCmd.MarkFlagRequired("image")

	// domain create:
	domainCreateCmd.Flags().BoolVar(&flagTLS, "TLS", false, "enable TLS for domain")

	// volume create:
	volumeCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeCreateCmd.Flags().StringP("volume-type", "v", "", "the type of volume")
	volumeCreateCmd.MarkFlagRequired("name")
	volumeCreateCmd.MarkFlagRequired("volume-type")

	// Application Update
	updateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	updateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	updateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 0, "min scale of application")
	updateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc, default to http")
	updateCmd.Flags().StringArrayVarP(&flagVariableArray, "routes", "r", nil, "Routes of application")

	rootCmd.AddCommand(
		createCmd,
		updateCmd)

	createCmd.AddCommand(
		appCreateCmd,
		srvCreateCmd,
		domainCreateCmd,
		volumeCreateCmd,
		bucketCreateCmd)
}
