package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/pkg/browser"
	ybApi "github.com/yottab/proto-api/proto"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	//For windows support
)

func printCount(size int) {
	fmt.Println(secTxtColor("# Count: %d ", size))
}

func uiList(list interface{}) {
	switch list.(type) {
	case *ybApi.ActivityListRes:
		itemList := list.(*ybApi.ActivityListRes)
		printCount(len(itemList.Activities))
		for _, v := range itemList.Activities {
			printTitle("", v.Name)
			printKeyVal(whiteSpace, "Updated", toTime(v.Time))
			printKeyVal(whiteSpace, "Description", v.Description)
			printKeyVal(whiteSpace, "Attachment", v.Attachment)
			printKeyVal(whiteSpace, "Tag", v.Tag.String())
			printKeyVal(whiteSpace, "Type", v.Type)
			printKeyVal(whiteSpace, "Email", v.Email)
		}
		return
	case *ybApi.SrvListRes:
		itemList := list.(*ybApi.SrvListRes)
		printCount(len(itemList.Services))
		for _, v := range itemList.Services {
			cond := fmt.Sprintf("%s (%s)", v.Condition.GetCondition(), v.Condition.GetReason())
			printTitleByStatus("", v.Name, v.Condition.GetCondition())
			printKeyVal(whiteSpace, "Condition", cond)
			printKeyVal(whiteSpace, "Plan", v.Plan)
			printKeyVal(whiteSpace, "Service_Refrence", v.ServiceRefrence)
			printKeyVal(whiteSpace, "Updated", toTime(v.Updated))
			mainTxtPrintln(whiteSpace)
		}
		return
	case *ybApi.PrdListRes:
		itemList := list.(*ybApi.PrdListRes)
		printCount(len(itemList.Rows))
		for _, v := range itemList.Rows {
			printTitle("", v.GetName())
			printKeyVal(whiteSpace, "Description", v.Description)
			mainTxtPrintln(whiteSpace)
		}
		return
	case *ybApi.ImgListRes:
		itemList := list.(*ybApi.ImgListRes)
		printCount(len(itemList.Imgs))
		for _, v := range itemList.Imgs {
			printTitle("", v.GetName())
			printKeyVal(whiteSpace, "Created", toTime(v.Created))
			printKeyVal(whiteSpace, "Updated", toTime(v.Updated))
			mainTxtPrintln(whiteSpace)
		}
		return
	case *ybApi.VolumeSpecListRes:
		itemList := list.(*ybApi.VolumeSpecListRes)
		printCount(len(itemList.VolumeSpecs))
		for _, v := range itemList.VolumeSpecs {
			uiVolumeSpec("", v)
		}
		return
	case *ybApi.VolumeListRes:
		itemList := list.(*ybApi.VolumeListRes)
		printCount(len(itemList.Volumes))
		for _, v := range itemList.Volumes {
			uiVolumeStatus(v)
		}
		return
	case *ybApi.DomainListRes:
		itemList := list.(*ybApi.DomainListRes)
		printCount(len(itemList.Domains))
		for _, v := range itemList.Domains {
			uiDomainStatus(v)
		}
		return
	default:
		return
	}
}

func uiVolumeMount(space string, vm []*ybApi.VolumeMount) {
	if len(vm) == 0 {
		return
	}
	printKeyVal(space, "Mount", "")
	for _, m := range vm {
		printKeyVal(space+whiteSpaceDash, "To", m.GetAttachment())
		printKeyVal(space+whiteSpace, "Path", m.GetMountPath())
	}
}

func uiImageInfo(res *ybApi.ImgStatusRes) {
	printTitle("", res.GetName())
	printKeyVal(whiteSpace, "Tags", strings.Join(res.Tags, ","))
	printKeyVal(whiteSpace, "Created", toTime(res.Created))
	printKeyVal(whiteSpace, "Updated", toTime(res.Updated))
}

func uiConditions(space string, conditions []*ybApi.ServiceCondition) {
	if len(conditions) == 0 {
		return
	}
	if len(conditions) == 1 {
		c := conditions[0]
		cond := fmt.Sprintf("%s (%s)", c.GetCondition(), c.GetReason())
		printKeyVal(space, "Conditions", cond)
		return
	}
	for i, c := range conditions {
		cond := fmt.Sprintf("%s (%s)", c.GetCondition(), c.GetReason())
		printKeyVal(space+whiteSpaceDash, "Index", strconv.Itoa(i))
		printKeyVal(space+whiteSpace, "Condition", cond)
	}
}

