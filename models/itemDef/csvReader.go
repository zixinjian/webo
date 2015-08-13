package itemDef

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReadDefFromCsv() map[string]ItemDef {
	var lEntityDefMap = make(map[string]ItemDef)
	filepath.Walk("conf/item", func(filePath string, f os.FileInfo, err error) error {
		if strings.HasSuffix(filePath, ".csv") {
			oItemDef := readCsv(filePath)
			lEntityDefMap[oItemDef.Name] = oItemDef
		}
		return nil
	})
	//	fmt.Println("entity", lEntityDefMap)
	return lEntityDefMap
}
func readCsv(fileName string) ItemDef {
	cntb, err := ioutil.ReadFile(fileName)
	fmt.Println("read file: ", fileName)
	if err != nil {
		panic(err)
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	rows, _ := r2.ReadAll()
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

		typo := strings.Trim(row[2], " ")
		switch typo {
		case "string", "int":
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
		case "text", "select", "password", "date", "time", "none":
			field.Input = input
		default:
			panic(fmt.Sprintf("File:%s, row:%d input :[%s] is not vaild", fileName, ridx, input))
		}
		model := strings.Trim(row[7], " ")
		switch model {
		case "":
			field.Model = "text"
		case "sn", "text", "password", "enum", "curtime", "curuser", "time", "date", "int":
			field.Model = model
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
		field.Enum = getEnumValue(strings.Trim(row[9], " "), field.Type)
		oItemDef.Fields = append(oItemDef.Fields, field)
	}
	return oItemDef
}

func getEnumValue(enumStr string, t string) []string {
	if enumStr == "" {
		return nil
	}
	enumList := strings.Split(enumStr, ",")
	var retList []string
	for _, v := range enumList {
		retList = append(retList, strings.Trim(v, " "))
	}
	return retList
}
