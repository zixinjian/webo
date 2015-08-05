package itemDef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Field struct {
	Name    string      `json:name`
	Type    string      `json:type`
	Label   string      `json:label`
	Input   string      `json:input`
	Require string      `json:require`
	Unique  string      `json:unique`
	Model   string      `json:model`
	Enum    []string    `json:enum`
	Default interface{} `json:default`
}

var fieldSn = Field{"sn", "string", "sn", "none", "false", "false", "text"}
var fieldCreater = Field{"creater", "string", "创建人", "none", "false", "false", "text"}
var fieldCreateTime = Field{"createtime", "time", "创建时间", "none", "false", "false", "curtime"}

type ItemDef struct {
	Name      string  `json:name`
	Fields    []Field `json:fields`
	fieldMaps map[string]Field
}

func (this *ItemDef) initAccDate() map[string]Field {
	fieldMap := make(map[string]Field, len(this.Fields))
	for _, field := range this.Fields {
		fieldMap[field.Name] = field
	}
	return fieldMap
}

func (this *ItemDef) IsValidField(fieldName string) bool {
	_, ok := this.fieldMaps[fieldName]
	return ok
}

func (this *ItemDef) GetFieldMap() map[string]Field {
	fieldMap := make(map[string]Field, len(this.Fields))
	for _, field := range this.Fields {
		fieldMap[field.Name] = field
	}
	return fieldMap
}

func (field *Field) GetValue(valueString string) (interface{}, bool) {
	switch field.Type {
	case "string":
		return valueString, true
	case "int64":
		value, err := strconv.ParseInt(valueString, 10, 64)
		if err != nil {
			return value, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

var EntityDefMap = make(map[string]ItemDef)

func initDefault(oItemDef ItemDef) {
	nField := len(oItemDef.Fields)
	fields := make([]Field, nField+3)
	fields[0] = fieldSn
	for idx, field := range oItemDef.Fields {
		if field.Type == "" {
			field.Type = "string"
		}
		if field.Model == "" {
			field.Model = "text"
		}
		if field.Input == "" {
			field.Input = "text"
		}
		if field.Require == "" {
			field.Require = "false"
		}
		if field.Unique == "" {
			field.Unique = "false"
		}
		fields[idx+1] = field
	}
	fields[nField-2] = fieldCreater
	fields[nField-1] = fieldCreateTime
}

func init() {
	//fmt.Println("initItemDefMap")
	bytes, err := ioutil.ReadFile("conf/item.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
	}
	//    fmt.Println("readFile", bytes)
	if err := json.Unmarshal(bytes, &EntityDefMap); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
	}
	//fmt.Println("itemde", EntityDefMap)
	for _, oDef := range EntityDefMap {
		oDef.initAccDate()

		//		odefd, ok := EntityDefMap["user"]
		//		if ok {
		//			for idx, v := range odefd.Fields {
		//				fmt.Println(v.Name, idx, v.Model)
		//			}
		//		}
	}
	//	fmt.Println(EntityDefMap)
}
