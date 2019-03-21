package cmd

import (
	"github.com/spf13/cobra"
	uvApi "github.com/uvcloud/uv-api-go/proto"
)

var (
	workerListCmd = &cobra.Command{
		Use:   "worker:list",
		Short: "list worker of an application",
		Long:  `Lists application workers visible to the current user.`,
		Run:   workerList}

	workerInfoCmd = &cobra.Command{
		Use:   "worker:info",
		Short: "view info about an application worker",
		Long:  `This subcommand prints info about some application worker`,
		Run:   workerInfo}

	workerCreateCmd = &cobra.Command{
		Use:   "worker:create",
		Short: "Add a new worker to an application",
		Long:  `This subcommand adds a new worker to an application, be notify that there is a upper limit to worker counts`,
		Run:   workerCreate}

	workerDestroyCmd = &cobra.Command{
		Use:   "worker:destroy",
		Short: "Delete an application worker",
		Long:  `This subcommand delete a worker from an application`,
		Run:   workerDestroy}

	workerPortforwardCmd = &cobra.Command{
		Use:   "worker:portforward",
		Short: "port-forward to connect to an application worker running in a cluster",
		Long: `Port-forward to connect to an application worker running in a cluster.
		This type of connection can be useful for admin panels, monitoring tools`,
		Run: workerPortforward}
)

func workerList(cmd *cobra.Command, args []string) {
	req := reqIdentity(cmd)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppListWorkers(client.Context(), req)
	uiCheckErr("Could not List the Applications: %v", err)
	uiList(res)
}

func workerInfo(cmd *cobra.Command, args []string) {
	req := new(uvApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiWorkerStatus(res)
}

func workerPortforward(cmd *cobra.Command, args []string) {
	req := new(uvApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func workerDestroy(cmd *cobra.Command, args []string) {
	req := new(uvApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("attachment").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppRemoveWorker(client.Context(), req)
	uiCheckErr("Could not Destroy the Application: %v", err)
	log.Printf("app %s deleted", req.Name)
	uiList(res)
}

func workerCreate(cmd *cobra.Command, args []string) {
	req := new(uvApi.WorkerReq)
	req.Identities = new(uvApi.AttachIdentity)
	req.Identities.Name = cmd.Flag("name").Value.String()
	req.Identities.Attachment = cmd.Flag("attachment").Value.String()
	req.Config = new(uvApi.WorkerConfig)
	req.Config.Port = flagVarPort
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppAddWorker(client.Context(), req)
	uiCheckErr("Could not Create the Application: %v", err)
	uiList(res)
}

func init() {
	// app List:
	workerListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")
	workerListCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerListCmd.MarkFlagRequired("name")

	// app Info:
	workerInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerInfoCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	workerInfoCmd.MarkFlagRequired("name")
	workerInfoCmd.MarkFlagRequired("attachment")

	// worker portforward:
	workerPortforwardCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerPortforwardCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	workerPortforwardCmd.MarkFlagRequired("name")
	workerPortforwardCmd.MarkFlagRequired("attachment")

	// worker Create:
	workerCreateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	workerCreateCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	workerCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	workerCreateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	workerCreateCmd.MarkFlagRequired("image")
	workerCreateCmd.MarkFlagRequired("port")
	workerCreateCmd.MarkFlagRequired("name")
	workerCreateCmd.MarkFlagRequired("attachment")

	// worker Destroy:
	workerDestroyCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	workerDestroyCmd.Flags().StringP("attachment", "a", "", "name of attachment")
	workerDestroyCmd.MarkFlagRequired("name")
	workerDestroyCmd.MarkFlagRequired("attachment")

	rootCmd.AddCommand(
		workerListCmd,
		workerInfoCmd,
		workerPortforwardCmd,
		workerCreateCmd,
		workerDestroyCmd,
	)
}
