package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"google.golang.org/grpc/codes"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/browser"
	ybApi "github.com/yottab/proto-api/proto"

	"google.golang.org/grpc/status"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	//For windows support
	"github.com/mattn/go-colorable"
)

var log = logrus.New()

func init() {
	//Configure logging formatter
	customFormatter := new(logrus.TextFormatter)
	customFormatter.ForceColors = true
	customFormatter.DisableTimestamp = false
	customFormatter.FullTimestamp = false
	customFormatter.DisableColors = false
	log.Formatter = customFormatter

	//Windows color support
	log.SetOutput(colorable.NewColorableStdout())
}

func uiStringArray(title string, arr []string) {
	if len(arr) == 0 {
		log.Printf("%s: None", title)
		return
	}
	log.Printf("%s: [%s]", title, strings.Join(arr, ","))

}

func uiList(list interface{}) {
	switch list.(type) {
	case *ybApi.ActivityListRes:
		itemList := list.(*ybApi.ActivityListRes)
		log.Printf("Count: %d\t", len(itemList.Activities))
		for _, v := range itemList.Activities {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Description: %s", v.Description)
			log.Printf("  Attachment: %s", v.Attachment)
			log.Printf("  Tag: %s", v.Tag.String())
			log.Printf("  Type: %s", v.Type)
			log.Printf("  Email: %s", v.Email)
			log.Printf("  Updated: %v ", toTime(v.Time))
		}
		return
	case *ybApi.SrvListRes:
		itemList := list.(*ybApi.SrvListRes)
		log.Printf("Count: %d\t", len(itemList.Services))
		for _, v := range itemList.Services {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  State: %s", v.State.String())
			log.Printf("  Service refrence: %s", v.ServiceRefrence)
			log.Printf("  Updated: %v ", toTime(v.Updated))
		}
		return
	case *ybApi.PrdListRes:
		itemList := list.(*ybApi.PrdListRes)
		log.Printf("Count: %d\t", len(itemList.Rows))
		for _, v := range itemList.Rows {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Description: %s", v.Description)
		}
		return
	case *ybApi.ImgListRes:
		itemList := list.(*ybApi.ImgListRes)
		log.Printf("Count: %d\t", len(itemList.Imgs))
		for _, v := range itemList.Imgs {
			log.Printf("- Name: %s", v.Name)
			//TODO: Currently no tags will be recieved from the server
			// log.Printf("  Tags: [%s]", strings.Join(v.Tags, ","))
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	case *ybApi.VolumeSpecListRes:
		itemList := list.(*ybApi.VolumeSpecListRes)
		log.Printf("Count: %d\t", len(itemList.VolumeSpecs))
		for _, v := range itemList.VolumeSpecs {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Class: %s", v.Class)
			log.Printf("  Size: %s", v.Size)
		}
		return
	case *ybApi.VolumeListRes:
		itemList := list.(*ybApi.VolumeListRes)
		log.Printf("Count: %d\t", len(itemList.Volumes))
		for _, v := range itemList.Volumes {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Spec: %s", v.Spec.Name)
			log.Printf("  AttachedTo: %s", v.AttachedTo)
			log.Printf("  MountPath: %s", v.MountPath)
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	case *ybApi.DomainListRes:
		itemList := list.(*ybApi.DomainListRes)
		log.Printf("Count: %d\t", len(itemList.Domains))
		for _, v := range itemList.Domains {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Spec: %s", v.Domain)
			log.Printf("  TLS: %s", v.Tls)
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	case *ybApi.WorkerListRes:
		itemList := list.(*ybApi.WorkerListRes)
		log.Printf("Count: %d\t", len(itemList.Services))
		for _, v := range itemList.Services {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  State: %s", v.State.String())
			log.Printf("  Image: %s", v.Config.Image)
			log.Printf("  Port: %d", v.Config.Port)
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	default:
		return
	}
}

func uiImageInfo(res *ybApi.ImgStatusRes) {
	log.Printf("Name: %s", res.Name)
	log.Printf("Tags: %s", strings.Join(res.Tags, ","))
	log.Printf("Created: %v , Updated: %v ", toTime(res.Created), toTime(res.Updated))
}

func uiServicStatus(srv *ybApi.SrvStatusRes) {
	log.Printf("Service Name: %s ", srv.Name)
	log.Printf("Plan Name: %s ", srv.Plan)
	log.Printf("State: %v ", srv.State.String())
	log.Printf("Created: %v , Updated: %v ", toTime(srv.Created), toTime(srv.Updated))
	uiMap(srv.Variable, "Variable")
	uiStringArray("List of endpoints", srv.Endpoints)
	uiAttachedDomains(srv.Domains)
}

func uiNFSMount(in *ybApi.PortforwardRes) {
	log.Printf("NFS portforwarding is READY!")
	log.Printf("NFSv4 now is available at host 127.0.0.1, port 2049")
	log.Printf(`How to mount:
				On archlinux: see https://wiki.archlinux.org/index.php/NFS
				On Ubuntu: see https://help.ubuntu.com/community/SettingUpNFSHowTo#NFSv4_client
				On OSX operating systems, you can connect using cmd+k in the Finder application
				On Windows operating systems use the following steps:
				Go to Control Panel → Programs → Programs and Features
				Select: Turn Windows features on or off" from the left hand navigation.
				Scroll down to "Services for NFS" and click the "plus" on the left
				Check "Client for NFS"
				Select "Ok"
				Windows should install the client. Once the client package is install you will have the "mount" command available.
				For example you can use the following cmd command to mount to "L" drive:
				mount \127.0.0.1\ L:`)

}

func uiPortforward(in *ybApi.PortforwardRes) {
	bearer := string(in.Token)
	proxyURL, _ := url.Parse(in.ProxyHost)
	conf := &rest.Config{
		BearerToken:     bearer,
		Host:            fmt.Sprintf("%s://%s", proxyURL.Scheme, proxyURL.Host),
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	transport, upgrader, err := spdy.RoundTripperFor(conf)
	if err != nil {
		panic(err.Error())
	}

	var done = make(chan struct{}, 1)
	var rdy = make(chan struct{})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	go func() {
		fmt.Printf("Forwarding ports: %s", in.Ports)
		<-signals
		fmt.Print("closing the opened ports...")
		if done != nil {
			close(done)
		}
		os.Exit(1)
	}()
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, proxyURL)
	pf, err := portforward.New(dialer, in.Ports, done, rdy, &stdout, &stderr)
	if err != nil {
		panic(err.Error())
	}
	err = pf.ForwardPorts()
	if err != nil {
		panic(err.Error())
	}
}

func uiPlan(plan []*ybApi.Plan) {
	log.Println("Plan:")
	for _, p := range plan {
		log.Printf("- Name: %s ", p.Name)
		log.Printf("  Description: %v ", p.Description)
	}
}
func uiMap(mapVar map[string]string, name string) {
	if len(mapVar) == 0 {
		log.Printf("%s: None", name)
		return
	}
	log.Printf("%s:", name)
	for k, v := range mapVar {
		log.Printf("\t %s: %s ", k, v)
	}
}

func uiProduct(prd *ybApi.ProductRes) {
	log.Printf("Product Name: %s ", prd.Name)
	descLines := strings.Split(strings.Replace(prd.Description, "\r\n", "\n", -1), "\n")
	log.Print("Description: ")
	for _, line := range descLines {
		log.Print(line)
	}
	uiPlan(prd.Plan)
	if len(prd.Variables) == 0 {
		log.Print("Variables: []")
		return
	}
	log.Print("Variables:")
	for _, vari := range prd.Variables {
		log.Printf("- Name: %s", vari.Name)
		log.Printf("  Type: %s, default: %s", vari.Type, vari.DefaultValue)
		log.Printf("  Description: %s", vari.Description)
		log.Printf("  Is required: %v", vari.IsRequired)
	}
}

func uiSettingByDetail(set *ybApi.SettingRes) {
	log.Printf("Setting Name: %s ", set.Name)
	log.Printf("Application: %s Path: %s", set.App, set.Path)
	log.Print("Value: ")
	log.Print(set.File)
}

func uiSetting(set *ybApi.SettingRes) {
	log.Println(set.File)
}

func uiApplicationOpen(app *ybApi.AppStatusRes) {
	if len(app.Config.Routes) == 0 {
		print("No endpoint provided by the app!")
	}
	//We only recieve one route at this moment!
	route := app.Config.Routes[0]
	switch app.Config.EndpointType {
	//We only handle http endpoint at this moment!
	case "http":
		route = fmt.Sprintf("https://%s:443", route)
		if err := browser.OpenURL("https://" + route); err != nil {
			fmt.Printf("Can't open this endpoint, error: %v!", err)
		}
		print("Opened in default browser!")
		return
	default:
		print("Can't open this type of endpoints right now!")
		return
	}

}
func uiApplicationStatus(app *ybApi.AppStatusRes) {
	log.Printf("Service Name: %s ", app.Name)
	log.Printf("State: %v ", app.State)
	log.Printf("Image: %v", app.Config.Image)
	log.Printf("Internal-port: %v ", app.Config.Port)
	log.Printf("Minimum-scale: %v", app.Config.MinScale)

	//Print routes
	log.Printf("Endpoints(Public URLs):")
	for idx, route := range app.Config.Routes {
		switch app.Config.EndpointType {
		case "http":
			route = fmt.Sprintf("https://%s:443", route)
			break
		case "grpc":
			route = fmt.Sprintf("dns://%s:443", route)
			break
		default:
			break
		}
		log.Printf("\t%d. %v -> (%v endpoint type)", idx+1, route, app.Config.EndpointType)
	}
	log.Printf("Created: %v , Updated: %v ", toTime(app.Created), toTime(app.Updated))
	uiMap(app.EnvironmentVariables, "Environment variables")
	// uiAttachedDomains(app.Domains)
	if app.VcapServices == "" {
		log.Printf("VCAP_SERVICES: None")
		return
	}
	log.Printf("VCAP_SERVICES: ")
	lines := strings.Split(strings.Replace(jsonPrettyPrint(app.VcapServices), "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		log.Print(line)
	}

}

func uiWorkerStatus(worker *ybApi.WorkerRes) {
	log.Printf("Service Name: %s ", worker.Name)
	log.Printf("State: %v ", worker.State)
	log.Printf("Config: %v ", worker.Config)
	log.Printf("Created: %v , Updated: %v ", toTime(worker.Created), toTime(worker.Updated))
}

func uiAttachedDomains(domains []*ybApi.AttachedDomainInfo) {
	if len(domains) == 0 {
		log.Println("Attached domains: None")
		return
	}
	log.Println("Attached domains:")
	log.Println("Domain | Endpoint | Type")
	for i, d := range domains {
		log.Printf("%d. %s | %s | %s ", i, d.Domain, d.Endpoint, d.EndpointType)
	}
}

func uiApplicationLog(client ybApi.YB_AppLogClient) {
	var byteRecieved = 0
	for {
		c, err := client.Recv()
		if err != nil {
			if status.Code(err) == codes.OutOfRange {
				log.Printf("Transfer of %d bytes done", byteRecieved)
				return
			}
			log.Fatal(err)
		}
		byteRecieved += len(c.Chunk)
		log.Printf(string(c.Chunk))
	}
}

func uiDomainStatus(dom *ybApi.DomainStatusRes) {
	log.Printf("Domain Name: %s ", dom.Domain)
	log.Printf("Created: %v , Update: %v", toTime(dom.Created), toTime(dom.Updated))
	log.Printf("AttachedTo: %s ", dom.AttachedTo)
	log.Printf("TLS: %s ", dom.Tls)
}

func uiVolumeSpec(vol *ybApi.VolumeSpec) {
	log.Printf("Volume Spec Name: %s ", vol.Name)
	log.Printf("Spec Class: %s ", vol.Class)
	log.Printf("Spec Size: %v ", vol.Size)
}

func uiVolumeStatus(vol *ybApi.VolumeStatusRes) {
	log.Printf("Volume Name: %s ", vol.Name)
	log.Printf("Created: %v , Update: %v", toTime(vol.Created), toTime(vol.Updated))
	log.Printf("AttachedTo: %s ", vol.AttachedTo)
	uiVolumeSpec(vol.Spec)
}

func uiCheckErr(info string, err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
