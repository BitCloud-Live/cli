package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kubernetes/kompose/pkg/kobject"
	"github.com/kubernetes/kompose/pkg/loader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: run command and arg

/*
support keyword:
	- image
	- volumes
	- labels
	- environment
	- ports

	// TODO
	- command
	- Expose
	- Dockerfile
	- healthcheck
*/

const (
	composeProductName    = "io.yottab.product"
	composeServicePlan    = "io.yottab.plan"
	composeServiceValues  = "io.yottab.values"
	composeDomainName     = "io.yottab.domain.name"
	composeDomainPath     = "io.yottab.domain.path"     // default is '/'
	composeDomainProtocol = "io.yottab.domain.protocol" // default is http
	composeDomains        = "io.yottab.domains"         // [{domain, path, endpoint}]
)

var (
	flagComposeFilesArray = make([]string, 0, 2)

	// ErrVolumeNotExist .
	ErrVolumeNotExist = errors.New("Docker Compose Err: Volume Not Exist")
	// ErrPortNotExist .
	ErrPortNotExist = errors.New("Docker Compose Err: Port Not Exist")
	// ErrVolumeBadFormat .
	ErrVolumeBadFormat = errors.New("Docker Compose Err: only support external volume")
)

func composeLoader(files []string) (ServiceConfigs map[string]kobject.ServiceConfig) {
	l, err := loader.GetLoader("compose")
	if err != nil {
		log.Fatal(err)
	}

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}

	komposeObject, err = l.LoadFile(files)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return komposeObject.ServiceConfigs
}

func checkComposeFormat(serviceConfigs map[string]kobject.ServiceConfig) {
	if err := applicationCheckExistPort(serviceConfigs); err != nil {
		log.Fatalf(err.Error())
	}

	if err := volumeCheck(serviceConfigs); err != nil {
		log.Fatalf(err.Error())
	}

	if err := domainCheckExist(serviceConfigs); err != nil {
		log.Fatalf(err.Error())
	}

	if err := composeProductCheck(serviceConfigs); err != nil {
		log.Fatalf(err.Error())
	}
}

// TODO: check product exist at Yb.product.list
func composeProductCheck(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		if !composeConfigIsApplication(conf) {
			if _, OK := conf.Labels[composeProductName]; !OK {
				log.Printf("Service %s not have ProductName at labels.%s", srv, composeProductName)
				return ErrPortNotExist
			}
		}
	}
	return nil
}

func applicationCheckExistPort(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) && len(conf.Port) == 0 {
			log.Printf("Application %s not have port", srv)
			return ErrPortNotExist
		}
	}
	return nil
}

func volumeCheck(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		// if image not null => srv.type is Application
		for _, v := range conf.Volumes {
			if v.VolumeName == "." ||
				strings.Contains(v.VolumeName, "/") ||
				strings.Contains(v.VolumeName, "\\") {
				log.Printf("Only support Volume, Err at [%s]", srv)
				return ErrVolumeBadFormat
			}
			volumeCheckExist(srv, v.VolumeName)
		}
	}
	return nil
}
func volumeCheckExist(srv, name string) {
	_, err := VolumeInfo(name)
	uiCheckErr("Volume not OK at Service "+srv, err)
}

func domainCheckExist(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			if domain, ok := conf.Labels[composeDomainName]; ok {
				_, err := DomainInfo(domain)
				uiCheckErr("Domain not OK at Service "+srv, err)
			}
		}
	}
	return nil
}

// TODO: panic if( Err != notFind )
func serviceIsExist(name string) bool {
	_, err := ServiceInfo(name)
	return err != nil
}

func getComposeProductName(conf kobject.ServiceConfig) string {
	if conf.Labels != nil {
		return conf.Labels[composeProductName]
	}

	return ""
}

func getComposeServiceValues(conf kobject.ServiceConfig) (val map[string]string) {
	val = make(map[string]string)
	if conf.Labels != nil {
		jsonVal, ok := conf.Labels[composeServiceValues]
		if ok {
			byteVal := []byte(jsonVal)
			if err := json.Unmarshal(byteVal, &val); err != nil {
				log.Fatalf(err.Error())
			}
		}
	}

	return
}

func getComposeProductPlan(conf kobject.ServiceConfig) string {
	if conf.Labels != nil {
		return conf.Labels[composeServicePlan]
	}

	return ""
}

