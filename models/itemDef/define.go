package itemDef

import (
	//	"encoding/json"
	//	"fmt"
	//	"io/ioutil"
	"fmt"
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

var fieldSn = Field{"sn", "string", "string", "none", "false", "false", "sn", nil, ""}
var fieldCreater = Field{"creater", "string", "创建人", "none", "false", "false", "curuser", nil, ""}
var fieldCreateTime = Field{"createtime", "time", "创建时间", "none", "false", "false", "curtime", nil, ""}

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
	this.fieldMaps = fieldMap
	return fieldMap
}

func (this *ItemDef) IsValidField(fieldName string) bool {
	_, ok := this.fieldMaps[fieldName]
	fmt.Println("ok", ok, this.fieldMaps)
	return ok
}

func (this *ItemDef) GetFieldMap() map[string]Field {
	fieldMap := make(map[string]Field, len(this.Fields))
	for _, field := range this.Fields {
		fieldMap[field.Name] = field
	}
	return fieldMap
}

func (this *ItemDef) GetField(filedName string) (Field, bool) {
	v, ok := this.fieldMaps[filedName]
	return v, ok
}
func (this *ItemDef) initDefault() {
	nField := len(this.Fields)
	fields := make([]Field, nField+3)
	fields[0] = fieldSn
	for idx, field := range this.Fields {
		field.initDefault()
		fields[idx+1] = field
	}
	fields[nField+1] = fieldCreater
	fields[nField+2] = fieldCreateTime
	this.Fields = fields
}
func (field *Field) GetValue(valueString string) (interface{}, bool) {
	switch field.Type {
	case "string":
		return valueString, true
	case "int":
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

func (field *Field) initDefault() {
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
	if field.Model == "" {
		field.Model = "text"
	}
	if field.Default == nil {
		if field.Type == "int" {
			field.Default = 0
		} else {
			field.Default = ""
		}
	}
}

var EntityDefMap = make(map[string]ItemDef)

func init() {
	//	fmt.Println("initItemDefMap")
	//	bytes, err := ioutil.ReadFile("conf/item.json")
	//	if err != nil {
	//		fmt.Println("ReadFile: ", err.Error())
	//	}
	//	//    fmt.Println("readFile", bytes)
	//	var lEntityDefMap = make(map[string]ItemDef)
	//	if err := json.Unmarshal(bytes, &lEntityDefMap); err != nil {
	//		fmt.Println("Unmarshal: ", err.Error())
	//	}
	lEntityDefMap := ReadDefFromCsv()
	//fmt.Println("itemde", EntityDefMap)
	for k, oItemDef := range lEntityDefMap {
		oItemDef.initDefault()
		oItemDef.initAccDate()
		EntityDefMap[k] = oItemDef
	}
	odefd, ok := EntityDefMap["user"]
	if ok {
		fmt.Println("ddd", odefd.Name, odefd.GetFieldMap())
	}
	fmt.Println(EntityDefMap)
}
