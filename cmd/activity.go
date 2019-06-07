package cmd

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	ybApi "github.com/yottab/proto-api/proto"
)

var (
	flagTag int32
)

func actList(cmd *cobra.Command, args []string) {
	req := &ybApi.ActivityReq{}
	strTag := getCliArg(args, 0, "0")
	tagVal, err := strconv.ParseUint(strTag, 0, 64)
	if err != nil {
		tagVal = 0
	}
	req.Tag = ybApi.ActivityTag(tagVal)
	req.Name = getCliArg(args, 1, "")
	req.Index = flagIndex
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
