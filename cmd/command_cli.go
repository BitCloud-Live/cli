package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

//version parameter parsing from compiler ldflags
var version string

var (
	cliCmd = &cobra.Command{
		Use:   "cli [command]",
		Short: "YOTTAb cli setting",
		Long:  ``}

	cliUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Self update cli to latest stable release",
		Long:  `Self update cli to latest stable release from github release page of YOTTAb cli`,
		Run:   checkAndSelfUpdate}

	cliCompletionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run
			
			. <(yb completion)
			
			To configure your bash shell to load completions for each session add to your bashrc
			
			# ~/.bashrc or ~/.profile
			. <(yb completion)
			`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO

			rootCmd.GenBashCompletion(os.Stdout)
		}}
)

// CheckNewerVersion .
func CheckNewerVersion(inform bool) (*selfupdate.Release, bool) {
	latest, found, err := selfupdate.DetectLatest("yottab/cli")
	if err != nil {
		if inform {
			log.Println("Error occurred while detecting version:", err)
		}
		return nil, false
	}

	//Convert to a semver compatible version
	semverCompatVersion := strings.TrimPrefix(version, "v")
	v := semver.MustParse(semverCompatVersion)
	if inform {
		log.Printf("Current version: %s", v)
	}
	if !found || latest.Version.Equals(v) || latest.Version.LT(v) {
		if inform {
			log.Printf("Latest stable version from upstream (github): %s", latest.Version)
			log.Print("Current version is the latest")
		}
		return nil, false
	}
	return latest, true
}

func checkAndSelfUpdate(cmd *cobra.Command, args []string) {
	latest, available := CheckNewerVersion(true)
	if !available {
		return
	}
	SelfUpdate(latest)
}

// SelfUpdate .
func SelfUpdate(latestAvailable *selfupdate.Release) {
	fmt.Printf("We found a newer version: %s", latestAvailable.Version)
	fmt.Print("Do you want to update? [Y/n]: ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Error("Error while reading from stdin!")
		return
	}
	input = strings.ToLower(strings.TrimSpace(input))
	if input != "y" && input != "" {
		log.Print("Update canceled!")
		return
	}

	if err := selfupdate.UpdateTo(latestAvailable.AssetURL, os.Args[0]); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latestAvailable.Version)
}

func init() {
	rootCmd.AddCommand(cliCmd)

	cliCmd.AddCommand(
		cliUpdateCmd,
		cliCompletionCmd)
}
