package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagVarPort         uint64
	flagVarMinScale     uint64
	minScale            uint64
	flagVarImage        string
	flagVarEndpointType string
	flagVarDebug        bool
)

func appList(cmd *cobra.Command, args []string) {
	req := getRequestIndex(flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppList(client.Context(), req)
	uiCheckErr("Could not List the Applications", err)
	uiList(res)
}

func appInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application", err)
	uiApplicationStatus(res)
}

func appOpen(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application", err)
	uiApplicationOpen(res)
}

func appLog(cmd *cobra.Command, args []string) {
	streamAppLog(args)
}

func appRun(cmd *cobra.Command, args []string) {
	req := new(ybApi.ShellReq)
	req.Name = getCliRequiredArg(args, 0)
	if len(args) < 2 {
		log.Fatal("no command run")
	}
	req.Command = args[1:]
	client := grpcConnect()
	defer client.Close()
	output, err := client.V2().AppShell(client.Context(), req)
	uiCheckErr("Could not Get Application remote command output", err)
	s := string(output.Chunk)
	log.Printf("output: %s", s)
}

// AppStart start the application by name
func AppStart(name string) (*ybApi.AppStatusRes, error) {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	return client.V2().AppStart(client.Context(), req)
}
func appStart(cmd *cobra.Command, args []string) {
	appName := getCliRequiredArg(args, 0)
	res, err := AppStart(appName)
	uiCheckErr("Could not Start the Application", err)
	uiApplicationStatus(res)
}

// AppStop stop the application by name
func AppStop(name string) (*ybApi.AppStatusRes, error) {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	return client.V2().AppStop(client.Context(), req)
}

func appStop(cmd *cobra.Command, args []string) {
	appName := getCliRequiredArg(args, 0)
	res, err := AppStop(appName)
	uiCheckErr("Could not Stop the Application", err)
	uiApplicationStatus(res)
}

func appDestroy(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().AppDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Application", err)
	log.Printf("app %s deleted", req.Name)
}

// ApplicationCreate create application by link of docker image
func ApplicationCreate(appName, image, plan, EndpointType string, port, minScale uint64, debugMode bool) (*ybApi.AppStatusRes, error) {
	if err := endpointTypeValid(EndpointType); err != nil {
		log.Panic(err)
	}
	req := new(ybApi.AppCreateReq)
	req.Name = appName
	req.Plan = plan
	req.Values = make(map[string]string)
	req.Values["ports"] = fmt.Sprintf("%d/%s", port, EndpointType)
	req.Values["minimum-scale"] = fmt.Sprintf("%d", minScale)
	req.Values["image"] = image
	req.Values["debug"] = fmt.Sprintf("%b", debugMode)

	client := grpcConnect()
	defer client.Close()
	return client.V2().AppCreate(client.Context(), req)
}

func appCreate(cmd *cobra.Command, args []string) {
	res, err := ApplicationCreate(
		cmd.Flag("name").Value.String(),
		cmd.Flag("image").Value.String(),
		cmd.Flag("plan").Value.String(),
		flagVarEndpointType,
		flagVarPort,
		flagVarMinScale,
		flagVarDebug)
	uiCheckErr("Could not Create the Application", err)
	uiApplicationStatus(res)
}

/*
func appChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = getCliRequiredArg(args, 0)
	req.Plan = getCliRequiredArg(args, 1)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan", err)
	uiApplicationStatus(res)
}*/

func appUpdate(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvConfigSetReq)
	req.Name = getCliRequiredArg(args, 0)
	req.Values = make(map[string]string)
	if flagVarMinScale != 0 {
		req.Values["minimum-scale"] = fmt.Sprintf("%d", flagVarMinScale)
	}
	if flagVarPort != 0 {
		req.Values["ports"] = fmt.Sprintf("%d/%s", flagVarPort, flagVarEndpointType)
	}
	if flagVarImage != "" {
		req.Values["image"] = flagVarImage
	}
	client := grpcConnect()
	defer client.Close()

	res, err := client.V2().AppConfigSet(client.Context(), req)
	uiCheckErr("Could not Set the Config for Application", err)
	uiApplicationStatus(res)
}

// ApplicationAddEnvironmentVariable add Environment Variable to Application
func ApplicationAddEnvironmentVariable(serviceName string, variable map[string]string) (*ybApi.AppStatusRes, error) {
	req := new(ybApi.AppAddEnvironmentVariableReq)
	req.Name = serviceName
	req.Variables = variable

	client := grpcConnect()
	defer client.Close()
	return client.V2().AppAddEnvironmentVariable(client.Context(), req)
}

