package config

import (
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	DEFAULTE_CONTOROLLER = "controller.uvcloud.ir:8443"
	KEY_USER             = "uv_user"
	KEY_TOKEN            = "token"
	KEY_HOST             = "host"
	APP_NAME             = "uv"
	CONFIG_NAME          = "config.json"
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
		viper.SetConfigFile(filepath.Join(configPath, CONFIG_NAME))
	}
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found!")
		log.Printf("Creating a new config file in %s", viper.ConfigFileUsed())
		viper.WriteConfig()
		if err := viper.ReadInConfig(); err != nil {
			log.Print("Failed!")
		}
		log.Print("Succeed")
	}
}

func ResetConfigFile() (err error) {
	return viper.WriteConfig()
}
