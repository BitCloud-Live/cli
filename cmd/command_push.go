package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	flagGitCommitHash bool
	flagGitTag        bool
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// return name of current folder
// the last element of path
func defaultRepositoryName(path string) string {
	var ss []string
	if runtime.GOOS == "windows" {
		ss = strings.Split(path, "\\")
	} else {
		ss = strings.Split(path, "/")
	}

	currentDirName := ss[len(ss)-1]

	return currentDirName

}

func init() {
	repositoryPushLogCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the Repository")
	repositoryPushLogCmd.Flags().StringP("tag", "T", "latest", "the uniquely identifiable name for the Repository")

	repositoryPushCmd.Flags().StringP("name", "n", "", "the uniquely identifiable name for the Repository")
	repositoryPushCmd.Flags().StringP("tag", "T", "latest", "docker image tag")
	repositoryPushCmd.Flags().BoolVar(&flagGitTag, "tag-git", false, "set docker image tag from recent git tag")
	repositoryPushCmd.Flags().BoolVar(&flagGitCommitHash, "tag-commit", false, "set docker image tag from git commit hash")
	repositoryPushCmd.Flags().StringP("path", "p", getPath(), "the uniquely identifiable name for the Repository")

	repositoryPushCmd.AddCommand(repositoryPushLogCmd)
	rootCmd.AddCommand(repositoryPushCmd)
}
