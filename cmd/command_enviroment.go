package cmd

import (
	"github.com/spf13/cobra"
)

// Comment
var (
	commentShort     = "environment variables [Set/Unset] for an Application or Service"
	commentSetLong   = `Add list of environment variables to an Application or Service.`
	commentUnsetLong = `Remove list of environment variables for an Application or Service.`
	commentExample   = `
  $: yb environment unset application <app_name>
        -v="<key>"
		-v="<key>"

  $: yb environment set application  <app_name>
        -v="<key>=<value>"
        -v="<key>=<value>"`
)

// Set Environment
var (
	envCmd = &cobra.Command{
		Use:     "environment [command]",
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample}

	setEnvCmd = &cobra.Command{
		Use:     "set [type]",
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample}

	appSetEnvCmd = &cobra.Command{
		Use:     "application [name]",
		Run:     appAddEnvironmentVariable,
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample}

	srvSetEnvCmd = &cobra.Command{
		Use:     "service [name]",
		Run:     srvEnvironmentSet,
		Short:   commentShort,
		Long:    commentSetLong,
		Example: commentExample}
)

// UnSet Environment
var (
	unsetEnvCmd = &cobra.Command{
		Use:     "unset [type]",
		Short:   commentShort,
		Long:    commentUnsetLong,
		Example: commentExample}

	appUnsetEnvCmd = &cobra.Command{
		Use:     "application [name]",
		Run:     appRemoveEnvironmentVariable,
		Short:   commentShort,
		Long:    commentUnsetLong,
		Example: commentExample}

	srvUnsetEnvCmd = &cobra.Command{
		Use:     "service [name]",
		Run:     srvEnvironmentUnset,
		Short:   commentShort,
		Long:    commentUnsetLong,
		Example: commentExample}
)

func init() {
	// srv set flag:
	srvSetEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of the service")
	srvSetEnvCmd.MarkFlagRequired("variable")

	// srv unset flag:
	srvUnsetEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of the service")
	srvUnsetEnvCmd.MarkFlagRequired("variable")

	// app set flag:
	appSetEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	appSetEnvCmd.MarkFlagRequired("variable")

	// app unset flag:
	appUnsetEnvCmd.Flags().StringArrayVarP(&flagVariableArray, "variable", "v", nil, "Environment Variable of application")
	appUnsetEnvCmd.MarkFlagRequired("variable")

	// Add the commands Environment Variable
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(
		setEnvCmd,
		unsetEnvCmd)

	setEnvCmd.AddCommand(
		appSetEnvCmd,
		srvSetEnvCmd)

	unsetEnvCmd.AddCommand(
		appUnsetEnvCmd,
		srvUnsetEnvCmd)
}
