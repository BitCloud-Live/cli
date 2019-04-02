package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
)

// RootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     config.APP_NAME,
		Short:   "YOTTAb cli",
		Long:    `YOTTAb cli for client side usage`,
		Version: version,
	}
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

	rootCmd.PersistentFlags().StringVar(&config.ConfigManualAddress, "config", "", "Config file (default is $HOME/.yb/config.json)")
	rootCmd.PersistentFlags().StringP(config.KEY_USER, "u", "", "Provide a username for the new account")
	rootCmd.PersistentFlags().StringP(config.KEY_HOST, "l", config.DEFAULTE_CONTOROLLER, "Address of Controller. a fully-qualified controller URI")
	rootCmd.PersistentFlags().StringP(config.KEY_TOKEN, "t", "", "Manual Send 'TOKEN' for Authentication")

	viper.BindPFlag(config.KEY_USER, rootCmd.PersistentFlags().Lookup(config.KEY_USER))
	viper.BindPFlag(config.KEY_HOST, rootCmd.PersistentFlags().Lookup(config.KEY_HOST))
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
