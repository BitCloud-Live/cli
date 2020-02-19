package cmd

import (
	"context"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagVarPort         uint64
	flagVarMinScale     uint64
	minScale            uint64
	flagVarImage        string
	flagVarEndpointType string
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

	client := grpcConnect()
	defer client.Close()
	req := getCliRequestIdentity(args, 0)
	logClient, err := client.V2().AppLog(context.Background(), req)
	uiCheckErr("Could not Get Application log", err)
	uiApplicationLog(logClient)
	//TODO: Api testing
	// req := new(ybApi.TailRequest)
	// req.Name = "app-name"
	// req.Tail = 1000
	// output, err := client.V2().AppLogTail(context.Background(), req)
	// log.Debug(err)
	// s := string(output.Chunk)
	// log.Printf("output: %s", s)
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
func ApplicationCreate(appName, image, plan, EndpointType string, port, minScale uint64) (*ybApi.AppStatusRes, error) {
	if err := endpointTypeValid(EndpointType); err != nil {
		log.Panic(err)
	}
	req := new(ybApi.AppCreateReq)
	req.Name = appName
	req.Plan = plan
	req.Config = new(ybApi.AppConfig)
	req.Config.Port = port
	req.Config.EndpointType = EndpointType
	req.Config.MinScale = minScale
	req.Config.Image = image

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
		flagVarMinScale)
	uiCheckErr("Could not Create the Application", err)
	uiApplicationStatus(res)
}

func appChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = getCliRequiredArg(args, 0)
	req.Plan = getCliRequiredArg(args, 1)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan", err)
	uiApplicationStatus(res)
}

func appUpdate(cmd *cobra.Command, args []string) {
	req := new(ybApi.ConfigSetReq)
	req.Name = getCliRequiredArg(args, 0)
	req.Config = new(ybApi.AppConfig)
	req.Config.MinScale = flagVarMinScale
	req.Config.Port = flagVarPort
	req.Config.Image = flagVarImage
	req.Config.Routes = flagVariableArray
	req.Config.EndpointType = flagVarEndpointType
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
	req := new(ybApi.VolumeAttachReq)
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
		app  = cmd.Flag("application").Value.String()
		dom  = cmd.Flag("domain").Value.String()
		path = cmd.Flag("path").Value.String()
		ep   = cmd.Flag("endpoint").Value.String()
	)

	if app == "" {
		app = readFromConsole("Enter Application Name:")
	}
	if dom == "" {
		dom = readFromConsole("Enter Domain:")
	}
	if path == "" {
		dom = readFromConsole("Enter Path:")
	}

	res, err := AppAttachDomain(app, dom, path, ep)

	uiCheckErr("Could not Attach the Domain for Application", err)
	uiApplicationStatus(res)
}

func appDetachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("application").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("domain").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachDomain(client.Context(), req)
	uiCheckErr("Could not Detach the Domain for Application", err)
	uiApplicationStatus(res)
}
