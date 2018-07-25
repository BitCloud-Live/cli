package cmd

import (
	"github.com/spf13/cobra"
	uvApi "github.com/uvcloud/uv-api-go/proto"
)

var (
	flagTLS bool
)

var (
	domainListCmd = &cobra.Command{
		Use:   "dom:list",
		Short: "",
		Long:  ``,
		Run:   domainList}

	domainInfoCmd = &cobra.Command{
		Use:   "dom:info",
		Short: "",
		Long:  ``,
		Run:   domainInfo}

	domainCreateCmd = &cobra.Command{
		Use:   "dom:create",
		Short: "",
		Long:  ``,
		Run:   domainCreate}

	domainDeleteCmd = &cobra.Command{
		Use:   "dom:delete",
		Short: "",
		Long:  ``,
		Run:   domainDelete}
)

func domainList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().DomainList(client.Context(), req)
	uiCheckErr("Could not List the domain: %v", err)
	uiList(res)
}

func domainInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().DomainInfo(client.Context(), req)
	uiCheckErr("Could not get the Domains: %v", err)
	uiDomainStatus(res)
}

func domainCreate(cmd *cobra.Command, args []string) {
	req := new(uvApi.DomainCreateReq)
	req.Domain = cmd.Flag("domain").Value.String()
	req.Tls = flagTLS

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().DomainCreate(client.Context(), req)
	uiCheckErr("Could not Create the Domain: %v", err)
	uiDomainStatus(res)
}

func domainDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V1().DomainDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Domain: %v", err)
	log.Println("Task is done.")
}

func init() {
	// domain list:
	domainListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")

	// domain info:
	domainInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the domain.")
	domainInfoCmd.MarkFlagRequired("name")

	// domain create:
	domainCreateCmd.Flags().StringP("domain", "d", "", "the name of domain's spac.")
	domainCreateCmd.Flags().BoolVar(&flagTLS, "TLS", false, "enable TLS for domain")
	domainCreateCmd.MarkFlagRequired("domain")

	// domain delete:
	domainDeleteCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the domain.")
	domainDeleteCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(
		domainListCmd,
		domainInfoCmd,
		domainCreateCmd,
		domainDeleteCmd)
}
