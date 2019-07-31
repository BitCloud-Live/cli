package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagTLS bool
)

func domainList(cmd *cobra.Command, args []string) {
	req := getCliRequestIndexForApp(args, 0, flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().DomainList(client.Context(), req)
	uiCheckErr("Could not List the domain: %v", err)
	uiList(res)
}

func domainInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().DomainInfo(client.Context(), req)
	uiCheckErr("Could not get the Domains: %v", err)
	uiDomainStatus(res)
}

// DomainCreate create domain
func DomainCreate(domain string, tls bool) (*ybApi.DomainStatusRes, error) {
	req := new(ybApi.DomainCreateReq)
	req.Domain = domain
	req.Tls = tls

	client := grpcConnect()
	defer client.Close()
	return client.V2().DomainCreate(client.Context(), req)
}
func domainCreate(cmd *cobra.Command, args []string) {
	res, err := DomainCreate(
		getCliRequiredArg(args, 0),
		flagTLS)

	uiCheckErr("Could not Create the Domain: %v", err)
	uiDomainStatus(res)
}

func domainDelete(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().DomainDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Domain: %v", err)
	log.Println("Task is done.")
}
