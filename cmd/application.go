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
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppList(client.Context(), req)
	uiCheckErr("Could not List the Applications: %v", err)
	uiList(res)
}

func appInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiApplicationStatus(res)
}

func appOpen(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiApplicationOpen(res)
}

func appLog(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	logClient, err := client.V2().AppLog(context.Background(), req)
	if err != nil {
		panic(err)
	}
	uiCheckErr("Could not Get Application log: %v", err)
	uiApplicationLog(logClient)
}

func appFTPMount(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppFTPPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
	uiNFSMount(res)
}

func appStart(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppStart(client.Context(), req)
	uiCheckErr("Could not Start the Application: %v", err)
	uiApplicationStatus(res)
}

func appStop(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppStop(client.Context(), req)
	uiCheckErr("Could not Stop the Application: %v", err)
	uiApplicationStatus(res)
}

func appDestroy(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().AppDestroy(client.Context(), req)
	uiCheckErr("Could not Destroy the Application: %v", err)
	log.Printf("app %s deleted", req.Name)
}

func appCreate(cmd *cobra.Command, args []string) {
	if err := endpointTypeValid(flagVarEndpointType); err != nil {
		log.Panic(err)
	}
	req := new(ybApi.AppCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Plan = cmd.Flag("plan").Value.String()
	req.Config = new(ybApi.AppConfig)
	req.Config.Port = flagVarPort
	req.Config.EndpointType = flagVarEndpointType
	req.Config.MinScale = flagVarMinScale
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppCreate(client.Context(), req)
	uiCheckErr("Could not Create the Application: %v", err)
	uiApplicationStatus(res)
}

func appChangePlane(cmd *cobra.Command, args []string) {
	req := new(ybApi.ChangePlanReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Plan = argValue(args, 1, RequiredArg, "")

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppChangePlan(client.Context(), req)
	uiCheckErr("Could not Change the Plan: %v", err)
	uiApplicationStatus(res)
}

func appUpdate(cmd *cobra.Command, args []string) {
	req := new(ybApi.ConfigSetReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Config = new(ybApi.AppConfig)
	req.Config.MinScale = flagVarMinScale
	req.Config.Port = flagVarPort
	req.Config.Image = flagVarImage
	req.Config.Routes = flagVariableArray
	req.Config.EndpointType = flagVarEndpointType
	client := grpcConnect()
	defer client.Close()

	res, err := client.V2().AppConfigSet(client.Context(), req)
	uiCheckErr("Could not Set the Config for Application: %v", err)
	uiApplicationStatus(res)
}

func appAddEnvironmentVariable(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppAddEnvironmentVariableReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Variables = arrayFlagToMap(flagVariableArray)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAddEnvironmentVariable(client.Context(), req)
	uiCheckErr("Could not Add the Environment Variable for Application: %v", err)
	uiApplicationStatus(res)
}

func appRemoveEnvironmentVariable(cmd *cobra.Command, args []string) {
	req := new(ybApi.UnsetReq)
	req.Name = argValue(args, 0, RequiredArg, "")
	req.Variables = flagVariableArray

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppRemoveEnvironmentVariable(client.Context(), req)
	uiCheckErr("Could not Remove the Environment Variable for Application: %v", err)
	uiApplicationStatus(res)
}

func appReset(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppReset(client.Context(), req)
	uiCheckErr("Could not Reset the Application: %v", err)
	uiApplicationStatus(res)
}

func appSrvBind(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = cmd.Flag("application").Value.String()
	req.Service = cmd.Flag("service").Value.String()
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppSrvBind(client.Context(), req)
	uiCheckErr("Could not Bind the Service for Application: %v", err)
	uiApplicationStatus(res)
}

func appSrvUnBind(cmd *cobra.Command, args []string) {
	req := new(ybApi.AppSrvBindReq)
	req.Name = cmd.Flag("application").Value.String()
	req.Service = cmd.Flag("service").Value.String()
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppSrvUnBind(client.Context(), req)
	uiCheckErr("Could not Unbind the Service for Application: %v", err)
	uiApplicationStatus(res)
}

func appAttachVolume(cmd *cobra.Command, args []string) {
	req := new(ybApi.VolumeAttachReq)
	req.Name = cmd.Flag("application").Value.String()
	req.Attachment = cmd.Flag("volume").Value.String()
	req.MountPath = cmd.Flag("path").Value.String()

	client := grpcConnect()

	defer client.Close()
	res, err := client.V2().AppAttachVolume(client.Context(), req)
	uiCheckErr("Could not Attach the Volume for Application: %v", err)
	uiApplicationStatus(res)
}

func appDetachVolume(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("application").Value.String()
	req.Attachment = cmd.Flag("volume").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppDetachVolume(client.Context(), req)
	uiCheckErr("Could not Detach the Volume for Application: %v", err)
	uiApplicationStatus(res)
}

func appAttachDomain(cmd *cobra.Command, args []string) {
	req := new(ybApi.SrvDomainAttachReq)
	req.AttachIdentity = new(ybApi.AttachIdentity)
	req.AttachIdentity.Name = cmd.Flag("application").Value.String()
	req.AttachIdentity.Attachment = cmd.Flag("domain").Value.String()
	req.Path = cmd.Flag("path").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAttachDomain(client.Context(), req)
	uiCheckErr("Could not Attach the Domain for Application: %v", err)
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
	uiCheckErr("Could not Detach the Domain for Application: %v", err)
	uiApplicationStatus(res)
}
