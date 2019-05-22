package cmd

import (
	"context"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagVarPort         uint64
	flagVarMinScale     uint64
	minScale            uint64
	flagVarImage        string
	flagVarEndpointType string
)

var (
	appListCmd = &cobra.Command{
		Use:   "app:list",
		Short: "list accessible applications",
		Long:  `Lists applications visible to the current user.`,
		Run:   appList}

	appInfoCmd = &cobra.Command{
		Use:   "app:info",
		Short: "view info about an application",
		Long:  `This subcommand prints info about the current application`,
		Run:   appInfo}

	appOpenCmd = &cobra.Command{
		Use:   "app:open",
		Short: "Open an application in the default browser",
		Long:  `This subcommand will open an application in the default browser in every os.`,
		Run:   appOpen}

	appLogCmd = &cobra.Command{
		Use:   "app:log",
		Short: "tail application log",
		Long:  `This subcommand tails the current application logs`,
		Run:   appLog}

	appFTPMountCmd = &cobra.Command{
		Use:   "app:ftp-mount",
		Short: "Connect to the remote file system using ftp protocol",
		Long:  `This subcommand connects to the remote file system using ftp protocol`,
		Run:   appFTPMount}

	appCreateCmd = &cobra.Command{
		Use:   "app:create",
		Short: "creates a new application",
		Long: `Creates a new application.
		if no <name> is provided, one will be generated automatically.`,
		Run: appCreate}

	appConfigSetCmd = &cobra.Command{
		Use:   "app:cfg-set",
		Short: "sets configuration variables for an application",
		Long:  `This subcommand sets configuration variables for an application.`,
		Run:   appConfigSet}

	appConfigUnsetCmd = &cobra.Command{
		Use:   "app:cfg-unset",
		Short: "unsets configuration variables for an application",
		Long:  `This subcommand unsets configuration variables for an application.`,
		Run:   appConfigUnset}

	appAddEnvCmd = &cobra.Command{
		Use:   "app:env-set",
		Short: "sets environment variables for an application",
		Long: `Sets environment variables for an application.
		Usage: yb app:env-set <var>=<value> [<var>=<value>...] [options]

		Arguments:
		  <var>
			  the uniquely identifiable name for the environment variable.
		  <value>
			  the value of said environment variable.`,
		Run: appAddEnvironmentVariable}

	appRemoveEnvCmd = &cobra.Command{
		Use:   "app:env-unset",
		Short: "unset environment variables for an application",
		Long: `Unset environment variables for an application.
		Usage: deis app:env-unset <key>... [options]

		Arguments:
		  <key>
		    the variable to remove from the application's environment.`,
		Run: appRemoveEnvironmentVariable}

	appChangePlaneCmd = &cobra.Command{
		Use:   "app:plan",
		Short: "change the Plan of application",
		Long: `set Plan for an application.
		This limit isn't applied to each individual pod, 
		so setting a plan for an application means that 
		each pod can gets more resourse and overused pay per consume.`,
		Run: appChangePlane}

	appResetCmd = &cobra.Command{
		Use:   "app:restart",
		Short: "restart of application",
		Long:  `Destroys an application instance and create a new instance of application.`,
		Run:   appReset}

	appStartCmd = &cobra.Command{
		Use:   "app:start",
		Short: "start the stopped application",
		Long:  `This subcommand start the stopped application.`,
		Run:   appStart}

	appStopCmd = &cobra.Command{
		Use:   "app:stop",
		Short: "stop the running application",
		Long:  `This subcommand stop the running application.`,
		Run:   appStop}

	appDestroyCmd = &cobra.Command{
		Use:   "app:destroy",
		Short: "destroy an application",
		Long:  `This subcommand destroy an application.`,
		Run:   appDestroy}

	appSrvBindCmd = &cobra.Command{
		Use:   "app:link",
		Short: "add link to another service",
		Long: `This subcommand add link to another service
		when starting a new application container in the cluster, 
		then the application can access the service via a private networking interface.`,
		Run: appSrvBind}

	appSrvUnBindCmd = &cobra.Command{
		Use:   "app:unlink",
		Short: "remove link to binded service and restart application",
		Long:  `This subcommand remove link to binded service and restart application.`,
		Run:   appSrvUnBind}

	appAttachVolumeCmd = &cobra.Command{
		Use:   "vol:attach",
		Short: "attach volume to application",
		Long:  `This subcommand attach mounted volume.`,
		Run:   appAttachVolume}

	appDetachVolumeCmd = &cobra.Command{
		Use:   "vol:detach",
		Short: "detach volume of the application",
		Long:  `This subcommand detach volume of the application.`,
		Run:   appDetachVolume}

	appAttachDomainCmd = &cobra.Command{
		Use:   "dom:attach",
		Short: "attach domain to application",
		Long:  `This subcommand attach domain to application`,
		Run:   appAttachDomain}

	appDetachDomainCmd = &cobra.Command{
		Use:   "dom:detach",
		Short: "detach domain of the application",
		Long:  `This subcommand detach domain of the application`,
		Run:   appDetachDomain}
)

