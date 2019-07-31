package cmd

import (
	"github.com/spf13/cobra"
)

func prdList(cmd *cobra.Command, args []string) {
	req := getRequestIndex(flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().PrdList(client.Context(), req)
	uiCheckErr("Could not List the Products: %v", err)
	uiList(res)
}

func prdInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().PrdInfo(client.Context(), req)
	uiCheckErr("Could not Get the Product Info: %v", err)
	uiProduct(res)
}
