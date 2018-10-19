package cmd

import (
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
	uvApi "github.com/uvcloud/uv-api-go/proto"
)

func reqIndex(cmd *cobra.Command) *uvApi.ListReq {
	req := new(uvApi.ListReq)
	req.Index = flagIndex
	return req
}

func reqIndexForApp(cmd *cobra.Command) *uvApi.AppListReq {
	req := new(uvApi.AppListReq)
	req.Index = flagIndex
	req.App = flagAppName
	return req
}

func reqIdentity(cmd *cobra.Command) *uvApi.Identity {
	req := new(uvApi.Identity)
	req.Name = cmd.Flag("name").Value.String()
	return req
}
