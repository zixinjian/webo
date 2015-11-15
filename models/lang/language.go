package lang

import (
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var langMap_cn = make(map[string]string)
func readFromCsv(filePath string) {
	cntb, err := ioutil.ReadFile(filePath)
	//	fmt.Println("read Ui file: ", fileName)
	//beego.Error("ok ,this is error", filePath)
	if err != nil {
		panic(err)
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	rows, err := r2.ReadAll()
	if err != nil {
		beego.Error(fmt.Sprintf("Read lang from %s error", filePath))
		return
	}
	for idx, row := range rows {
		key := strings.TrimSpace(row[0])
		if key == "" {
			beego.Error(fmt.Sprintf("Read lang key from %s error row %d is null", filePath, idx))
			continue
		}
		value := strings.TrimSpace(row[1])
		if value == "" {
			beego.Error(fmt.Sprintf("Read lang value from %s error row %d is null", filePath, idx))
		}
		langMap_cn[key] = value
	}
}
func AddLabel(key, label string){
	langMap_cn[key]=label
}
func GetLabel(key string) string {
	if v, ok := langMap_cn[key]; ok {
		return v
	}
	beego.Error(fmt.Sprintf("GetLang error key:%s", key))
	return key
}
func init() {
	filepath.Walk("conf/lang", func(filePath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filePath, ".csv") {
			readFromCsv(filePath)
		}
		return nil
	})
}

