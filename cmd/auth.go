package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	uvApi "github.com/uvcloud/uv-api-go/proto"
	"github.com/uvcloud/uv-cli/config"
)

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to a controller",
		Long:  `This subcommand logs in by authenticating against a controller.`,
		Run:   login}

	logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "logout from the uvCloud",
		Long:  `This subcommand logs out from a controller and clears the user session.`,
		Run:   logout}
)

func login(cmd *cobra.Command, args []string) {
	email := viper.GetString(config.KEY_USER)

	if len(email) < 5 {
		email = readFromConsole("Username: ")
	}
	password := readPasswordFromConsole("Password: ")
	client := grpcConnect()
	defer client.Close()
	req := &uvApi.LoginReq{Email: email, Password: password}
	res, err := client.V1().Login(client.Context(), req)
	if err != nil {
		log.Fatalf("Could not Login: %v", err)
	}

	log.Printf("Login successful!")
	viper.Set(config.KEY_TOKEN, res.Token)

	// Save TOKEN to config file
	if err = config.ResetConfigFile(); err != nil {
		log.Fatalf("Could not Save ConfigFile: %v", err)
	}
}

func logout(cmd *cobra.Command, args []string) {
	req := &uvApi.Empty{}
	client := grpcConnect()
	defer client.Close()
	_, err := client.V1().Logout(client.Context(), req)
	if err != nil {
		log.Fatalf("Could not Logout: %v", err)
	}

	viper.Set(config.KEY_TOKEN, "")
}

func init() {
	rootCmd.AddCommand(loginCmd, logoutCmd)
}