func uiStringArray(space, title string, arr []string) {
	if len(arr) == 0 {
		return
	}

	tags := fmt.Sprintf("[%s]", strings.Join(arr, ","))
	printKeyVal(space, title, tags)
}

func uiServicStatus(srv *ybApi.SrvStatusRes) {
	printTitle("", srv.Name)
	uiConditions(whiteSpace, srv.Conditions)
	printKeyVal(whiteSpace, "Plan", srv.Plan)
	uiMapGeneralVariable(whiteSpace, srv.Variables)
	uiStringArray(whiteSpace, "Endpoints", srv.Endpoints)
	uiAttachedDomains(whiteSpace, srv.Domains)
	printKeyVal(whiteSpace, "Created", toTime(srv.Created))
	printKeyVal(whiteSpace, "Updated", toTime(srv.Updated))
}

func uiPortforward(in *ybApi.PortforwardRes) {
	bearer := string(in.Token)
	localPorts := []string{}
	for _, p := range in.Ports {
		pRemote, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
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
		colorfulPrintln(whiteSpace)
		colorfulPrintln(whiteSpace)
		colorfulPrintln(whiteSpace)
		colorfulPrintln(whiteSpace, mainTxtSprint("  Now local ports are accessible from localhost: "))
		for _, p := range localPorts {
			fports := strings.Split(p, ":")
			colorfulPrint(whiteSpace)
			fmt.Printf("%s%s%s%s\r\n",
				mainTxtColor("  Forwarding from "),
				secTxtColor(" 127.0.0.1:%s ", fports[0]),
				mainTxtBlink(" ~> "),
				secTxtColor(" %s ", fports[1]))
		}

		<-signals
		fmt.Print(secTitleColor("    Closing ports... "))
		if done != nil {
			close(done)
		}
		fmt.Println(secTitleBlink(" done  "))
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

func uiPlan(space string, plan []*ybApi.Plan) {
	if len(plan) == 0 {
		return
	}

	printKeyVal(space, "Plan", "")
	for _, p := range plan {
		printKeyVal(space+whiteSpaceDash, "Name", p.Name)
		printKeyVal(space+whiteSpace, "Description", p.Description)
		uiMapStrStr(space+whiteSpace, "Extras", p.Extras)
	}
}

func uiMapGeneralVariable(space string, mapVar map[string]*ybApi.GeneralVariable) {
	if len(mapVar) == 0 {
		return
	}

	printKeyVal(space, "General_Variables", "")
	for k, v := range mapVar {
		printKeyVal(space+whiteSpace, k, v.GetValue())
	}
}

func uiMapStrStr(space, name string, mapVar map[string]string) {
	if len(mapVar) == 0 {
		return
	}

	printKeyVal(space, name, "")
	for k, v := range mapVar {
		printKeyVal(space+whiteSpace, k, v)
	}
}

func uiRoutes(space string, mapVar []*ybApi.DomainAttachedTo) {
	if len(mapVar) == 0 {
		return
	}

	printKeyVal(space, "Variables", "")
	for _, v := range mapVar {
		printKeyVal(space+whiteSpaceDash, "To", v.GetName())
		printKeyVal(space+whiteSpace, "Path", v.GetPath())
		printKeyVal(space+whiteSpace, "EndPoint", v.GetEndpoint())
		mainTxtPrintln(space + whiteSpace)
	}
}

func uiYamlMultiLineStr(space, name, val string) {
	if val == "" {
		return
	}

	printKeyVal(space, name, "|")
	lines := strings.Split(strings.Replace(val, "\r\n", "\n", -1), "\n")
	for _, line := range lines {
		fmt.Println(
			space+whiteSpace,
			secTxtColor(line))
	}
}

func uiProduct(prd *ybApi.ProductRes) {
	printTitle("", prd.Name)
	uiYamlMultiLineStr(whiteSpace, "Description", prd.Description)
	uiPlan(whiteSpace, prd.Plan)
	uiPrdVar(whiteSpace, prd.Variables)

}

func uiPrdVar(space string, variables map[string]*ybApi.GeneralVariable) {
	if len(variables) == 0 {
		printKeyVal(space, "Variables", "[]")
	} else {
		printKeyVal(space, "Variables", "")
		for _, v := range variables {
			printKeyVal(space+whiteSpaceDash, "Name", v.Name)
			printKeyVal(space+whiteSpace, "Required", strconv.FormatBool(v.IsRequired))
			printKeyVal(space+whiteSpace, "Type", v.Type)
			printKeyVal(space+whiteSpace, "Default", v.DefaultValue)
			printKeyVal(space+whiteSpace, "Description", v.Description)
			mainTxtPrintln(space + whiteSpace)
		}
	}
}

func uiApplicationOpen(app *ybApi.AppStatusRes) {
	if len(app.Routes) == 0 {
		println(secTitleColor("No endpoint provided by the app!"))
	}
	//We only recieve one route at this moment!
	route := app.Routes[0]
	if !strings.HasPrefix("http", route) {
		println(secTitleColor("Can't open this type of endpoints right now!"))
	} else if err := browser.OpenURL(route); err != nil {
		fmt.Printf("Can't open this endpoint, error: %v!", err)
	} else {
		print("Opened in default browser!")
	}
}

func uiAppInstances(space string, instances []*ybApi.Instance) {
	if len(instances) == 0 {
		return
	}

	printKeyVal(space, "Instances", "")
	for _, v := range instances {
		printKeyVal(space+whiteSpaceDash, "ID", v.Name)
		printKeyVal(space+whiteSpace, "CPU", v.Cpu)
		printKeyVal(space+whiteSpace, "RAM", v.Ram)
		printKeyVal(space+whiteSpace, "Created", toTime(v.Created))
		mainTxtPrintln(space + whiteSpace)
	}
}

func uiApplicationStatus(app *ybApi.AppStatusRes) {
	printTitle("", app.Name)
	printKeyVal(whiteSpace, "Plan", app.Plan)
	uiConditions(whiteSpace, app.Conditions)
	uiAppInstances(whiteSpace, app.Instances)
	uiMapGeneralVariable(whiteSpace, app.Variables)
	uiMapStrStr(whiteSpace, "Environment_Variables", app.EnvironmentVariables)
	printKeyVal(whiteSpace, "Created", toTime(app.Created))
	printKeyVal(whiteSpace, "Updated", toTime(app.Updated))
	// uiAttachedDomains(whiteSpace, app.Routes)
	uiYamlMultiLineStr(whiteSpace, "VCAP_SERVICES", jsonPrettyPrint(app.VcapServices))
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func uiAttachedDomains(space string, domains []*ybApi.AttachedDomainInfo) {
	if len(domains) == 0 {
		return
	}

	printKeyVal(space, "Attached_Domains", "")
	for _, v := range domains {
		printKeyVal(space+whiteSpaceDash, "Domain", v.Domain)
		printKeyVal(space+whiteSpace, "Endpoint", v.Endpoint)
		printKeyVal(space+whiteSpace, "Type", v.EndpointType)
		mainTxtPrintln(space + whiteSpace)
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

func uiDomainStatus(v *ybApi.DomainStatusRes) {
	printTitle("", v.Domain)
	printKeyVal(whiteSpace, "TLS", v.Tls)
	uiRoutes(whiteSpace, v.AttachedTo)
	printKeyVal(whiteSpace, "Created", toTime(v.Created))
	printKeyVal(whiteSpace, "Updated", toTime(v.Updated))
	mainTxtPrintln(whiteSpace)
}

func uiVolumeSpec(space string, vol *ybApi.VolumeSpec) {
	printKeyVal(space+whiteSpaceDash, "Name", vol.Name)
	printKeyVal(space+whiteSpace, "Class", vol.Class)
	printKeyVal(space+whiteSpace, "Size", vol.Size)
}

func uiVolumeStatus(v *ybApi.VolumeStatusRes) {
	printTitle("", v.GetName())
	printKeyVal(whiteSpace, "Spec", "")
	uiVolumeSpec(whiteSpace, v.Spec)
	uiVolumeMount(whiteSpace, v.Mounts)
	printKeyVal(whiteSpace, "Created", toTime(v.Created))
	printKeyVal(whiteSpace, "Updated", toTime(v.Updated))
	mainTxtPrintln(whiteSpace)
}

func uiCheckErr(info string, err error) {
	if err != nil {
		log.Fatalf("%s, Err: %v", info, err)
	}
}
