package ctrl

import (
	"encoding/xml"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"wb/lg"
)

type SvcsInjectXml struct {
	Svc []UiSvcXml `xml:"Ui>Svc"`
}
type UiSvcXml struct {
	Url    string    `xml:"Url,attr"`
	Data   []DataXml `xml:"Datas>Data"`
	Params []Param   `xml:"Params>Param"`
	Item   ItemXml
	Tpl    string `xml:"Tpl"`
}
type Param struct {
	Name string `xml:"Name,attr"`
}
type DataXml struct {
	Key  string `xml:"Key,attr"`
	Data string `xml:",innerxml"`
}
type ItemXml struct {
	Name       string
	Injections []string `xml:"Injections>Injection"`
}

var UiInjectionMap = make(map[string]UiSvcXml)
var isLoad = false

func LoadSvcConfig(path string) {
	if isLoad == true {
		lg.Critical("SvcConfig Already Loaded")
		return
	}
	lg.Info("Load SvcConfig from:", path)
	filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filePath, ".xml") {
			readSvcConfigFromXml(filePath)
		}
		return nil
	})
	beego.Router("/sui/*", &SvcController{}, "*:UiSvc")
	isLoad = true
}
func readSvcConfigFromXml(path string) {
	lg.Info("readSvcConfigFromXml: Read svcInjectXml from:", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var svcInjectXml SvcsInjectXml
	err = xml.Unmarshal(content, &svcInjectXml)
	if err != nil {
		panic(err)
	}
	for _, uiSvcXml := range svcInjectXml.Svc {
		uiSvcXml.Url = strings.TrimSpace(uiSvcXml.Url)
		uiSvcXml.Tpl = strings.TrimSpace(uiSvcXml.Tpl)
		lstData := make([]DataXml, 0)
		for _, data := range uiSvcXml.Data {
			if strings.TrimSpace(data.Key) != "" {
				lstData = append(lstData, data)
			}
		}
		uiSvcXml.Data = lstData
		UiInjectionMap[uiSvcXml.Url] = uiSvcXml
	}
	lg.Info("UiInjectionMap", UiInjectionMap)
}
