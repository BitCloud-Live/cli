package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	setListCmd = &cobra.Command{
		Use:   "set:list",
		Short: "list of all settings",
		Long:  `This subcommand can pageing the setting name.`,
		Run:   setList}

	setInfoCmd = &cobra.Command{
		Use:   "set:info",
		Short: "detail of setting",
		Long:  `This subcommand show the information of a setting.`,
		Run:   setInfo}

	setAddCmd = &cobra.Command{
		Use:   "set:add",
		Short: "add a setting to the application",
		Long:  `This subcommand add a setting to the application.`,
		Run:   setAdd}

	setUpdateCmd = &cobra.Command{
		Use:   "set:update",
		Short: "update a setting for application",
		Long:  `This subcommand show the information of a setting.`,
		Run:   setUpdate}

	setDeleteCmd = &cobra.Command{
		Use:   "set:delete",
		Short: "delete a setting from application",
		Long:  `This subcommand delete a setting from application.`,
		Run:   setDelete}
)

func setList(cmd *cobra.Command, args []string) {
	req := reqIndexForApp(cmd)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SetList(client.Context(), req)
	uiCheckErr("Could not List the Settings: %v", err)
	uiList(res)
}

func setInfo(cmd *cobra.Command, args []string) {
	req := new(ybApi.SettingInfoReq)
	req.Name = cmd.Flag("name").Value.String()
	req.App = cmd.Flag("application").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SetInfo(client.Context(), req)
	uiCheckErr("Could not Get the Setting Info: %v", err)
	uiSettingByDetail(res)
}

func setAdd(cmd *cobra.Command, args []string) {
	var (
		err error
		req = new(ybApi.SettingReq)
	)
	req.Name = cmd.Flag("name").Value.String()
	req.App = cmd.Flag("application").Value.String()
	req.Path = cmd.Flag("path").Value.String()
	file := cmd.Flag("file").Value.String()
	req.File, err = readFile(file)
	uiCheckErr("Could not Read the File: %v", err)

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SetAdd(client.Context(), req)
	uiCheckErr("Could not Add the Setting: %v", err)
	uiSettingByDetail(res)
}

func setUpdate(cmd *cobra.Command, args []string) {
	var (
		err error
		req = new(ybApi.SettingReq)
	)
	req.Name = cmd.Flag("name").Value.String()
	req.App = cmd.Flag("application").Value.String()
	req.Path = cmd.Flag("path").Value.String()
	if file := cmd.Flag("file").Value.String(); file != "" {
		req.File, err = readFile(file)
		uiCheckErr("Could not Read the File: %v", err)
	}

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().SetUpdate(client.Context(), req)
	uiCheckErr("Could not Update the Setting: %v", err)
	uiSettingByDetail(res)
}

func setDelete(cmd *cobra.Command, args []string) {
	req := new(ybApi.SettingInfoReq)
	req.Name = cmd.Flag("name").Value.String()
	req.App = cmd.Flag("application").Value.String()

	client := grpcConnect()
	defer client.Close()
	_, err := client.V2().SetDelete(client.Context(), req)
	uiCheckErr("Could not Delete the Setting Info: %v", err)
	log.Println("Task is done.")
}

func init() {
	// setting list:
	setListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number list")
	setListCmd.Flags().StringVarP(&flagAppName, "app", "n", "", "page number list")

	// setting info:
	setInfoCmd.Flags().StringP("name", "n", "", "name of Setting")
	setInfoCmd.Flags().StringP("application", "a", "", "name of Application")
	setInfoCmd.MarkFlagRequired("application")

	// setting Add:
	setAddCmd.Flags().StringP("name", "n", "", "name of Setting")
	setAddCmd.Flags().StringP("application", "a", "", "name of Application")
	setAddCmd.Flags().StringP("path", "p", "", "mounted settings path")
	setAddCmd.Flags().StringP("file", "f", "", "settings file address")
	setAddCmd.MarkFlagRequired("name")
	setAddCmd.MarkFlagRequired("application")
	setAddCmd.MarkFlagRequired("path")
	setAddCmd.MarkFlagRequired("file")

	// setting Update:
	setUpdateCmd.Flags().StringP("name", "n", "", "name of Setting")
	setUpdateCmd.Flags().StringP("application", "a", "", "name of Application")
	setUpdateCmd.Flags().StringP("path", "p", "", "mounted settings path")
	setUpdateCmd.Flags().StringP("file", "f", "", "settings file address")
	setUpdateCmd.MarkFlagRequired("name")
	setUpdateCmd.MarkFlagRequired("application")

	// setting Delete:
	setDeleteCmd.Flags().StringP("name", "n", "", "name of Setting")
	setDeleteCmd.Flags().StringP("application", "a", "", "name of Application")
	setDeleteCmd.MarkFlagRequired("name")
	setDeleteCmd.MarkFlagRequired("application")

	rootCmd.AddCommand(
		setListCmd,
		setInfoCmd,
		setAddCmd,
		setUpdateCmd,
		setDeleteCmd)
}
