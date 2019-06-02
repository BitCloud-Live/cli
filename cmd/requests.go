package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

const (
	// RequiredArg : this arg must be set
	RequiredArg = true

	// NotRequiredArg : this arg Not required to set
	NotRequiredArg = false
)

func reqIndex(cmd *cobra.Command) *ybApi.ListReq {
	req := new(ybApi.ListReq)
	req.Index = flagIndex
	return req
}

func reqIndexForApp(args []string, index int, required bool) *ybApi.AppListReq {
	req := new(ybApi.AppListReq)
	req.Index = flagIndex
	req.App = argValue(args, index, required, "NO_NAME")
	return req
}

func reqIdentity(args []string, index int, required bool) *ybApi.Identity {
	val := argValue(args, index, required, "NO_NAME")
	req := new(ybApi.Identity)
	req.Name = val
	return req
}

func argValue(args []string, index int, required bool, defaultValue string) string {
	if len(args) < index+1 {
		if required {
			log.Fatalf("not enugh required arg")
		} else {
			return defaultValue
		}
	}
	return args[index]
}
