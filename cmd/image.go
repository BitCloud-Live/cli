package cmd

import (
	"github.com/spf13/cobra"
)

var (
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
)

func imgList(cmd *cobra.Command, args []string) {
	req := getCliRequestIndexForApp(args, 0, flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().ImgList(client.Context(), req)
	uiCheckErr("Could not List the Products: %v", err)
	uiList(res)
}

func imgInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().ImgInfo(client.Context(), req)
	uiCheckErr("Could not Get the Image Info: %v", err)
	uiImageInfo(res)
}

func imgImport(cmd *cobra.Command, args []string) {} //TODO

func imgBuild(cmd *cobra.Command, args []string) {} //TODO

func imgDelete(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().ImgDelete(client.Context(), req)
	uiCheckErr("Could not Destroy the Image: %v", err)
	log.Printf("image %s deleted", req.Name)
}
