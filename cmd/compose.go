package cmd

// import (
// 	"io/ioutil"
// 	"os"

// 	"github.com/aanand/compose-file/loader"
// 	"github.com/aanand/compose-file/types"
// 	"github.com/spf13/cobra"
// )

// func loadConfigFile(fileYmlName string) (*types.Config, error) {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}

// 	bytes, err := ioutil.ReadFile(fileYmlName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dict, err := loader.ParseYAML(bytes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return loader.Load(types.ConfigDetails{
// 		WorkingDir: workingDir,
// 		ConfigFiles: []types.ConfigFile{
// 			{Filename: "filename.yml", Config: dict},
// 		},
// 		Environment: nil,
// 	})
// }

// func getDockerComposeObject(cmd *cobra.Command) *types.Config {
// 	ymlFile := cmd.Flag("compose-file").Value.String()
// 	confObj, err := loadConfigFile(ymlFile)
// 	if err != nil {
// 		log.Fatal("Err load File: ", err)
// 	}
// 	return confObj
// }

// func composeCreate(cmd *cobra.Command, args []string) {
// 	confObj := getDockerComposeObject(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// 	/*
// 		for _, srv := range confObj.Services {
// 			log.Println("Name", srv.Name)

// 			log.Println("Volumes", srv.Volumes)
// 			log.Println("DomainName", srv.DomainName) // string
// 			log.Println("Image", srv.Image) // string
// 			log.Println("Labels", srv.Labels) // map[string]string
// 			//log.Println("Deploy", srv.Deploy) // DeployConfig
// 			log.Println("DependsOn", srv.DependsOn)
// 			//log.Println("Ports", srv.Ports)

// 			log.Println("Environment", srv.Environment) // map[string]string

// 			log.Println("Links", srv.Links)
// 			// Command         []string
// 			// Ports           []string
// 			// Restart         string
// 			// User            string
// 			// Entrypoint      []string
// 			// Expose          []string
// 			// ExternalLinks   []string
// 			// Links           []string
// 			// ExtraHosts      map[string]string
// 			// Hostname        string
// 			// HealthCheck     *HealthCheckConfig
// 			// Logging         *LoggingConfig
// 			// NetworkMode     string
// 		}

// 		log.Println("-----------------------------------------------")

// 		for key, vol := range confObj.Volumes {
// 			log.Println("key", key)
// 			_, err := VolumeCreate(key, "persistant-2Gi")

// 			log.Println("Err", err)

// 			//log.Println("Driver", vol.Driver)
// 			log.Println("DriverOpts", vol.DriverOpts)

// 			//log.Println("Labels", vol.Labels)
// 			//log.Println("External", vol.External)
// 		}
// 	*/
// }

// func composeStart(cmd *cobra.Command, args []string) {
// 	confObj := getDockerComposeObject(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// }

// func composeStop(cmd *cobra.Command, args []string) {
// 	confObj := getDockerComposeObject(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// }
