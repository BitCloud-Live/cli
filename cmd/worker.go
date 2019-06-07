package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

func workerList(cmd *cobra.Command, args []string) {
	// TODO get index of page
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppListWorkers(client.Context(), req)
	uiCheckErr("Could not List the Applications: %v", err)
	uiList(res)
}

func workerInfo(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = getCliRequiredArg(args, 0)
	req.Attachment = getCliRequiredArg(args, 1)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiWorkerStatus(res)
}

func workerPortforward(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = getCliRequiredArg(args, 0)
	req.Attachment = getCliRequiredArg(args, 1)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func workerDestroy(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = getCliRequiredArg(args, 0)
	req.Attachment = getCliRequiredArg(args, 1)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppRemoveWorker(client.Context(), req)
	uiCheckErr("Could not Destroy the Application: %v", err)
	log.Printf("worker %s deleted", req.Attachment)
	uiList(res)
}

func workerCreate(cmd *cobra.Command, args []string) {
	req := new(ybApi.WorkerReq)
	req.Identities = new(ybApi.AttachIdentity)
	req.Identities.Name = getCliRequiredArg(args, 0)
	req.Identities.Attachment = cmd.Flag("name").Value.String()
	req.Config = new(ybApi.WorkerConfig)
	req.Config.Port = flagVarPort
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAddWorker(client.Context(), req)
	uiCheckErr("Could not Create the Application: %v", err)
	uiList(res)
}

func workerUpdate(cmd *cobra.Command, args []string) {
	req := new(ybApi.WorkerReq)
	req.Identities = new(ybApi.AttachIdentity)
	req.Identities.Name = getCliRequiredArg(args, 0)
	req.Identities.Attachment = cmd.Flag("name").Value.String()
	req.Config = new(ybApi.WorkerConfig)
	req.Config.Port = flagVarPort
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerUpdate(client.Context(), req)
	uiCheckErr("Could not Update the Application: %v", err)
	uiList(res)
}
