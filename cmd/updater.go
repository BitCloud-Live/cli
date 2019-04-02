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
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Self update cli to latest stable release",
		Long:  `Self update cli to latest stable release from github release page of YOTTAb cli`,
		Run:   confirmAndSelfUpdate}
)

func confirmAndSelfUpdate(cmd *cobra.Command, args []string) {
	latest, found, err := selfupdate.DetectLatest("yottab/cli")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}

	//Convert to a semver compatible version
	semverCompatVersion := strings.TrimPrefix(version, "v")
	v := semver.MustParse(semverCompatVersion)
	log.Printf("Current version: %s", v)
	if !found || latest.Version.Equals(v) || latest.Version.LT(v) {
		log.Printf("Latest stable version from upstream (github): %s", latest.Version)
		log.Print("Current version is the latest")
		return
	}

	fmt.Printf("We found a newer version: %s", latest.Version)
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

	if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0]); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
