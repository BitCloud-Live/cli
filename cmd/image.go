package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

func imgList(cmd *cobra.Command, args []string) {
	req := getCliRequestIndexForApp(args, 0, flagIndex)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().ImgList(client.Context(), req)
	uiCheckErr("Could not List the Products", err)
	uiList(res)
}

func imgInfo(cmd *cobra.Command, args []string) {
	req := getCliRequestIdentity(args, 0)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().ImgInfo(client.Context(), req)
	uiCheckErr("Could not Get the Image Info", err)
	uiImageInfo(res)
}

func imgDelete(cmd *cobra.Command, args []string) {
	imageTag := cmd.Flag("tag").Value.String()    // Repository tag
	imageName := cmd.Flag("name").Value.String()  // Repository name

	req := new(ybApi.ImgBuildReq)
	req.RepositoryName = imageName
	req.RepositoryTag = imageTag

	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().ImgDelete(client.Context(), req)
	uiCheckErr("Could not Destroy the Image", err)
	log.Printf("image %s:%s deleted", req.RepositoryName, req.RepositoryTag)
}
