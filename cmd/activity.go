package cmd

import (
	"sort"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagTag int32
)

var (
	actListCmd = &cobra.Command{
		Use:   "act:list",
		Short: "show all activities",
		Long:  `This subcommand shows list of activities.`,
		Run:   actList}
	actTagListCmd = &cobra.Command{
		Use:   "act:tags",
		Short: "show all available tags",
		Long:  `This subcommand shows list of available activity tags.`,
		Run:   actTags}
)

func actList(cmd *cobra.Command, args []string) {
	req := &ybApi.ActivityReq{}
	req.Name = cmd.Flag("name").Value.String()
	req.Tag = ybApi.ActivityTag(flagTag)
	client := grpcConnect()
	defer client.Close()
	res, err := client.V2().ActivityList(client.Context(), req)
	uiCheckErr("Could not get the list of activities: %v", err)
	uiList(res)
}

func actTags(cmd *cobra.Command, args []string) {
	// To store the keys in ybApi.ActivityTag_name in sorted order
	var keys []int
	for k := range ybApi.ActivityTag_name {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		log.Printf("%d: %s", k, ybApi.ActivityTag_name[int32(k)])
	}

}

func init() {
	// imgage List:
	actListCmd.Flags().Int32VarP(&flagIndex, "index", "i", 0, "page number list")
	actListCmd.Flags().StringP("name", "n", "", "name filter for activity")
	actListCmd.Flags().Int32VarP(&flagTag, "tag-id", "d", 0, "see act:tags for all available tags, default to 0 (None)")

	rootCmd.AddCommand(
		actListCmd,
		actTagListCmd,
	)
}
