package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Repository Build
var repositoryBuildCmd = &cobra.Command{
	Use:   "build",
	Run:   repositoryBuild,
	Short: "build a Repository from a Dockerfile",
	Long:  `This subcommand Build an image from a dockerfile.`,
	Example: `
		$: yb build \
			  --name=my-application \
			  --tag=0.0.1 \
			  --path="~/Desktop/my-application/"`}

// return current path
func getPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// return name of current folder
// the last element of path
func defaultRepositoryName() string {
	path := getPath()
	return filepath.Base(path)
}

func init() {
	repositoryBuildCmd.Flags().StringP("name", "n", defaultRepositoryName(), "the uniquely identifiable name for the Repository")
	repositoryBuildCmd.Flags().StringP("tag", "T", "latest", "the uniquely identifiable name for the Repository")
	repositoryBuildCmd.Flags().StringP("path", "p", getPath(), "the uniquely identifiable name for the Repository")

	rootCmd.AddCommand(repositoryBuildCmd)
}
