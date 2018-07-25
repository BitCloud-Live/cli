package cmd

import (
	"github.com/spf13/cobra"
)

var (
	imgListCmd = &cobra.Command{
		Use:   "img:list",
		Short: "show all images",
		Long:  `This subcommand can pageing the images name.`,
		Run:   imgList}

	imgInfoCmd = &cobra.Command{
		Use:   "img:info",
		Short: "detail of Image",
		Long:  `This subcommand show the information of a images.`,
		Run:   imgInfo}

	imgImportCmd = &cobra.Command{
		Use:   "img:import",
		Short: "import image file",
		Long:  `This subcommand import image from exported docker format.`,
		Run:   imgImport}

	imgBuildCmd = &cobra.Command{
		Use:   "img:build",
		Short: "build an image from a Dockerfile",
		Long:  `This subcommand Build an image from a dockerfile.`,
		Run:   imgBuild}

	imgDeleteCmd = &cobra.Command{
		Use:   "img:delete",
		Short: "delete the image",
		Long:  `This subcommand delete the Image.`,
		Run:   imgDelete}
)

func imgList(cmd *cobra.Command, args []string) {
	req := reqIndex(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V1().ImgList(client.Context(), req)
	uiCheckErr("Could not List the Products: %v", err)
	uiList(res)
}

// TODO add to UV.proto
func imgInfo(cmd *cobra.Command, args []string) {

	//	req := reqIdentity(cmd)
	//	res, err := grpcClinet.UV.ImgInfo(grpcClinet.Ctx, req)
	//	uiCheckErr("Could not Get the Image Info: %v", err)
	//	uiImageInfo(res)
}

func imgImport(cmd *cobra.Command, args []string) {} //TODO

func imgBuild(cmd *cobra.Command, args []string) {} //TODO

func imgDelete(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V1().ImgDelete(client.Context(), req)
	uiCheckErr("Could not Destroy the Image: %v", err)
	log.Println("Work is done!")
}

func init() {
	// imgage List:
	imgListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")

	// imgage Info:
	imgInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the image")
	imgInfoCmd.MarkFlagRequired("name")

	// imgage Delete:
	imgDeleteCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the image")
	imgDeleteCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(
		imgListCmd,
		imgInfoCmd,
		// imgImportCmd,
		// imgBuildCmd,
		imgDeleteCmd)
}
