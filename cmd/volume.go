package cmd

import (
	"log"

	"github.com/spf13/cobra"
	uvApi "github.com/uvcloud/uv-api-go/proto"
)

var (
	volumeSpecListCmd = &cobra.Command{
		Use:   "vol:specList",
		Short: "",
		Long:  ``,
		Run:   volumeSpecList}

	volumeSpecInfoCmd = &cobra.Command{
		Use:   "vol:spec",
		Short: "",
		Long:  ``,
		Run:   volumeSpecInfo}

	volumeListCmd = &cobra.Command{
		Use:   "vol:list",
		Short: "",
		Long:  ``,
		Run:   volumeList}

	volumeInfoCmd = &cobra.Command{
		Use:   "vol:info",
		Short: "",
		Long:  ``,
		Run:   volumeInfo}

	volumeCreateCmd = &cobra.Command{
		Use:   "vol:create",
		Short: "",
		Long:  ``,
		Run:   volumeCreate}

	volumeDeleteCmd = &cobra.Command{
		Use:   "vol:delete",
		Short: "",
		Long:  ``,
		Run:   volumeDelete}
)

func volumeSpecList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().VolumeSpecList(client.Context(), req)
	uiCheckErr("Could not List the Volumes Spec: %v", err)
	uiList(res)
}

func volumeSpecInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().VolumeSpecInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes Spec: %v", err)
	uiVolumeSpec(res)
}

func volumeList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().VolumeList(client.Context(), req)
	uiCheckErr("Could not List the volume: %v", err)
	uiList(res)
}

func volumeInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().VolumeInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes: %v", err)
	//TODO
	log.Println(res)
}

func volumeCreate(cmd *cobra.Command, args []string) {
	req := new(uvApi.VolumeCreateReq)
	req.Name = cmd.Flag("name").Value.String()
	req.Spec = cmd.Flag("spec").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().VolumeCreate(client.Context(), req)
	uiCheckErr("Could not Create the Volume: %v", err)
	//TODO
	log.Println(res)
}

func volumeDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V1().VolumeDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Volume: %v", err)
	log.Println("Task is done.")
}

func init() {
	// volume Spec list:
	volumeSpecListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")

	// volume Spec info:
	volumeSpecInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeSpecInfoCmd.MarkFlagRequired("name")

	// volume list:
	volumeListCmd.Flags().Int32("index", 0, "page number list")

	// volume info:
	volumeInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeInfoCmd.MarkFlagRequired("name")

	// volume create:
	volumeCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeCreateCmd.Flags().StringP("spec", "s", "", "the name of volume's spac.")
	volumeCreateCmd.MarkFlagRequired("name")
	volumeCreateCmd.MarkFlagRequired("spec")

	// volume delete:
	volumeDeleteCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeDeleteCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(
		volumeSpecListCmd,
		volumeSpecInfoCmd,
		volumeListCmd,
		volumeInfoCmd,
		volumeCreateCmd,
		volumeDeleteCmd)
}