func appList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppList(client.Context(), req)
	uiCheckErr("Could not List the Applications: %v", err)
	uiList(res)
}

func appInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiApplicationStatus(res)
}

func appOpen(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiApplicationOpen(res)
}

func appLog(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	logClient, err := client.V2().AppLog(context.Background(), req)
	if err != nil {
		panic(err)
	}
	uiCheckErr("Could not Get Application log: %v", err)
	uiApplicationLog(logClient)
}

func appFTPMount(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppFTPPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
	uiNFSMount(res)
}

func appStart(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppStart(client.Context(), req)
	uiCheckErr("Could not Start the Application: %v", err)
	uiApplicationStatus(res)
}

func appStop(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppStop(client.Context(), req)
	uiCheckErr("Could not Stop the Application: %v", err)
	uiApplicationStatus(res)
}

func appDestroy(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().AppDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Application: %v", err)
	log.Printf("app %s deleted", req.Name)
}

func appCreate(cmd *cobra.Command, args []string) {
	if err := endpointTypeValid(flagVarEndpointType); err != nil {
		log.Panic(err)
	}
	req := new(ybApi.AppCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()
	req.Config = new(ybApi.AppConfig)
	req.Config.Port = flagVarPort
	req.Config.EndpointType = flagVarEndpointType
	req.Config.MinScale = flagVarMinScale
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppCreate(client.Context(), req)
	uiCheckErr("Could not Create the Application: %v", err)
	uiApplicationStatus(res)
}

func appChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan: %v", err)
	uiApplicationStatus(res)
}

func appConfigSet(cmd *cobra.Command, args []string) {
	req := new(ybApi.ConfigSetReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Config = new(ybApi.AppConfig)
	req.Config.MinScale = flagVarMinScale
	req.Config.Port = flagVarPort
	req.Config.Image = flagVarImage
	client := grpcConnect()
	defer client.Close()

	res, err := client.V2().AppConfigSet(client.Context(), req)
	uiCheckErr("Could not Set the Config for Application: %v", err)
	uiApplicationStatus(res)
}

func appConfigUnset(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppConfigUnset(client.Context(), req)
	uiCheckErr("Could not Unset the Config for Application: %v", err)
	uiApplicationStatus(res)
}

func appAddEnvironmentVariable(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppAddEnvironmentVariableReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Variables = arrayFlagToMap(flagVariableArray)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAddEnvironmentVariable(client.Context(), req)
	uiCheckErr("Could not Add the Environment Variable for Application: %v", err)
	uiApplicationStatus(res)
}

func appRemoveEnvironmentVariable(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppRemoveEnvironmentVariable(client.Context(), req)
	uiCheckErr("Could not Remove the Environment Variable for Application: %v", err)
	uiApplicationStatus(res)
}

func appReset(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppReset(client.Context(), req)
	uiCheckErr("Could not Reset the Application: %v", err)
	uiApplicationStatus(res)
}

func appSrvBind(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Service = cmd.Flag("service").Value.String()
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppSrvBind(client.Context(), req)
	uiCheckErr("Could not Bind the Service for Application: %v", err)
	uiApplicationStatus(res)
}

func appSrvUnBind(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Service = cmd.Flag("service").Value.String()
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppSrvUnBind(client.Context(), req)
	uiCheckErr("Could not Unbind the Service for Application: %v", err)
	uiApplicationStatus(res)
}

func appAttachVolume(cmd *cobra.Command, args []string) {
	req := new(ybApi.VolumeAttachReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()
	req.MountPath = cmd.Flag("mount-path").Value.String()

	client := grpcConnect()

	defer client.Close()
	res, err := client.V2().AppAttachVolume(client.Context(), req)
	uiCheckErr("Could not Attach the Volume for Application: %v", err)
	uiApplicationStatus(res)
}

func appDetachVolume(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachVolume(client.Context(), req)
	uiCheckErr("Could not Detach the Volume for Application: %v", err)
	uiApplicationStatus(res)
}

func appAttachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("name").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("attachment").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAttachDomain(client.Context(), req)
	uiCheckErr("Could not Attach the Domain for Application: %v", err)
	uiApplicationStatus(res)
}

func appDetachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("name").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("attachment").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Application: %v", err)
	uiApplicationStatus(res)
}

