package cmd

// import (
// "io/ioutil"
// "os"

// "github.com/aanand/compose-file/loader"
// "github.com/aanand/compose-file/types"
// 	"github.com/spf13/cobra"
// )

// func buildConfigDetails(source types.Dict) types.ConfigDetails {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return types.ConfigDetails{
// 		WorkingDir: workingDir,
// 		ConfigFiles: []types.ConfigFile{
// 			{Filename: "filename.yml", Config: source},
// 		},
// 		Environment: nil,
// 	}
// }

// func load(fileYmlName string) (*types.Config, error) {
// 	bytes, err := ioutil.ReadFile(fileYmlName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dict, err := loader.ParseYAML(bytes)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return loader.Load(buildConfigDetails(dict))
// }

// func getYamlFilePath(cmd *cobra.Command) *types.Config {
// 	ymlFile := cmd.Flag("compose-file").Value.String()
// 	confObj, err := load(ymlFile)
// 	if err != nil {
// 		log.Fatal("Err load File: ", err)
// 	}
// 	return confObj
// }

// func composeCreate(cmd *cobra.Command, args []string) {
// 	confObj := getYamlFilePath(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// 	/*
// 		for _, srv := range confObj.Services {
// 			log.Println("Name", srv.Name)

// 			log.Println("Volumes", srv.Volumes)
// 			log.Println("DomainName", srv.DomainName)

// 			log.Println("Image", srv.Image)
// 			log.Println("Labels", srv.Labels)
// 			//log.Println("Deploy", srv.Deploy)
// 			log.Println("DependsOn", srv.DependsOn)
// 			//log.Println("Ports", srv.Ports)

// 			log.Println("Environment", srv.Environment)

// 			log.Println("Links", srv.Links)
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
// 	confObj := getYamlFilePath(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// }

// func composeStop(cmd *cobra.Command, args []string) {
// 	confObj := getYamlFilePath(cmd)
// 	log.Println("TODO: Kill Me", confObj)
// 	// TODO
// }
