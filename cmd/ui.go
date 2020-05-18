package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"

	"github.com/pkg/browser"
	"github.com/sirupsen/logrus"
	ybApi "github.com/yottab/proto-api/proto"

	"google.golang.org/grpc/status"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	//For windows support
)

var log = logrus.New()

type clearTextFormatter struct{}

func (f *clearTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\r\n"), nil
}

func init() {
	/*
		//Configure logging formatter
		customFormatter := new(logrus.TextFormatter)
		customFormatter.ForceColors = false
		customFormatter.DisableTimestamp = true
		customFormatter.FullTimestamp = false
		customFormatter.DisableColors = false
		log.Formatter = customFormatter

		//Windows color support
		log.SetOutput(colorable.NewColorableStdout())
	*/

	log.SetFormatter(new(clearTextFormatter))
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
		log.Printf("# Count: %d\t", len(itemList.Activities))
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
		log.Printf("# Count: %d\t", len(itemList.Services))
		for _, v := range itemList.Services {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Condition: %s", v.Condition.GetCondition())
			log.Printf("             %s", v.Condition.GetReason())
			log.Printf("  Plan: %s", v.Plan)
			log.Printf("  Service refrence: %s", v.ServiceRefrence)
			log.Printf("  Updated: %v ", toTime(v.Updated))
		}
		return
	case *ybApi.PrdListRes:
		itemList := list.(*ybApi.PrdListRes)
		log.Printf("# Count: %d\t", len(itemList.Rows))
		for _, v := range itemList.Rows {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Description: %s", v.Description)
		}
		return
	case *ybApi.ImgListRes:
		itemList := list.(*ybApi.ImgListRes)
		log.Printf("# Count: %d\t", len(itemList.Imgs))
		for _, v := range itemList.Imgs {
			log.Printf("- Name: %s", v.Name)
			//TODO: Currently no tags will be recieved from the server
			// log.Printf("  Tags: [%s]", strings.Join(v.Tags, ","))
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	case *ybApi.VolumeSpecListRes:
		itemList := list.(*ybApi.VolumeSpecListRes)
		log.Printf("# Count: %d\t", len(itemList.VolumeSpecs))
		for _, v := range itemList.VolumeSpecs {
			log.Printf("- Name: %s", v.Name)
			log.Printf("  Class: %s", v.Class)
			log.Printf("  Size: %s", v.Size)
		}
		return
	case *ybApi.VolumeListRes:
		itemList := list.(*ybApi.VolumeListRes)
		log.Printf("# Count: %d\t", len(itemList.Volumes))
		for _, v := range itemList.Volumes {
			log.Printf("- Name: %s", v.GetName())
			log.Printf("  Spec: %s", v.Spec.GetName())
			uiVolumeMount(v.Mounts)
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	case *ybApi.DomainListRes:
		itemList := list.(*ybApi.DomainListRes)
		log.Printf("# Count: %d\t", len(itemList.Domains))
		for _, v := range itemList.Domains {
			log.Printf("- Domain Name: %s", v.Domain)
			log.Printf("  TLS: %s", v.Tls)
			uiRoutes(v.AttachedTo, "  Routes")
			log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
		}
		return
	// case *ybApi.WorkerListRes:
	// 	itemList := list.(*ybApi.WorkerListRes)
	// 	log.Printf("# Count: %d\t", len(itemList.Services))
	// 	for _, v := range itemList.Services {
	// 		log.Printf("- Name: %s", v.Name)
	// 		log.Printf("  State: %s", v.State.String())
	// 		log.Printf("  Image: %s", v.Config.Image)
	// 		log.Printf("  Port: %d", v.Config.Port)
	// 		log.Printf("  Created: %v , Updated: %v ", toTime(v.Created), toTime(v.Updated))
	// 	}
	// 	return
	default:
		return
	}
}

func uiVolumeMount(vm []*ybApi.VolumeMount){
	for _, m := range vm {
		log.Printf("  AttachedTo: %s", m.GetAttachment())
		log.Printf("  MountPath: %s", m.GetMountPath())
	}
}

func uiImageInfo(res *ybApi.ImgStatusRes) {
	log.Printf("Name: %s", res.Name)
	log.Printf("Tags: %s", strings.Join(res.Tags, ","))
	log.Printf("Created: %v , Updated: %v ", toTime(res.Created), toTime(res.Updated))
}

func uiConditions(conditions []*ybApi.ServiceCondition){
	log.Println("Conditions:")
	for i, c := range conditions {
		log.Printf(" - index: %d", i)
		log.Printf("   condition: %s", c.GetCondition())
		log.Printf("   reason:%s", c.GetReason())
	}
}

func uiServicStatus(srv *ybApi.SrvStatusRes) {
	log.Printf("Service Name: %s ", srv.Name)
	log.Printf("Plan Name: %s ", srv.Plan)
	uiConditions(srv.Conditions)
	log.Printf("Plan: %v ", srv.Plan)
	log.Printf("Created: %v , Updated: %v ", toTime(srv.Created), toTime(srv.Updated))
	uiMapGeneralVariable(srv.Variables, "Variable")
	uiStringArray("List of endpoints", srv.Endpoints)
	uiAttachedDomains(srv.Domains)
}

func uiNFSMount(in *ybApi.PortforwardRes) {
	log.Printf("FTP portforwarding is ready @ localhost:21")
	log.Printf("Now you can connect using your favorite ftp client, e.g. filezilla...")
}

func uiPortforward(in *ybApi.PortforwardRes) {
	bearer := string(in.Token)
	localPorts := []string{}
	for _, p := range in.Ports {
		pRemote, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		//FIXME: better ux needed here
		pLocal := pRemote + 2000
		pNew := fmt.Sprintf("%d:%d", pLocal, pRemote)
		localPorts = append(localPorts, pNew)
	}
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
		fmt.Println("#####################################################")
		fmt.Printf("### Forwarding ports [local:remote]: %s", localPorts)
		fmt.Println("")
		fmt.Println("### Now local ports are accessible from localhost")
		fmt.Println("### For example: localhost:3306, 127.0.0.1:5432")
		fmt.Println("#####################################################")
		<-signals
		fmt.Println("closing ports...")
		if done != nil {
			close(done)
		}
		fmt.Println("done.")
		os.Exit(1)

	}()
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, proxyURL)
	pf, err := portforward.New(dialer, localPorts, done, rdy, &stdout, &stderr)
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
		log.Printf("  Extras: %v ", p.Extras)
	}
}
func uiMapGeneralVariable(mapVar map[string]*ybApi.GeneralVariable, name string) {
	if len(mapVar) == 0 {
		log.Printf("%s: None", name)
		return
	}
	log.Printf("%s:", name)
	for k, v := range mapVar {
		log.Printf("\t %s: %s ", k, v.GetValue())
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

func uiRoutes(mapVar []*ybApi.DomainAttachedTo, name string) {
	if len(mapVar) == 0 {
		log.Printf("%s: None", name)
		return
	}
	log.Printf("%s:", name)
	for _, v := range mapVar {
		log.Printf("\t %s -> %s (%s)", v.GetName(), v.GetPath(), v.GetEndpoint())
	}
}

func uiProduct(prd *ybApi.ProductRes) {
	log.Printf("Product: %v", prd)
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
		log.Printf("  Type: %s", vari.Type)
		log.Printf("  Default: %s", vari.DefaultValue)
		log.Printf("  Description: %s", vari.Description)
		log.Printf("  Is required: %v", vari.IsRequired)
	}
}

//Deprecated
// func uiSettingByDetail(set *ybApi.SettingRes) {
// 	log.Printf("Setting Name: %s ", set.Name)
// 	log.Printf("Application: %s Path: %s", set.App, set.Path)
// 	log.Print("Value: ")
// 	log.Print(set.File)
// }

//Deprecated
// func uiSetting(set *ybApi.SettingRes) {
// 	log.Println(set.File)
// }

func uiApplicationOpen(app *ybApi.AppStatusRes) {
	if len(app.Routes) == 0 {
		print("No endpoint provided by the app!")
	}
	//We only recieve one route at this moment!
	route := app.Routes[0]
	if !strings.HasPrefix("http", route){
		print("Can't open this type of endpoints right now!")
	}else if err := browser.OpenURL(route); err != nil {
		fmt.Printf("Can't open this endpoint, error: %v!", err)
	}else {
		print("Opened in default browser!")
	}
}

func uiApplicationStatus(app *ybApi.AppStatusRes) {
	log.Printf("Service Name: %s ", app.Name)
	log.Printf("Plan: %v", app.Plan)
	uiMapGeneralVariable(app.Variables, "Variables")
	uiConditions(app.Conditions)
	log.Printf("Created: %v , Updated: %v ", toTime(app.Created), toTime(app.Updated))
	uiMap(app.EnvironmentVariables, "Environment variables")
	// uiAttachedDomains(app.Domains)
	uiVCAP(app.VcapServices)
}

func uiVCAP(vcap string){
	if vcap == "" {
		log.Printf("VCAP_SERVICES: None")
		return
	}
	log.Printf("VCAP_SERVICES: ")
	lines := strings.Split(strings.Replace(jsonPrettyPrint(vcap), "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		log.Print(line)
	}
}

//Deprecated
// func uiWorkerStatus(worker *ybApi.WorkerRes) {
// 	log.Printf("Service Name: %s ", worker.Name)
// 	log.Printf("State: %v ", worker.State)
// 	log.Printf("Config: %v ", worker.Config)
// 	log.Printf("Created: %v , Updated: %v ", toTime(worker.Created), toTime(worker.Updated))
// }

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

func uiStreamLog(client ybApi.YB_AppLogClient) error {
	var byteRecieved = 0
	for {
		c, err := client.Recv()
		if err != nil {
			if status.Code(err) == codes.OutOfRange {
				log.Printf("Transfer of %d bytes done", byteRecieved)
				return nil
			}
			return err
		}
		byteRecieved += len(c.Chunk)
		log.Printf(string(c.Chunk))
	}
}

func uiDomainStatus(dom *ybApi.DomainStatusRes) {
	log.Printf("Domain Name: %s ", dom.Domain)
	log.Printf("Created: %v , Updated: %v", toTime(dom.Created), toTime(dom.Updated))
	uiRoutes(dom.AttachedTo, "Routes")
	log.Printf("TLS: %s ", dom.Tls)
}

func uiVolumeSpec(vol *ybApi.VolumeSpec) {
	log.Printf("Volume Spec Name: %s ", vol.Name)
	log.Printf("Spec Class: %s ", vol.Class)
	log.Printf("Spec Size: %v ", vol.Size)
}

func uiVolumeStatus(vol *ybApi.VolumeStatusRes) {
	log.Printf("Volume Name: %s ", vol.Name)
	log.Printf("Created: %v , Updated: %v", toTime(vol.Created), toTime(vol.Updated))
	uiVolumeMount(vol.Mounts)
	uiVolumeSpec(vol.Spec)
}

func uiCheckErr(info string, err error) {
	if err != nil {
		log.Fatalf("%s, Err: %v", info, err)
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
