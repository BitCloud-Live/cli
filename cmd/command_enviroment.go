package cmd

import (
	"github.com/spf13/cobra"
)

var (
	commentShort     = "environment variables [Set/Unset] for an Application"
	commentSetLong   = `Add list of environment variables to an Application.`
	commentUnsetLong = `Remove list of environment variables for an Application.`
	commentExample   = `
  $: yb environment unset my-admin
        -v="key1"
		-v="key2"

  $: yb environment set my-admin
        -v="key1=value1"
        -v="key2=value2"`
)

var (
	envCmd = &cobra.Command{
		Use:     "environment [command]",
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		}}

	setEnvCmd = &cobra.Command{
		Use:     "set [app.name]",
		Run:     appAddEnvironmentVariable,
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample}

	unsetEnvCmd = &cobra.Command{
		Use:     "unset [app.name]",
		Run:     appRemoveEnvironmentVariable,
		Short:   commentShort,
		Long:    commentUnsetLong,
		Example: commentExample}
)

func init() {
	// app set flag:
	setEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	setEnvCmd.MarkFlagRequired("variable")

	// app unset flag:
	unsetEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	unsetEnvCmd.MarkFlagRequired("variable")

	// Add the commands Environment Variable
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(
		setEnvCmd,
		unsetEnvCmd)
}
