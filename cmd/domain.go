package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagTLS bool
)

func domainList(cmd *cobra.Command, args []string) {
	req := reqIndexForApp(args, 0, NotRequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().DomainList(client.Context(), req)
	uiCheckErr("Could not List the domain: %v", err)
	uiList(res)
}

func domainInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().DomainInfo(client.Context(), req)
	uiCheckErr("Could not get the Domains: %v", err)
	uiDomainStatus(res)
}

func domainCreate(cmd *cobra.Command, args []string) {
	req := new(ybApi.DomainCreateReq)
	req.Domain = argValue(args, 0, RequiredArg, "")
	req.Tls = flagTLS

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().DomainCreate(client.Context(), req)
	uiCheckErr("Could not Create the Domain: %v", err)
	uiDomainStatus(res)
}

func domainDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().DomainDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Domain: %v", err)
	log.Println("Task is done.")
}
