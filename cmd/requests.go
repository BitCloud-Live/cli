package cmd

import (
	"log"
	ybApi "github.com/yottab/proto-api/proto"
)

func getRequestIndex(pageIndex int32) *ybApi.ListReq {
	req := new(ybApi.ListReq)
	req.Index = pageIndex //cli get value by flagIndex
	return req
}

func getCliRequestIndexForApp(cliArgs []string, index int, pageIndex int32) *ybApi.AppListReq {
	appName := getCliArg(cliArgs, index, "")         // not Required have application.Name
	return getRequestIndexForApp(appName, pageIndex) //cli get pageIndex value by flagIndex
}
func getRequestIndexForApp(appName string, pageIndex int32) *ybApi.AppListReq {
	req := new(ybApi.AppListReq)
	req.Index = pageIndex //cli get value by flagIndex
	req.App = appName     // not Required have application.Name
	return req
}

func getCliRequestIdentity(cliArgs []string, index int) *ybApi.Identity {
	val := getCliRequiredArg(cliArgs, index)
	return getRequestIdentity(val)
}

func getRequestIdentity(name string) *ybApi.Identity {
	ri := checkRequiredArg(name)
	req := new(ybApi.Identity)
	req.Name = ri
	return req
}

// get input Arg and check it

func getCliArg(cliArgs []string, index int, defaultValue string) string {
	if len(cliArgs) < index+1 {
		return defaultValue
	}
	return cliArgs[index]
}
func getCliRequiredArg(cliArgs []string, index int) string {
	if len(cliArgs) < index+1 {
		log.Fatalf("not enugh required arg")
	}
	return cliArgs[index]
}
func checkArg(arg string, defaultValue string) string {
	if arg == "" {
		return defaultValue
	}
	return arg
}
func checkRequiredArg(arg string) string {
	if arg == "" {
		log.Fatalf("not enugh required arg")
	}
	return arg
}
