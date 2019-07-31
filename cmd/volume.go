package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

func volumeSpecList(cmd *cobra.Command, args []string) {
	req := getRequestIndex(flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeSpecList(client.Context(), req)
	uiCheckErr("Could not List the Volumes Spec: %v", err)
	uiList(res)
}

func volumeSpecInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeSpecInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes Spec: %v", err)
	uiVolumeSpec(res)
}

func volumeList(cmd *cobra.Command, args []string) {
	req := getCliRequestIndexForApp(args, 0, flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeList(client.Context(), req)
	uiCheckErr("Could not List the volume: %v", err)
	uiList(res)
}

func volumeInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes: %v", err)
	uiVolumeStatus(res)
}

// VolumeCreate crate a volume by name and type
func VolumeCreate(name, volumeType string) (*ybApi.VolumeStatusRes, error) {
	req := new(ybApi.VolumeCreateReq)
	req.Name = name
	req.Spec = volumeType

	client := grpcConnect()

	log.Println(client)

	defer client.Close()

	return client.V2().VolumeCreate(client.Context(), req)
}
func volumeCreate(cmd *cobra.Command, args []string) {
	res, err := VolumeCreate(
		cmd.Flag("name").Value.String(),
		cmd.Flag("volume-type").Value.String())

	uiCheckErr("Could not Create the Volume: %v", err)
	uiVolumeStatus(res)
}

// VolumeDelete delete a volume by name
func VolumeDelete(name string) error {
	req := getRequestIdentity(name)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().VolumeDelete(client.Context(), req)
	return err
}
func volumeDelete(cmd *cobra.Command, args []string) {
	name := getCliRequiredArg(args, 0)
	err := VolumeDelete(name)
	uiCheckErr("Could not Delete the Volume: %v", err)
	log.Println("Task is done.")
}
