package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
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
	req := reqIndexForApp(args, 0, NotRequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvList(client.Context(), req)
	uiCheckErr("Could not List the Services: %v", err)
	uiList(res)
}

func srvInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvInfo(client.Context(), req)
	uiCheckErr("Could not Get Service: %v", err)
	uiServicStatus(res)
}

func srvPortforward(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func srvStart(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvStart(client.Context(), req)
	uiCheckErr("Could not Start the Service: %v", err)
	uiServicStatus(res)
}

func srvStop(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvStop(client.Context(), req)
	uiCheckErr("Could not Stop the Service: %v", err)
	uiServicStatus(res)
}

func srvDestroy(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().SrvDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Service: %v", err)
	log.Printf("service %s deleted", req.Name)
}

func srvCreate(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.ProductName = argValue(args, 0, RequiredArg, "")
	req.Plan = cmd.Flag("plan").Value.String()
	req.Variable = arrayFlagToMap(flagVariableArray)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvCreate(client.Context(), req)
	uiCheckErr("Could not Create the Service: %v", err)
	uiServicStatus(res)
}

///////////////////////////////////////////////////////

func srvConfigSet(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvConfigSetReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Variables = arrayFlagToMap(flagVariableArray)
	client := grpcConnect()
	defer client.Close()

	res, err := client.V2().SrvConfigSet(client.Context(), req)
	uiCheckErr("Could not Set the Config for Service: %v", err)
	uiServicStatus(res)
}

func srvConfigUnset(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvConfigUnset(client.Context(), req)
	uiCheckErr("Could not Unset the Config for Service: %v", err)
	uiServicStatus(res)
}

func srvChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Plan = argValue(args, 1, RequiredArg, "")

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan: %v", err)
	uiServicStatus(res)
}

func srvAttachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("service").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("domain").Value.String()
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
	req.AttachIdentity.Name = cmd.Flag("service").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("domain").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Service: %v", err)
	uiServicStatus(res)
}