func appAddEnvironmentVariable(cmd *cobra.Command, args []string) {
	res, err := ApplicationAddEnvironmentVariable(
		getCliRequiredArg(args, 0),
		arrayFlagToMap(flagVariableArray))

	uiCheckErr("Could not Add the Environment Variable for Application", err)
	uiApplicationStatus(res)
}

func appRemoveEnvironmentVariable(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = getCliRequiredArg(args, 0)
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppRemoveEnvironmentVariable(client.Context(), req)
	uiCheckErr("Could not Remove the Environment Variable for Application", err)
	uiApplicationStatus(res)
}

func appReset(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppReset(client.Context(), req)
	uiCheckErr("Could not Reset the Application", err)
	uiApplicationStatus(res)
}

// ApplicationLinkService Application <-> Service
func ApplicationLinkService(applicationName, serviceName string) (*ybApi.AppStatusRes, error) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = applicationName
	req.Service = serviceName
	client := grpcConnect()
	defer client.Close()
	return client.V2().AppSrvBind(client.Context(), req)
}

func appSrvBind(cmd *cobra.Command, args []string) {
	res, err := ApplicationLinkService(
		cmd.Flag("application").Value.String(),
		cmd.Flag("service").Value.String())
	uiCheckErr("Could not Bind the Service for Application", err)
	uiApplicationStatus(res)
}

func appSrvUnBind(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = cmd.Flag("application").Value.String()
	req.Service = cmd.Flag("service").Value.String()
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppSrvUnBind(client.Context(), req)
	uiCheckErr("Could not Unbind the Service for Application", err)
	uiApplicationStatus(res)
}

// AppAttachVolume link the application and volume
func AppAttachVolume(appName, volumeName, path string) (*ybApi.AppStatusRes, error) {
	req := new(ybApi.VolumeMount)
	req.Name = appName
	req.Attachment = volumeName
	req.MountPath = path

	client := grpcConnect()
	defer client.Close()
	return client.V2().AppAttachVolume(client.Context(), req)
}

func appAttachVolume(cmd *cobra.Command, args []string) {
	res, err := AppAttachVolume(
		cmd.Flag("application").Value.String(),
		cmd.Flag("volume").Value.String(),
		cmd.Flag("path").Value.String())
	uiCheckErr("Could not Attach the Volume for Application", err)
	uiApplicationStatus(res)
}

func appDetachVolume(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("application").Value.String()
	req.Attachment = cmd.Flag("volume").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachVolume(client.Context(), req)
	uiCheckErr("Could not Detach the Volume for Application", err)
	uiApplicationStatus(res)
}

// AppAttachDomain link the application and domain
func AppAttachDomain(appName, domainName, path, endpoint string) (*ybApi.AppStatusRes, error) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = appName
	req.AttachIdentity.Attachment = domainName
	req.Path = path
	req.Endpoint = endpoint

	client := grpcConnect()
	defer client.Close()
	return client.V2().AppAttachDomain(client.Context(), req)
}

func appAttachDomain(cmd *cobra.Command, args []string) {
	var (
		app  = forceFlagGetStrValue(cmd, "application", "Enter Application Name:")
		dom  = forceFlagGetStrValue(cmd, "domain", "Enter Domain:")
		path = forceFlagGetStrValue(cmd, "path", "Enter Path:")
		ep   = forceFlagGetStrValue(cmd, "endpoint", "Enter Endpoint [ format: 8080/http ]:")
	)

	res, err := AppAttachDomain(app, dom, path, ep)

	uiCheckErr("Could not Attach the Domain for Application", err)
	uiApplicationStatus(res)
}

func appDetachDomain(cmd *cobra.Command, args []string) {
	var (
		app  = forceFlagGetStrValue(cmd, "application", "Enter Application Name:")
		dom  = forceFlagGetStrValue(cmd, "domain", "Enter Domain:")
		path = forceFlagGetStrValue(cmd, "path", "Enter Path:")
		ep   = forceFlagGetStrValue(cmd, "endpoint", "Enter Endpoint [ format: 8080/http ]:")
	)

	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = app
	req.AttachIdentity.Attachment = dom
	req.Path = path
	req.Endpoint = ep

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Application", err)
	uiApplicationStatus(res)
}
