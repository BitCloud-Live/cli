package cmd

import (
	"github.com/spf13/cobra"
)

// Create Command
var (
	createCmd = &cobra.Command{
		Use:   "create [command]",
		Short: "creates new [service|application|domain|volume]",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	appCreateCmd = &cobra.Command{
		Use:   "application",
		Run:   appCreate,
		Short: "creates a new application",
		Long:  `Creates a new application.`,
		Example: `
  $: yb create application \
        --image="hub.yottab.io/test/dotnetcore:aspnetapp" \
        --name="myaspnetapp" \
		--port=80
		

  $: yb create application \
        --image="hub.yottab.io/test/myworker:v1.2.3" \
		--name="myaspnetapp" \
		--min-scale=4 \
		--debug=true \
        --port=8080 \
        --endpoint-type=private`}

	srvCreateCmd = &cobra.Command{
		Use:   "service [PRODUCT.name]",
		Run:   srvCreate,
		Short: "creates a new servive from product list",
		Long: `Creates a new servive.
  if no 'name' is provided, one will be generated automatically.`,
		Example: `
  ## Get a list of Service Plan and Required Variables
  $: yb product service mysql

  $: yb create service mysql \
        --name=mydb \
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
		Run:   volumeCreate,
		Short: "create new volume",
		Long:  `create new volume`,
		Example: `
  ## Get a list of Volume Plans 
  $: yb product volume

  $: yb create volume \
		--name="my-volume" \
		--volume-type="persistant-2Gi"`}

	bucketCreateCmd = &cobra.Command{
		Use:   "bucket [name]",
		Run:   bucketCreate,
		Short: "create an Object Storage bucket over the YOTTAb",
		Long:  `create an Object Storage bucket over the YOTTAb`}
)

// Update Command
var (
	updateCmd = &cobra.Command{
		Use:   "update [APP.name]",
		Run:   appUpdate,
		Short: "update the config of existing application",
		Long:  `This subcommand Update an existing application.`,
		Example: `
  ## The Application Version is updated
  $: yb update myadmin \
        --image="hub.yottab.io/test/myadmin:v1.2.3"
		

  ## The Application Scale is updated
  $: yb update myadmin \
		--min-scale=4`}
)

func init() {
	createCmd.Flags().StringP("compose-file", "f", "", "the path of compose file")

	// service Create flag:
	srvCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvCreateCmd.Flags().StringP("plan", "", "starter", "the plan of sell")
	srvCreateCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "variable of service")
	srvCreateCmd.MarkFlagRequired("name")

	// application Create flag:
	appCreateCmd.Flags().StringP("plan", "", "default", "name of plan")
	appCreateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	appCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	appCreateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	appCreateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc|private")
	appCreateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 1, "min scale of application")
	appCreateCmd.MarkFlagRequired("name")
	appCreateCmd.MarkFlagRequired("image")

	// domain create:
	domainCreateCmd.Flags().BoolVar(&flagTLS, "TLS", true, "enable TLS for domain")
	appCreateCmd.MarkFlagRequired("TLS")

	// volume create:
	volumeCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeCreateCmd.Flags().StringP("volume-type", "v", "", "the type of volume")
	volumeCreateCmd.MarkFlagRequired("name")
	volumeCreateCmd.MarkFlagRequired("volume-type")

	// Application Update
	updateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	updateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	updateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 0, "min scale of application")
	updateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc|private")

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
