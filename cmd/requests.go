package cmd

import (
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	ybApi "github.com/yottab/proto-api/proto"
)

func reqIndex(cmd *cobra.Command) *ybApi.ListReq {
	req := new(ybApi.ListReq)
	req.Index = flagIndex
	return req
}

func reqIndexForApp(cmd *cobra.Command) *ybApi.AppListReq {
	req := new(ybApi.AppListReq)
	req.Index = flagIndex
	req.App = flagAppName
	return req
}

func reqIdentity(cmd *cobra.Command) *ybApi.Identity {
	req := new(ybApi.Identity)
	req.Name = cmd.Flag("name").Value.String()
	return req
}
