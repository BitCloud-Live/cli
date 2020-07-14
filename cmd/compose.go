package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/kubernetes/kompose/pkg/kobject"
	"github.com/kubernetes/kompose/pkg/loader"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// ErrSrviceNotExist .
	ErrSrviceNotExist = errors.New("Docker Compose Err: Service Not Exist")
	// ErrDomainNotExist .
	ErrDomainNotExist = errors.New("Docker Compose Err: Domain Not Exist")
)

func composeLoader(files []string) (ServiceConfigs map[string]kobject.ServiceConfig) {
	l, err := loader.GetLoader("compose")
	uiCheckErr("Compose Loader Object Err", err)

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}

	komposeObject, err = l.LoadFile(files)
	uiCheckErr("Compose LoadFiles Err", err)

	return komposeObject.ServiceConfigs
}

func checkComposeFormat(serviceConfigs map[string]kobject.ServiceConfig) {
	err := applicationCheckExistPort(serviceConfigs)
	uiCheckErr("application Check Exist Port ", err)

	err = volumeCheck(serviceConfigs)
	uiCheckErr("volume Check", err)

	err = appCheckExistDomain(serviceConfigs)
	uiCheckErr("domain Check", err)

	err = composeProductCheck(serviceConfigs)
	uiCheckErr("product Check", err)
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
			log.Printf("Volume Check at service %s", srv)
			if v.VolumeName == "." ||
				strings.Contains(v.VolumeName, "/") ||
				strings.Contains(v.VolumeName, "\\") {
				log.Printf("Only support Volume, Err at [%s]", srv)
				return ErrVolumeBadFormat

			} else if !volumeCheckExist(v.VolumeName) {
				log.Printf("Volume '%s' Not Exist", v.VolumeName)
				return ErrVolumeNotExist
			}

			log.Printf("Volume %s is OK", v.VolumeName)
		}
	}
	return nil
}
func volumeCheckExist(name string) bool {
	_, err := VolumeInfo(name)
	return !errIsGrpcNotFound(name, err)
}

func serviceIsExist(name string) bool {
	_, err := ServiceInfo(name)
	return !errIsGrpcNotFound(name, err)
}

func appIsExist(name string) bool {
	_, err := AppInfo(name)
	return !errIsGrpcNotFound(name, err)
}

// if have domain, check exist
func appCheckExistDomain(configs map[string]kobject.ServiceConfig) error {
	for _, conf := range configs {
		if composeConfigIsApplication(conf) {
			if domain, ok := conf.Labels[composeDomainName]; ok {
				_, err := DomainInfo(domain)
				if errIsGrpcNotFound(domain, err) {
					return ErrDomainNotExist
				}
			}
		}
	}
	return nil
}

func checkExistLinkVolume(app, volume string) bool {
	List, err := VolumeList(app, 0)
	uiCheckErr("Err At check Linked The Volume to Applicaation", err)
	for _, v := range List.Volumes {
		if volume == v.Name {
			return true
		}
	}
	return false
}

func checkExistLinkDomain(app, domain string) bool {
	List, err := DomainList(app, 0)
	uiCheckErr("Err At check Linked The Domain to Applicaation", err)
	for _, d := range List.Domains {
		if domain == d.Domain {
			return true
		}
	}
	return false
}

func errIsGrpcNotFound(name string, err error) bool {
	errStatus, isOK := status.FromError(err)
	if !isOK {
		log.Fatalf("Check '%s' is not OK, Err: %v", name, err)
	}

	return errStatus.Code() == codes.NotFound
}

func getComposeProductName(conf kobject.ServiceConfig) string {
	if conf.Labels != nil {
		return conf.Labels[composeProductName]
	}

	return ""
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
	domain, OK := conf.Labels[composeDomainName]
	if !OK {
		return
	}
	path, OK = conf.Labels[composeDomainPath]
	if !OK {
		path = "/"
	}
	return
}

func composeAppLinkDomain(configs map[string]kobject.ServiceConfig) {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			domain, path := getDomainVariable(conf)
			if domain != "" {
				endpoint := getEndPoint(conf)
				if !checkExistLinkDomain(srv, domain) {
					_, err := AppAttachDomain(srv, domain, path, endpoint)
					if err != nil {
						log.Fatalf("Could not Link the Domain '%s' at '%s' and path '%s', to Application '%s': %v",
							domain, path, endpoint, srv, err)
					}
				}
				log.Printf("Domain '%s' at '%s' and path '%s', Linked to Application '%s'", domain, path, endpoint, srv)
			}
		}
	}
}

func composeAppLinkVolume(configs map[string]kobject.ServiceConfig) {
	for srv, conf := range configs {
		if composeConfigIsApplication(conf) {
			for _, vol := range conf.Volumes {
				mountPath := vol.MountPath[1:]
				if !checkExistLinkVolume(srv, vol.VolumeName) {
					_, err := AppAttachVolume(srv, vol.VolumeName, mountPath)
					if err != nil {
						log.Fatalf("Could not Link the Volume '%s' at '%s' to Application '%s': %v",
							vol.VolumeName, mountPath, srv, err)
					}
				}
				log.Printf("Volume '%s' at '%s' Linked to Application '%s'", vol.VolumeName, mountPath, srv)
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
		log.Printf("chech CreateService %s, isSrv:%v", srv, !composeConfigIsApplication(conf))
		if !composeConfigIsApplication(conf) {
			if serviceIsExist(srv) {
				log.Printf("Service '%s' is Exist", srv)
				continue
			}

			product := getComposeProductName(conf)
			plan := getComposeProductPlan(conf)
			variables := convConfigEnvVarToMap(conf)
			log.Printf("Service %s, product=%s, plan=%s, var=%v", srv, product, plan, variables)
			res, err := ServiceCreate(product, srv, plan, variables)
			if err != nil {
				log.Printf("Service '%s' Could not Create, Err: %v", srv, err)
				return err
			}
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
		log.Printf("chech CreateApplication %s, isApp:%v", srv, composeConfigIsApplication(conf))
		if composeConfigIsApplication(conf) {
			if appIsExist(srv) {
				log.Printf("Application %s is exist", srv)
				continue
			}

			plan := getComposeProductPlan(conf)
			endpointProtocol := getEndPointProtocol(conf)
			port := uint64(conf.Port[0].ContainerPort)
			minScale := getReplicas(conf)
			log.Printf("Application %s: img=%s, endpointProtocol=%s, plan=%s, minScale=%d", srv, conf.Image, endpointProtocol, plan, minScale)

			// TODO set otherValues
			res, err := ApplicationCreate(srv, conf.Image, plan, endpointProtocol, port, minScale, nil)
			if err != nil {
				log.Printf("Application '%s' Could not Create, Err: %v", srv, err)
				return err
			}
			uiApplicationStatus(res)
		}
	}
	return nil
}

func composeStart(cmd *cobra.Command, args []string) {
	log.Printf("Compose Files: %v", flagComposeFilesArray)

	configs := composeLoader(flagComposeFilesArray)

	// check metadata befor start
	checkComposeFormat(configs)

	// first create Service
	err := composeCreateService(configs)
	uiCheckErr("Could not Create Service", err)

	// then create Application, add Link the (Volumes, Domains and Env)
	err = composeAppCreate(configs)
	uiCheckErr("Could not Create App", err)
	composeAppLinkVolume(configs)
	composeAppLinkDomain(configs)
	composeAppLinkEnv(configs)
}

// TODO
func composeStop(cmd *cobra.Command, args []string) {}