func composeConfigIsApplication(conf kobject.ServiceConfig) bool {
	// if image not null => srv.type is Application
	return conf.Image != ""
}

func getEndPointProtocol(conf kobject.ServiceConfig) (protocol string) {
	protocol, OK := conf.Labels[composeDomainProtocol]
	if !OK {
		protocol = "http"
	}
	return
}

func getEndPoint(conf kobject.ServiceConfig) (endPoint string) {
	if len(conf.Port) != 1 || conf.Port[0].ContainerPort == 0 {
		log.Fatal("Err: Port format")
	}

	protocol := getEndPointProtocol(conf)
	return fmt.Sprintf("%s/%d", protocol, conf.Port[0].ContainerPort)
}

func getDomainVariable(conf kobject.ServiceConfig) (domain, path string) {
	domain = conf.Labels[composeDomainName]
	path, OK := conf.Labels[composeDomainPath]
	if !OK {
		path = "/"
	}
	return
}

func composeAppLinkDomain(configs map[string]kobject.ServiceConfig) {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			domain, path := getDomainVariable(conf)
			endpoint := getEndPoint(conf)
			_, err := AppAttachDomain(srv, domain, path, endpoint)
			uiCheckErr("Could not Add the Domain for Application", err)
		}
	}
}

func composeAppLinkVolume(configs map[string]kobject.ServiceConfig) {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			for _, vol := range conf.Volumes {
				_, err := AppAttachVolume(srv, vol.VolumeName, vol.MountPath)
				uiCheckErr("Could not Add the Volume for Application", err)
			}
		}
	}
}

// Load environment variables from file
func getEnvsFromFile(config kobject.ServiceConfig) map[string]string {
	env := make(map[string]string)

	for _, filename := range config.EnvFile {
		envLoad := parseDotEnvFile(filename)
		env = mapMerge(env, envLoad)
	}

	return env
}
func convConfigEnvVarToMap(config kobject.ServiceConfig) map[string]string {
	env := make(map[string]string, len(config.Environment))
	for _, v := range config.Environment {
		env[v.Name] = v.Value
	}
	return env
}
func composeAppLinkEnv(configs map[string]kobject.ServiceConfig) {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			env := convConfigEnvVarToMap(conf)
			envFile := getEnvsFromFile(conf)
			env = mapMerge(env, envFile)

			_, err := ApplicationAddEnvironmentVariable(srv, env)
			uiCheckErr("Could not Add the Environment Variable for Application", err)
		}
	}
}

// if not exist, create service
func composeCreateService(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		if !composeConfigIsApplication(conf) && !serviceIsExist(srv) {
			product := getComposeProductName(conf)
			plan := getComposeProductPlan(conf)
			val := getComposeServiceValues(conf)
			res, err := ServiceCreate(product, srv, plan, val)
			uiCheckErr("Could not Create the Service", err)
			uiServicStatus(res)
		}
	}
	return nil
}

func getReplicas(conf kobject.ServiceConfig) uint64 {
	if conf.Replicas == 0 {
		return 1
	}

	return uint64(conf.Replicas)
}

// if not exist, Create Application
func composeAppCreate(configs map[string]kobject.ServiceConfig) error {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			// TODO check err == not exist
			if _, err := AppInfo(srv); err != nil {
				plan := getComposeProductPlan(conf)
				endpointProtocol := getEndPointProtocol(conf)
				port := uint64(conf.Port[0].ContainerPort)
				minScale := getReplicas(conf)

				// TODO set otherValues
				res, err := ApplicationCreate(srv, conf.Image, plan, endpointProtocol, port, minScale, nil)
				uiCheckErr("Could not Create the Application", err)
				uiApplicationStatus(res)
			} else {
				log.Printf("Application %s is exist", srv)
			}
		}
	}
	return nil
}

func composeStart(cmd *cobra.Command, args []string) {
	configs := composeLoader(flagComposeFilesArray)

	// check metadata befor start
	checkComposeFormat(configs)

	// first create Service
	err := composeCreateService(configs)
	uiCheckErr("Could not Link Volume", err)

	// then create Application, add Link the (Volumes, Domains and Env)
	err = composeAppCreate(configs)
	uiCheckErr("Could not Link Volume", err)
	composeAppLinkVolume(configs)
	composeAppLinkDomain(configs)
	composeAppLinkEnv(configs)
}

// TODO
func composeStop(cmd *cobra.Command, args []string) {}
