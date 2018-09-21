package cmd

import (
	"github.com/spf13/cobra"
	uvApi "github.com/uvcloud/uv-api-go/proto"
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
		Short: "port-forward to connect to an application running in a cluster",
		Long: `Port-forward to connect to an application running in a cluster.
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
		Use:   "srv:attachDomain",
		Short: "attach domain to service",
		Long:  `This subcommand attach domain to a service`,
		Run:   srvAttachDomain}

	srvDetachDomainCmd = &cobra.Command{
		Use:   "srv:detachDomain",
		Short: "detach domain of the service",
		Long:  `This subcommand detach domain of the service`,
		Run:   srvDetachDomain}
)

func srvList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvList(client.Context(), req)
	uiCheckErr("Could not List the Services: %v", err)
	uiList(res)
}

func srvInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvInfo(client.Context(), req)
	uiCheckErr("Could not Get Service: %v", err)
	uiServicStatus(res)
}

func srvPortforward(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func srvStart(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvStart(client.Context(), req)
	uiCheckErr("Could not Start the Service: %v", err)
	uiServicStatus(res)
}

func srvStop(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvStop(client.Context(), req)
	uiCheckErr("Could not Stop the Service: %v", err)
	uiServicStatus(res)
}

func srvDestroy(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V1().SrvDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Service: %v", err)
	log.Printf("service %s deleted", req.Name)
}

func srvCreate(cmd *cobra.Command, args []string) {
	req := new(uvApi.SrvCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.ProductName = cmd.Flag("product").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()
	req.Variable = arrayFlagToMap(flagVariableArray)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvCreate(client.Context(), req)
	uiCheckErr("Could not Create the Service: %v", err)
	uiServicStatus(res)
}

func srvChangePlane(cmd *cobra.Command, args []string) {
	req := new(uvApi.ChangePlanReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan: %v", err)
	uiServicStatus(res)
}

func srvAttachDomain(cmd *cobra.Command, args []string) {
	req := new(uvApi.SrvDomainAttachReq)
	req.AttachIdentity = new(uvApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("name").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("attachment").Value.String()
	req.Endpoint = cmd.Flag("endpoint").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvAttachDomain(client.Context(), req)
	uiCheckErr("Could not Attach the Domain for Service: %v", err)
	uiServicStatus(res)
}

func srvDetachDomain(cmd *cobra.Command, args []string) {
	req := new(uvApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().SrvDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Service: %v", err)
	uiServicStatus(res)
}

func init() {
	// srv List:
	srvListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")

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
	srvAttachDomainCmd.MarkFlagRequired("name")
	srvAttachDomainCmd.MarkFlagRequired("endpoint")
	srvAttachDomainCmd.MarkFlagRequired("attachment")

	// srv Detach Domain:
	srvDetachDomainCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the service.")
	srvDetachDomainCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	srvDetachDomainCmd.MarkFlagRequired("name")
	srvDetachDomainCmd.MarkFlagRequired("attachment")

	// add service subcommands to root
	rootCmd.AddCommand(
		srvListCmd,
		srvInfoCmd,
		srvCreateCmd,
		srvChangePlaneCmd,
		srvPortforwardCmd,
		srvStartCmd,
		srvStopCmd,
		srvDestroyCmd,
		srvAttachDomainCmd,
		srvDetachDomainCmd)
}
