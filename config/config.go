package config

import (
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	DEFAULTE_CONTOROLLER = "controller.uvcloud.ir:8443"
	KEY_EMAIL            = "email"
	KEY_TOKEN            = "token"
	KEY_LINK             = "link"
	APP_NAME             = "uv"
	CONFIG_NAME          = "config"
)

var (
	// Get config file Path from the flag.
	ConfigManualAddress = ""

	// Find home directory for definition of archive folder
	// configPath == $HOME/.uv/
	configPath = filepath.Join(getHome(), "."+APP_NAME)
)

// Find home directory.
func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func UpdateVarByConfigFile() {
	// read config either from Flag --config or Path "$HOME/.uv/config.json"
	if ConfigManualAddress != "" {
		// Use config file from the flag.
		viper.SetConfigFile(ConfigManualAddress)
	} else {
		// Search config in home directory
		viper.AddConfigPath(configPath)
		viper.SetConfigName(CONFIG_NAME)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning - Can't read config file:", err)
	}
}

func ResetConfigFile() (err error) {
	return viper.WriteConfig()
}
