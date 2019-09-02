package cmd

import (
	"github.com/spf13/cobra"
)

// Repository Build
var (
	objectCmd = &cobra.Command{
		Use:   "object",
		Short: "subcommand and option switches for management on Object Storage",
		Long:  `This subcommand for management files on Object Storage.`,
		Example: `
		$: yb object cp "Source/Path" "Bucket_Name/Destination/Path"
		$: yb object rm "Bucket_Name/Object_Name"`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	objectRmCmd = &cobra.Command{
		Use:   "rm",
		Run:   objectRm,
		Short: "remove a file from Object Storage",
		Long:  `remove a file from Object Storage`,
		Example: `
		$: yb object rm "Bucket_Name/Object_Name"`,
	}
	objectCpCmd = &cobra.Command{
		Use:   "cp",
		Run:   objectCp,
		Short: "copy objects and option switches for management on Object Storage",
		Long:  `copy objects and option switches for management on Object Storage`,
		Example: `
		$: yb object cp "Source/Path" "Bucket_Name/Destination/Path"`,
	}
	objectLsCmd = &cobra.Command{
		Use:   "ls",
		Run:   objectLs,
		Short: "list of all on Object Storage by bucket name",
		Long:  `list of all on Object Storage by bucket name`,
		Example: `
		$: yb object ls "Bucket_Name"`,
	}
	bucketCmd = &cobra.Command{
		Use:   "bucket",
		Run:   bucketList,
		Short: "show all available bucket",
		Long:  `This subcommand shows list of available Object Storage Bucket.`,
		Example: `
		$: yb create bucket test-bucket-create
		$: yb object bucket `}
)

func init() {
	objectCmd.AddCommand(
		objectCpCmd,
		objectRmCmd,
		objectLsCmd,
		bucketCmd)

	rootCmd.AddCommand(objectCmd)
}
