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
	res, err := client.V2().VolumeSpecList(client.Context(), req)
	uiCheckErr("Could not List the Volumes Spec: %v", err)
	uiList(res)
}

func volumeSpecInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeSpecInfo(client.Context(), req)
	uiCheckErr("Could not get the Volumes Spec: %v", err)
	uiVolumeSpec(res)
}

func volumeList(cmd *cobra.Command, args []string) {
	req := reqIndexForApp(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().VolumeList(client.Context(), req)
	uiCheckErr("Could not List the volume: %v", err)
	uiList(res)
}

func volumeInfo(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
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
	uiVolumeStatus(res)
}

func volumeDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().VolumeDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Volume: %v", err)
	log.Println("Task is done.")
}

func init() {
	// volume Spec list:
	volumeSpecListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number list")
	volumeSpecListCmd.Flags().StringVarP(&flagAppName, "app", "n", "", "page number list")

	// volume Spec info:
	volumeSpecInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeSpecInfoCmd.MarkFlagRequired("name")

	// volume list:
	volumeListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number list")
	volumeListCmd.Flags().StringVarP(&flagAppName, "app", "n", "", "page number list")

	// volume info:
	volumeInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeInfoCmd.MarkFlagRequired("name")

	// volume create:
	volumeCreateCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the volume.")
	volumeCreateCmd.Flags().StringP("volume-type", "v", "", "the type of volume")
	volumeCreateCmd.MarkFlagRequired("name")
	volumeCreateCmd.MarkFlagRequired("volume-type")

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
