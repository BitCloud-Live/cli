package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Repository Push
var repositoryPushCmd = &cobra.Command{
	Use:   "push",
	Run:   pushRepository,
	Short: "build a Repository from a Dockerfile",
	Long:  `This subcommand Build an image from a dockerfile.`,
	Example: `
		$: yb push \
			  --name=my-application \
			  --tag=0.0.1 \
			  --path="~/Desktop/my-application/"`,
}

var repositoryPushLogCmd = &cobra.Command{
	Use:   "log",
	Run:   pushLog,
	Short: "tail builder log",
	Long:  `This subcommand Tail current builder log.`,
	Example: `
		$: yb push log \
			  --name=my-application \
			  --tag=0.0.1`,
}

// return current path
func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s%c", dir, os.PathSeparator)
}

// return name of current folder
// the last element of path
func defaultRepositoryName() string {
	path := getPath()
	return filepath.Base(path)
}

func init() {
	repositoryPushLogCmd.Flags().StringP("name", "n", defaultRepositoryName(), "the uniquely identifiable name for the Repository")
	repositoryPushLogCmd.Flags().StringP("tag", "T", "latest", "the uniquely identifiable name for the Repository")

	repositoryPushCmd.Flags().StringP("name", "n", defaultRepositoryName(), "the uniquely identifiable name for the Repository")
	repositoryPushCmd.Flags().StringP("tag", "T", "latest", "the uniquely identifiable name for the Repository")
	repositoryPushCmd.Flags().StringP("path", "p", getPath(), "the uniquely identifiable name for the Repository")

	repositoryPushCmd.AddCommand(repositoryPushLogCmd)
	rootCmd.AddCommand(repositoryPushCmd)
}
