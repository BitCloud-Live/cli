package cmd

import (
	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
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
		Long:  `This subcommand adds a new worker to an application, be notify that there is a upper limit to worker counts.`,
		Run:   workerCreate}

	workerUpdateCmd = &cobra.Command{
		Use:   "worker:update",
		Short: "Update an existing worker for an application",
		Long:  `This subcommand Update an existing worker for an application.`,
		Run:   workerUpdate}

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
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("worker").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerInfo(client.Context(), req)
	uiCheckErr("Could not Get Application: %v", err)
	uiWorkerStatus(res)
}

func workerPortforward(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("worker").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerPortforward(client.Context(), req)
	uiCheckErr("Could not Portforward the Service: %v", err)
	uiPortforward(res)
}

func workerDestroy(cmd *cobra.Command, args []string) {
	req := new(ybApi.AttachIdentity)
	req.Name = cmd.Flag("name").Value.String()
	req.Attachment = cmd.Flag("worker").Value.String()

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
	req.Identities.Name = cmd.Flag("name").Value.String()
	req.Identities.Attachment = cmd.Flag("worker").Value.String()
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
	req.Identities.Name = cmd.Flag("name").Value.String()
	req.Identities.Attachment = cmd.Flag("worker").Value.String()
	req.Config = new(ybApi.WorkerConfig)
	req.Config.Port = flagVarPort
	req.Config.Image = cmd.Flag("image").Value.String()

	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().AppWorkerUpdate(client.Context(), req)
	uiCheckErr("Could not Update the Application: %v", err)
	uiList(res)
}

func init() {
	// worker List:
	workerListCmd.Flags().Int32Var(&flagIndex, "index", 0, "page number list")
	workerListCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerListCmd.MarkFlagRequired("name")

	// worker Info:
	workerInfoCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerInfoCmd.Flags().StringP("worker", "w", "", "name of worker")
	workerInfoCmd.MarkFlagRequired("name")
	workerInfoCmd.MarkFlagRequired("worker")

	// worker portforward:
	workerPortforwardCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the application.")
	workerPortforwardCmd.Flags().StringP("worker", "w", "", "name of worker")
	workerPortforwardCmd.MarkFlagRequired("name")
	workerPortforwardCmd.MarkFlagRequired("worker")

	// worker Create:
	workerCreateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	workerCreateCmd.Flags().StringP("worker", "w", "", "name of worker")
	workerCreateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	workerCreateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	workerCreateCmd.MarkFlagRequired("image")
	workerCreateCmd.MarkFlagRequired("port")
	workerCreateCmd.MarkFlagRequired("name")
	workerCreateCmd.MarkFlagRequired("worker")

	// worker Update:
	workerUpdateCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	workerUpdateCmd.Flags().StringP("worker", "w", "", "name of worker")
	workerUpdateCmd.Flags().Uint64VarP(&flagVarPort, "port", "p", 0, "port of application")
	workerUpdateCmd.Flags().StringVarP(&flagVarImage, "image", "i", "", "image of application")
	workerUpdateCmd.MarkFlagRequired("name")
	workerUpdateCmd.MarkFlagRequired("worker")

	// worker Destroy:
	workerDestroyCmd.Flags().StringP("name", "n", "", "a uniquely identifiable name for the application. No other app can already exist with this name.")
	workerDestroyCmd.Flags().StringP("worker", "w", "", "name of worker")
	workerDestroyCmd.MarkFlagRequired("name")
	workerDestroyCmd.MarkFlagRequired("worker")

	rootCmd.AddCommand(
		workerListCmd,
		workerInfoCmd,
		workerPortforwardCmd,
		workerCreateCmd,
		workerUpdateCmd,
		workerDestroyCmd,
	)
}
