package cmd

import (
	"github.com/spf13/cobra"
)

var (
	prdListCmd = &cobra.Command{
		Use:   "prd:list",
		Short: "List of all products",
		Long:  `This subcommand can pageing the product name.`,
		Run:   prdList}

	prdInfoCmd = &cobra.Command{
		Use:   "prd:info",
		Short: "Detail of product",
		Long:  `This subcommand show the information, praice and ... of a product.`,
		Run:   prdInfo}
)

func prdList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().PrdList(client.Context(), req)
	uiCheckErr("Could not List the Products: %v", err)
	uiList(res)
}

func prdInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().PrdInfo(client.Context(), req)
	uiCheckErr("Could not Get the Product Info: %v", err)
	uiProduct(res)
}

func init() {
	// product list:
	prdListCmd.Flags().Int32("index", 0, "page of list")

	// product info:
	prdInfoCmd.Flags().StringP("name", "n", "", "name of Product")
	prdInfoCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(
		prdListCmd,
		prdInfoCmd)
}
