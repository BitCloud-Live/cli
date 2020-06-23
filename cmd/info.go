package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

func info(cmd *cobra.Command, args []string) {
	client := grpcConnect()
	defer client.Close()
	req := &ybApi.Empty{}
	res, err := client.V2().AccountInfo(client.Context(), req)
	uiCheckErr("Could not Get Application", err)
	uiAccountInfo(res)
}
