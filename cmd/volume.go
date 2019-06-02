package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	volumeSpecListCmd = &cobra.Command{
		Use:   "vol:type-list",
		Short: "",
		Long:  ``,
		Run:   volumeSpecList}

	volumeSpecInfoCmd = &cobra.Command{
		Use:   "vol:type",
		Short: "",
		Long:  ``,
		Run:   volumeSpecInfo}
)

func volumeSpecList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeSpecList(client.Context(), req)
	uiCheckErr("Could not List the Volumes Spec: %v", err)
	uiList(res)
}

func volumeSpecInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeSpecInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes Spec: %v", err)
	uiVolumeSpec(res)
}

func volumeList(cmd *cobra.Command, args []string) {
	req := reqIndexForApp(args, 0, NotRequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeList(client.Context(), req)
	uiCheckErr("Could not List the volume: %v", err)
	uiList(res)
}

func volumeInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes: %v", err)
	uiVolumeStatus(res)
}

func volumeCreate(cmd *cobra.Command, args []string) {
	req := new(ybApi.VolumeCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Spec = cmd.Flag("volume-type").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeCreate(client.Context(), req)
	uiCheckErr("Could not Create the Volume: %v", err)
	//TODO
	log.Println(res)
}

func volumeDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(args, 0, RequiredArg)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().VolumeDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Volume: %v", err)
	log.Println("Task is done.")
}