func init() {
	// app List:
	appListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")

	// app Info:
	appInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appInfoCmd.MarkFlagRequired("name")

	// app Open:
	appOpenCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appOpenCmd.MarkFlagRequired("name")

	// app log:
	appLogCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appLogCmd.MarkFlagRequired("name")

	// app ftp mount:
	appFTPMountCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appFTPMountCmd.MarkFlagRequired("name")

	// app Create:
	appCreateCmd.Flags().StringP("plan", "s", "", "name of plan")
	appCreateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	appCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	appCreateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	appCreateCmd.Flags().StringVarP(&flagVarEndpointType, "endpoint-type", "e", "http", "Accepted values: http|grpc, default to http")
	appCreateCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 1, "min scale of application")
	appCreateCmd.MarkFlagRequired("image")

	// app Config Set:
	appConfigSetCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appConfigSetCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	appConfigSetCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	appConfigSetCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 0, "min scale of application")
	appConfigSetCmd.MarkFlagRequired("name")

	// app Config Unset:
	appConfigUnsetCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appConfigUnsetCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 8080, "port of application")
	appConfigUnsetCmd.Flags().StringP("image", "i", "", "image of application")
	appConfigUnsetCmd.Flags().Uint64VarP(&flagVarMinScale, "min-scale", "m", 1, "min scale of application")
	appConfigUnsetCmd.MarkFlagRequired("name")

	// app Change Plan:
	appChangePlaneCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appChangePlaneCmd.Flags().StringP("plan", "p", "", "define the new plan of application")
	appChangePlaneCmd.MarkFlagRequired("name")

	// app Start:
	appStartCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appStartCmd.MarkFlagRequired("name")

	// app Stop:
	appStopCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appStopCmd.MarkFlagRequired("name")

	// app Reset:
	appResetCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appResetCmd.MarkFlagRequired("name")

	// app Destroy:
	appDestroyCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appDestroyCmd.MarkFlagRequired("name")

	// app Add Environment Variable
	appAddEnvCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appAddEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	appAddEnvCmd.MarkFlagRequired("name")
	appAddEnvCmd.MarkFlagRequired("variable")

	// app Remove Environment Variable
	appRemoveEnvCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appRemoveEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	appRemoveEnvCmd.MarkFlagRequired("name")
	appRemoveEnvCmd.MarkFlagRequired("variable")

	// app Service Bind:
	appSrvBindCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appSrvBindCmd.Flags().StringP("service", "s", "", "name of service")
	appSrvBindCmd.MarkFlagRequired("name")
	appSrvBindCmd.MarkFlagRequired("service")

	// app Service UnBind:
	appSrvUnBindCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appSrvUnBindCmd.Flags().StringP("service", "s", "", "name of service")
	appSrvUnBindCmd.MarkFlagRequired("name")
	appSrvUnBindCmd.MarkFlagRequired("service")

	// app Aettach Domain:
	appAttachDomainCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appAttachDomainCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	appAttachDomainCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")
	appAttachDomainCmd.MarkFlagRequired("name")
	appAttachDomainCmd.MarkFlagRequired("path")
	appAttachDomainCmd.MarkFlagRequired("attachment")

	// app Detach Domain:
	appDetachDomainCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appDetachDomainCmd.Flags().StringP("attachment", "a", "", "name of detachment")
	appDetachDomainCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")
	appDetachDomainCmd.MarkFlagRequired("name")
	appDetachDomainCmd.MarkFlagRequired("path")
	appDetachDomainCmd.MarkFlagRequired("attachment")

	// app Aettach Volume:
	appAttachVolumeCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appAttachVolumeCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	appAttachVolumeCmd.Flags().StringP("mount-path", "m", "", "name of attachment")
	appAttachVolumeCmd.MarkFlagRequired("name")
	appAttachVolumeCmd.MarkFlagRequired("attachment")
	appAttachVolumeCmd.MarkFlagRequired("mount-path")

	// app Detach Volume:
	appDetachVolumeCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	appDetachVolumeCmd.Flags().StringP("attachment", "a", "", "name of detachment")
	appDetachVolumeCmd.MarkFlagRequired("name")
	appDetachVolumeCmd.MarkFlagRequired("attachment")

	rootCmd.AddCommand(
		appListCmd,
		appInfoCmd,
		appOpenCmd,
		appLogCmd,
		appFTPMountCmd,
		appCreateCmd,
		appConfigSetCmd,
		appConfigUnsetCmd,
		appAddEnvCmd,
		appRemoveEnvCmd,
		appChangePlaneCmd,
		appResetCmd,
		appStartCmd,
		appStopCmd,
		appDestroyCmd,
		appSrvBindCmd,
		appSrvUnBindCmd,
		appAttachDomainCmd,
		appDetachDomainCmd,
		appAttachVolumeCmd,
		appDetachVolumeCmd)
}
