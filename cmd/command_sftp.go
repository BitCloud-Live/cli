package cmd

import (
	"github.com/spf13/cobra"
)

// sftp Command
var (
	sftpCmd = &cobra.Command{
		Use:   "sftp [volume-name]",
		Short: "Connect to a sftp agent running in a cluster",
		Long: `Connect to a sftp agent running in a cluster.
		This type of connection can be useful for volume mounting`,
		Run: volumeSftp,
	}
)

func init() {
	rootCmd.AddCommand(
		sftpCmd)
}
