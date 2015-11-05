package itemDef

import (
	"encoding/csv"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"webo/models/lang"
	"webo/models/s"
)

func ReadDefFromCsv() map[string]ItemDef {
	var lEntityDefMap = make(map[string]ItemDef)
	filepath.Walk("conf/item/fields", func(filePath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filePath, ".csv") {
			oItemDef := readItemDefCsv(filePath)
//			fmt.Println(oItemDef)
			lEntityDefMap[oItemDef.Name] = oItemDef
		}
		return nil
	})
	return lEntityDefMap
}
func LoadUiListFromCsv(lEntityDefMap map[string]ItemDef) map[string]ItemDef {
	filepath.Walk("conf/item/uilist", func(filePath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filePath, ".csv") {
			itemName := strings.TrimSuffix(f.Name(), ".csv")
			oItemDef, ok := lEntityDefMap[itemName]
			if !ok {
				return nil
			}
			uiListMap := readUiListCsv(filePath)
			for idx, field := range oItemDef.Fields {
				if u, ok := uiListMap[field.Name]; ok {
					field.UiList = u
					oItemDef.Fields[idx] = field
				}
			}
			lEntityDefMap[itemName] = oItemDef
		}
		return nil
	})
	return lEntityDefMap
}

func readUiListCsv(fileName string) map[string]UiListStruct {
	cntb, err := ioutil.ReadFile(fileName)
	beego.Debug(fmt.Sprintf("start parse uiList File:%s", fileName))
	if err != nil {
		panic(err)
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	rows, _ := r2.ReadAll()
	if len(rows) < 2 {
		panic(fmt.Sprintf("File:%s rows < 2", fileName))
	}
	retMap := make(map[string]UiListStruct)

	for ridx, row := range rows {
		if ridx < 1 {
			continue
		}
		filedName := strings.TrimSpace(row[0])
		if filedName == "" {
			continue
		}
		UiList := newUiList()

		if showInList := strings.TrimSpace(row[1]); strings.EqualFold(showInList, "false") {
			UiList.Shown = false
		} else {
			UiList.Shown = true
		}

		if sorable := strings.TrimSpace(row[2]); strings.EqualFold(sorable, "false") {
			UiList.Sortable = false
		} else {
			UiList.Sortable = true
		}

		if order := strings.TrimSpace(row[3]); strings.EqualFold(order, "desc") {
			UiList.Order = "desc"
		} else {
			UiList.Order = "asc"
		}
		if visiable := strings.TrimSpace(row[4]); strings.EqualFold(visiable, "false") {
			UiList.Visiable = false
		} else {
			UiList.Visiable = true
		}
		retMap[filedName] = UiList
	}
	return retMap
}

func readItemDefCsv(fileName string) ItemDef {
	cntb, err := ioutil.ReadFile(fileName)
	beego.Debug(fmt.Sprintf("start parse item File:%s", fileName))
	if err != nil {
		panic(err)
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	rows, _ := r2.ReadAll()
	//	fmt.Println("file rows", rows)
	if len(rows) < 3 {
		panic(fmt.Sprintf("File:%s rows < 3", fileName))
	}
	nameRow := rows[0]
	var itemName string
	if strings.Trim(nameRow[0], "") == "item" {
		itemName = strings.Trim(nameRow[1], " ")
	} else {
		panic(fmt.Sprintf("File:%s row0 not item", fileName))
	}

	if itemName == "" {
		panic(fmt.Sprintf("File:%s, row:%d itemName is none", fileName, 0))
	}
	oItemDef := ItemDef{}
	oItemDef.Name = itemName
	for ridx, row := range rows {
		if ridx < 2 {
			continue
		}
		if strings.Trim(row[0], " ") != "field" {
			fmt.Sprintf("File:%s, row:%d not field", fileName, ridx)
		}
		field := Field{}
		name := strings.Trim(row[1], " ")
		if name == "" {
			panic(fmt.Sprintf("File:%s, row:%d name is none", fileName, ridx))
		}
		field.Name = name

		typo := strings.TrimSpace(row[2])
		switch typo {
		case "string", "int", "float":
			field.Type = typo
		default:
			panic(fmt.Sprintf("File:%s, row:%d type :[%s] is not vaild", fileName, ridx, typo))
		}

		field.Label = strings.Trim(row[3], " ")
		require := strings.Trim(row[4], " ")
		if strings.EqualFold(require, "true") {
			field.Require = "true"
		} else {
			field.Require = "false"
		}

		unique := strings.Trim(row[5], " ")
		if strings.EqualFold(unique, "true") {
			field.Unique = "true"
		} else {
			field.Unique = "false"
		}

		input := strings.Trim(row[6], " ")
		switch input {
		case "":
			field.Input = "text"
		case "text", "select", "password", "date", "time", "none", "textarea", "datetime", "money", s.Upload, s.TFloat, s.Percent:
			field.Input = input
		case s.Autocomplete:
			field.Input = input
			field.Range = strings.TrimSpace(row[9])
		default:
			panic(fmt.Sprintf("File:%s, row:%d input :[%s] is not vaild", fileName, ridx, input))
		}
		model := strings.Trim(row[7], " ")
		switch model {
		case "":
			field.Model = "text"
		case "sn", "text", "password", "curtime", "curuser", "time", "date", "int", s.Upload:
			field.Model = model
		case s.Enum:
			field.Model = model
			field.Enum = getEnumValue(strings.Trim(row[9], " "), field.Type)
		default:
			panic(fmt.Sprintf("File:%s, row:%d model :[%s] is not vaild", fileName, ridx, model))
		}
		defaulto := strings.Trim(row[8], " ")
		switch field.Type {
		case "string":
			field.Default = defaulto
		case "int":
			if defaulto == "" {
				field.Default = nil
			}
		}

		field.UiList = newUiList()
		oItemDef.Fields = append(oItemDef.Fields, field)
	}
	//	fmt.Println("oItemDef", oItemDef)
	return oItemDef
}

func getEnumValue(enumStr string, t string) []EnumValue {
	if enumStr == "" {
		return nil
	}
	enumList := strings.Split(enumStr, ",")
	var retList []EnumValue
	for _, v := range enumList {
		value := strings.TrimSpace(v)
		retList = append(retList, EnumValue{value, value, lang.GetLabel(value)})
	}
	return retList
}
