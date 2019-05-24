package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	srvListCmd = &cobra.Command{
		Use:   "srv:list",
		Short: "list accessible services",
		Long:  `This subcommand can pageing the services.`,
		Run:   srvList}

	srvInfoCmd = &cobra.Command{
		Use:   "srv:info",
		Short: "view info about a service",
		Long:  `This subcommand prints info about the current service`,
		Run:   srvInfo}

	srvCreateCmd = &cobra.Command{
		Use:   "srv:create",
		Short: "creates a new servive",
		Long: `Creates a new servive.
		if no <name> is provided, one will be generated automatically.`,
		Run: srvCreate}

	srvChangePlaneCmd = &cobra.Command{
		Use:   "srv:plan",
		Short: "change the Plan of service",
		Long: `set Plan for an service.
		This limit isn't applied to each individual pod, 
		so setting a plan for an service means that 
		each pod can gets more resourse and overused pay per consume.`,
		Run: srvChangePlane}

	srvPortforwardCmd = &cobra.Command{
		Use:   "srv:portforward",
		Short: "portforward to connect to an application running in a cluster",
		Long: `Portforward to connect to an application running in a cluster.
		This type of connection can be useful for database debugging`,
		Run: srvPortforward}

	srvStartCmd = &cobra.Command{
		Use:   "srv:start",
		Short: "start the stopped service",
		Long:  `This subcommand start the stopped service.`,
		Run:   srvStart}

	srvStopCmd = &cobra.Command{
		Use:   "srv:stop",
		Short: "stop the running service",
		Long:  `This subcommand stop the running service.`,
		Run:   srvStop}

	srvDestroyCmd = &cobra.Command{
		Use:   "srv:destroy",
		Short: "destroy an service",
		Long:  `This subcommand destroy an service.`,
		Run:   srvDestroy}

	srvAttachDomainCmd = &cobra.Command{
		Use:   "dom:srv-attach",
		Short: "attach domain to service",
		Long:  `This subcommand attach domain to a service`,
		Run:   srvAttachDomain}

	srvDetachDomainCmd = &cobra.Command{
		Use:   "dom:srv-detach",
		Short: "detach domain of the service",
		Long:  `This subcommand detach domain of the service`,
		Run:   srvDetachDomain}

	srvConfigSetCmd = &cobra.Command{
		Use:   "srv:cfg-set",
		Short: "sets configuration variables for a service",
		Long:  `This subcommand sets configuration variables for a service.`,
		Run:   srvConfigSet}

	srvConfigUnsetCmd = &cobra.Command{
		Use:   "srv:cfg-unset",
		Short: "unsets configuration variables for a service",
		Long:  `This subcommand unsets configuration variables for a service.`,
		Run:   srvConfigUnset}
)

func srvList(cmd *cobra.Command, args []string) {
	req := reqIndexForApp(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvList(client.Context(), req)
	uiCheckErr("Could not List the Services: %v", err)
	uiList(res)
}

func srvInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvInfo(client.Context(), req)
	uiCheckErr("Could not Get Service: %v", err)
	uiServicStatus(res)
}

func srvPortforward(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func srvStart(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvStart(client.Context(), req)
	uiCheckErr("Could not Start the Service: %v", err)
	uiServicStatus(res)
}

func srvStop(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvStop(client.Context(), req)
	uiCheckErr("Could not Stop the Service: %v", err)
	uiServicStatus(res)
}

func srvDestroy(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().SrvDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Service: %v", err)
	log.Printf("service %s deleted", req.Name)
}

func srvCreate(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.ProductName = cmd.Flag("product").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()
	req.Variable = arrayFlagToMap(flagVariableArray)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvCreate(client.Context(), req)
	uiCheckErr("Could not Create the Service: %v", err)
	uiServicStatus(res)
}

func srvConfigSet(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvConfigSetReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Variables = arrayFlagToMap(flagVariableArray)
	client := grpcConnect()
	defer client.Close()

	res, err := client.V2().SrvConfigSet(client.Context(), req)
	uiCheckErr("Could not Set the Config for Service: %v", err)
	uiServicStatus(res)
}

func srvConfigUnset(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvConfigUnset(client.Context(), req)
	uiCheckErr("Could not Unset the Config for Service: %v", err)
	uiServicStatus(res)
}

func srvChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan: %v", err)
	uiServicStatus(res)
}

func srvAttachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("name").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("attachment").Value.String()
	req.Endpoint = cmd.Flag("endpoint").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvAttachDomain(client.Context(), req)
	uiCheckErr("Could not Attach the Domain for Service: %v", err)
	uiServicStatus(res)
}

func srvDetachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("name").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("attachment").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Service: %v", err)
	uiServicStatus(res)
}

func init() {
	// srv List:
	srvListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number list")
	srvListCmd.Flags().StringVarP(&flagAppName, "app", "n", "", "page number list")

	// srv Info:
	srvInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvInfoCmd.MarkFlagRequired("name")

	// srv Create:
	srvCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvCreateCmd.Flags().StringP("product", "p", "", "name of product")
	srvCreateCmd.Flags().StringP("plan", "P", "", "the plan of sell")
	srvCreateCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "variable of service")
	srvCreateCmd.MarkFlagRequired("name")
	srvCreateCmd.MarkFlagRequired("product")

	// srv Change Plan:
	srvChangePlaneCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvChangePlaneCmd.Flags().StringP("plan", "p", "", "define the new plan of service")
	srvChangePlaneCmd.MarkFlagRequired("name")
	srvChangePlaneCmd.MarkFlagRequired("plan")

	// srv Config Set:
	srvConfigSetCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service.")
	srvConfigSetCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of the service")
	srvConfigSetCmd.MarkFlagRequired("name")
	srvConfigSetCmd.MarkFlagRequired("variable")

	// srv Config Unset:
	srvConfigUnsetCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service.")
	srvConfigUnsetCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of the service")
	srvConfigUnsetCmd.MarkFlagRequired("name")
	srvConfigUnsetCmd.MarkFlagRequired("variable")

	// srv Portforward:
	srvPortforwardCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvPortforwardCmd.MarkFlagRequired("name")

	// srv Start:
	srvStartCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvStartCmd.MarkFlagRequired("name")

	// srv Stop:
	srvStopCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvStopCmd.MarkFlagRequired("name")

	// srv Destroy:
	srvDestroyCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service")
	srvDestroyCmd.MarkFlagRequired("name")

	// srv Attach Domain:
	srvAttachDomainCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service.")
	srvAttachDomainCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	srvAttachDomainCmd.Flags().StringP("endpoint", "e", "", "name of the service endpoint")
	srvAttachDomainCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")
	srvAttachDomainCmd.MarkFlagRequired("name")
	srvAttachDomainCmd.MarkFlagRequired("endpoint")
	srvAttachDomainCmd.MarkFlagRequired("path")
	srvAttachDomainCmd.MarkFlagRequired("attachment")

	// srv Detach Domain:
	srvDetachDomainCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service.")
	srvDetachDomainCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	srvDetachDomainCmd.Flags().StringP("path", "p", "", "http subpath to route traffic")
	srvDetachDomainCmd.MarkFlagRequired("path")
	srvDetachDomainCmd.MarkFlagRequired("name")
	srvDetachDomainCmd.MarkFlagRequired("attachment")

	// add service subcommands to root
	rootCmd.AddCommand(
		srvListCmd,
		srvInfoCmd,
		srvCreateCmd,
		srvChangePlaneCmd,
		srvConfigSetCmd,
		srvConfigUnsetCmd,
		srvPortforwardCmd,
		srvStartCmd,
		srvStopCmd,
		srvDestroyCmd,
		srvAttachDomainCmd,
		srvDetachDomainCmd)
}
