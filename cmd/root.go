package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uvcloud/uv-cli/config"
)

// RootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   config.APP_NAME,
		Short: "UVCloud cli",
		Long:  `UVCloud cli for client side usage`}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Here you will define your flags and configuration settings.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&config.ConfigManualAddress, "config", "", "Config file (default is $HOME/.uv/config.json)")
	rootCmd.PersistentFlags().StringP(config.KEY_EMAIL, "e", "", "Provide a Email address for the new account")
	rootCmd.PersistentFlags().StringP(config.KEY_LINK, "l", config.DEFAULTE_CONTOROLLER, "Address of Controller. a fully-qualified controller URI")
	rootCmd.PersistentFlags().StringP(config.KEY_TOKEN, "t", "", "Manual Send 'TOKEN' for Authentication")

	viper.BindPFlag(config.KEY_EMAIL, rootCmd.PersistentFlags().Lookup(config.KEY_EMAIL))
	viper.BindPFlag(config.KEY_LINK, rootCmd.PersistentFlags().Lookup(config.KEY_LINK))
	viper.BindPFlag(config.KEY_TOKEN, rootCmd.PersistentFlags().Lookup(config.KEY_TOKEN))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.UpdateVarByConfigFile()
	command := os.Args[1]
	token := viper.GetString(config.KEY_TOKEN)
	if len(token) == 0 && command != "login" {
		log.Fatal("You must login first.")
	}
}
