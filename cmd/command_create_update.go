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
				composeCreate(cmd, args)
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

	workerCreateCmd = &cobra.Command{
		Use:   "worker [APP.name]",
		Short: "Add a new worker to an application",
		Long:  `This subcommand adds a new worker to an application, be notify that there is a upper limit to worker counts.`,
		Run:   workerCreate}

	volumeCreateCmd = &cobra.Command{
		Use:   "volume",
		Short: "create new volume",
		Long:  `create new volume`,
		Run:   volumeCreate}
)

// Update Command
var (
	updateCmd = &cobra.Command{
		Use:   "update [command] [name]",
		Short: "update the existing [application|worker|plan]",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1)}

	appUpdateCmd = &cobra.Command{
		Use:   "application [APP.name]",
		Short: "Update an existing application",
		Long:  `This subcommand Update an existing application.`,
		Run:   appUpdate}

	workerUpdateCmd = &cobra.Command{
		Use:   "worker [APP.name]",
		Short: "Update an existing worker for an application",
		Long:  `This subcommand Update an existing worker for an application.`,
		Run:   workerUpdate}

	planeUpdateCmd = &cobra.Command{
		Use:   "plan [type]",
		Short: "change the Plan of [service|application]",
		Long: `set Plan for an service.
		This limit isn't applied to each individual pod, 
		so setting a plan for an service means that 
		each pod can gets more resourse and overused pay per consume.`}

	srvPlaneUpdateCmd = &cobra.Command{
		Use:   "service [Srv.name] [Plan.name]",
		Run:   srvChangePlane,
		Short: "change the Plan of service",
		Long: `set Plan for an service.
			This limit isn't applied to each individual pod, 
			so setting a plan for an service means that 
			each pod can gets more resourse and overused pay per consume.`}

	appPlaneUpdateCmd = &cobra.Command{
		Use:   "application [App.name] [Plan.name]",
		Run:   appChangePlane,
		Short: "change the Plan of application",
		Long: `set Plan for an application.
				This limit isn't applied to each individual pod, 
				so setting a plan for an application means that 
				each pod can gets more resourse and overused pay per consume.`}
)

func init() {
	createCmd.Flags().StringP("compose-file", "f", "", "the path of compose file")

	// service Create flag:
	srvCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvCreateCmd.Flags().StringP("plan", "P", "starter", "the plan of sell")
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

	// worker Create flag:
	workerCreateCmd.Flags().StringP("name", "n", "", "name of worker, a uniquely identifiable name for the worker. No other app can already exist with this name.")
	workerCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	workerCreateCmd.Flags().StringP("image", "i", "", "image of application")
	workerCreateCmd.MarkFlagRequired("image")
	workerCreateCmd.MarkFlagRequired("port")
	workerCreateCmd.MarkFlagRequired("name")

	// Application Update
	appUpdateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	appUpdateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	appUpdateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 0, "min scale of application")
	appUpdateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc, default to http")
	appUpdateCmd.Flags().StringArrayVarP(&flagVariableArray, "routes", "r", nil, "Routes of application")

	// worker Update:
	workerUpdateCmd.Flags().StringP("name", "n", "", "name of worker, a uniquely identifiable name for the worker. No other app can already exist with this name.")
	workerUpdateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	workerUpdateCmd.Flags().StringP("image", "i", "", "image of application")
	workerUpdateCmd.MarkFlagRequired("name")

	// volume create:
	volumeCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeCreateCmd.Flags().StringP("volume-type", "v", "", "the type of volume")
	volumeCreateCmd.MarkFlagRequired("name")
	volumeCreateCmd.MarkFlagRequired("volume-type")

	rootCmd.AddCommand(
		createCmd,
		updateCmd)

	createCmd.AddCommand(
		appCreateCmd,
		srvCreateCmd,
		domainCreateCmd,
		workerCreateCmd,
		volumeCreateCmd)

	updateCmd.AddCommand(
		workerUpdateCmd,
		appUpdateCmd,
		planeUpdateCmd)

	planeUpdateCmd.AddCommand(
		srvPlaneUpdateCmd,
		appPlaneUpdateCmd)
}
