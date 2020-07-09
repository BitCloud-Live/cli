package cmd

import (
	"log"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

func srvList(cmd *cobra.Command, args []string) {
	req := getCliRequestIndexForApp(args, 0, flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SrvList(client.Context(), req)
	uiCheckErr("Could not List the Services", err)
	uiList(res)
}

// ServiceInfo get information of Service
func ServiceInfo(name string) (*ybApi.SrvStatusRes, error) {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	return client.V2().SrvInfo(client.Context(), req)
}

func srvInfo(cmd *cobra.Command, args []string) {
	res, err := ServiceInfo(
		getCliRequiredArg(args, 0))
	uiCheckErr("Could not Get Service", err)
	uiServicStatus(res)
}

func srvPortforward(cmd *cobra.Command, args []string) {
	srvName := getCliRequiredArg(args, 0)
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	info, err := client.V2().SrvInfo(client.Context(), req)
	res, err := client.V2().SrvPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service", err)
	uiServicStatus(info)
	uiPortforward(srvName, nil, res)
}

// ServiceStart start service by name
func ServiceStart(name string) (*ybApi.SrvStatusRes, error) {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	return client.V2().SrvStart(client.Context(), req)
}
func srvStart(cmd *cobra.Command, args []string) {
	res, err := ServiceStart(
		getCliRequiredArg(args, 0))
	uiCheckErr("Could not Start the Service", err)
	uiServicStatus(res)
}

// ServiceStop stop service by name
func ServiceStop(name string) (*ybApi.SrvStatusRes, error) {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	return client.V2().SrvStop(client.Context(), req)
}
func srvStop(cmd *cobra.Command, args []string) {
	res, err := ServiceStop(
		getCliRequiredArg(args, 0))
	uiCheckErr("Could not Stop the Service", err)
	uiServicStatus(res)
}

func srvDestroy(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().SrvDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Service", err)
	log.Printf("service %s deleted", req.Name)
}

// ServiceCreate create Service by ProductName and PlanName
func ServiceCreate(productName, serviceName, plan string, values map[string]string) (*ybApi.SrvStatusRes, error) {
	req := new(ybApi.SrvCreateReq)
	req.Name = serviceName
	req.ProductName = productName
	req.Plan = plan
	req.Values = values

	client := grpcConnect()
	defer client.Close()
	return client.V2().SrvCreate(client.Context(), req)
}
func srvCreate(cmd *cobra.Command, args []string) {
	res, err := ServiceCreate(
		getCliRequiredArg(args, 0),      // ProductName
		cmd.Flag("name").Value.String(), // ServiceName
		cmd.Flag("plan").Value.String(), // Plan
		arrayFlagToMap(flagVariableArray))

	uiCheckErr("Could not Create the Service", err)
	uiServicStatus(res)
}

// ServiceLinkDomain like domain to service
func ServiceLinkDomain(serviceName, domainName, endpoint, path string) (*ybApi.SrvStatusRes, error) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = serviceName
	req.AttachIdentity.Attachment = domainName
	req.Endpoint = endpoint
	req.Path = path

	client := grpcConnect()
	defer client.Close()
	return client.V2().SrvAttachDomain(client.Context(), req)
}
func srvAttachDomain(cmd *cobra.Command, args []string) {
	res, err := ServiceLinkDomain(
		cmd.Flag("service").Value.String(),
		cmd.Flag("domain").Value.String(),
		cmd.Flag("endpoint").Value.String(),
		cmd.Flag("path").Value.String())

	uiCheckErr("Could not Attach the Domain for Service", err)
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
	uiCheckErr("Could not Detach the Domain for Service", err)
	uiServicStatus(res)
}
