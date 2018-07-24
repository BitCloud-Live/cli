package cmd

import (
	"bytes"
	"encoding/json"
	"log"

	uvApi "github.com/uvcloud/uv-api-go/proto"
)

func uiList(list *uvApi.ListStatusRes) {
	log.Printf("Count: %d,\t Next: %d,\t  Previous:%d \r\n", list.Count, list.Next, list.Previous)
	for i, v := range list.Names {
		log.Printf("%d. %s\r\n", i, v)
	}
}

func uiServicStatus(srv *uvApi.SrvStatusRes) {
	log.Printf("Service Name: %s \r\n", srv.Name)
	log.Printf("Plan Name: %s \r\n", srv.Plan)
	log.Printf("State: %v \r\n", srv.State.String())
	log.Printf("Created: %v,\t Updated: %v \r\n", srv.Created, srv.Updated)
	uiMap(srv.Variable, "Variable")
}

func uiPortforward(pf *uvApi.PortforwardRes) {
	log.Printf("SVC: %s \r\n", pf.Svc)
	log.Printf("Port: %v \r\n", pf.Port)
	log.Printf("Token: %v \r\n", pf.Token)
}

func uiPlan(plan []*uvApi.Plan) {
	log.Println("Plan:")
	for i, p := range plan {
		log.Printf("%d. \r\n", i)
		log.Printf("\t name: %s \r\n", p.Name)
		log.Printf("\t price: %v, off: %v \r\n", p.Price, p.Off)
		log.Printf("\t Description: %v \r\n", p.Description)
	}
}
func uiMap(mapVar map[string]string, name string) {
	log.Printf("%s:", name)
	for k, v := range mapVar {
		log.Printf("\t %s: %s \r\n", k, v)
	}
}

func uiProduct(prd *uvApi.ProductRes) {
	log.Printf("Product Name: %s \r\n", prd.Name)
	log.Printf("Description: %s \r\n", prd.Description)
	uiPlan(prd.Plan)
	uiMap(prd.VariableHints, "Variable Hints")
}

func uiSettingByDetail(set *uvApi.SettingRes) {
	log.Printf("Setting Name: %s \r\n", set.Name)
	log.Printf("Application: %s Path: %s\r\n", set.App, set.Path)
	log.Printf("value: \r\n%s", set.File)
}

func uiSetting(set *uvApi.SettingRes) {
	log.Println(set.File)
}

func uiApplicationStatus(app *uvApi.AppStatusRes) {
	log.Printf("Service Name: %s \r\n", app.Name)
	log.Printf("State: %v \r\n", app.State)
	log.Printf("Config: %v \r\n", app.Config)
	log.Printf("Created: %v,\t Updated: %v \r\n", app.Created, app.Updated)
	log.Printf("VCAP_SERVICES: \r\n%v\r\n", jsonPrettyPrint(app.VcapServices))
	uiMap(app.EnvironmentVariables, "Environment variables")
}

func uiDomainStatus(dom *uvApi.DomainStatusRes) {
	log.Printf("Domain Name: %s \r\n", dom.Name)
	log.Printf("Created: %v ,Update: %v\r\n", dom.Created, dom.Updated)
	log.Printf("Address: %s \r\n", dom.Domain)
	log.Printf("AttachedTo: %s \r\n", dom.AttachedTo)
	log.Printf("TLS: %s \r\n", dom.Tls)
}

func uiVolumeSpec(vol *uvApi.VolumeSpec) {
	log.Printf("Volume Spec Name: %s \r\n", vol.Name)
	log.Printf("\t Spec Class: %s \r\n", vol.Class)
	log.Printf("\t Spec Size: %v \r\n", vol.Size)
}

func uiVolumeStatus(vol *uvApi.VolumeStatusRes) {
	log.Printf("Volume Name: %s \r\n", vol.Name)
	log.Printf("Created: %v ,Update: %v\r\n", vol.Created, vol.Updated)
	log.Printf("AttachedTo: %s \r\n", vol.AttachedTo)
	uiVolumeSpec(vol.Spec)
}

func uiCheckErr(info string, err error) {
	if err != nil {
		log.Fatalf("Could not Start the Service: %v", err)
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
