package itemDef

import (
	//	"encoding/json"
	//	"fmt"
	//	"io/ioutil"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"webo/models/s"
)

type ItemDef struct {
	Name      string  `json:name`
	Fields    []Field `json:fields`
	fieldMaps map[string]Field
}

type EnumValue struct {
	Sn    string
	Name  string
	Label string
}

type Field struct {
	Name     string      `json:name`
	Type     string      `json:type`
	Label    string      `json:label`
	Input    string      `json:input`
	Require  string      `json:require`
	Unique   string      `json:unique`
	Model    string      `json:model`
	Enum     []EnumValue `json:enum`
	Range    string      `json:range`
	Default  interface{} `json:default`
	UiList   UiListStruct
	Editable bool
}
type UiListStruct struct {
	Shown    bool
	Sortable bool
	Order    string
	Visiable bool
}

func newUiList() UiListStruct {
	return UiListStruct{true, true, "", true}
}

var fieldSn = Field{"sn", "string", "编号", "none", "false", "false", "sn", nil, "", "", UiListStruct{true, false, "", false}, false}
var fieldCreater = Field{"creater", "string", "创建人", "none", "false", "false", "curuser", nil, "", "", UiListStruct{false, false, "", false}, false}
var fieldCreateTime = Field{"createtime", "time", "创建时间", "none", "false", "false", "curtime", nil, "", "", UiListStruct{false, false, "", false}, false}

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
	return ok
}

func (this *ItemDef) GetFieldModel(fieldName string) string {
	field, ok := this.fieldMaps[fieldName]
	if ok {
		return field.Model
	}
	beego.Error(fmt.Sprintf("GetFieldModel: Error filed name %s", fieldName))
	return ""
}

func (this *ItemDef) GetNeedTrans() map[string]bool {
	retMap := make(map[string]bool)
	for _, field := range this.Fields {
		if field.Model == "enum" {
			retMap[field.Name] = true
		}
	}
	return retMap
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

func (this *ItemDef) addDefaultFields() {
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
func (field *Field) IsEditable() bool {
	return field.Editable
}
func (field *Field) GetFormValue(valueString string) (interface{}, bool) {
	switch field.Type {
	case s.TypeString:
		return valueString, true
	case s.TypeInt:
		value, err := strconv.ParseInt(valueString, 10, 64)
		if err == nil {
			return value, true
		} else {
			beego.Error(fmt.Sprintf("Get field:%s varlue: %s as %s error:%s", field.Name, valueString, field.Type, err.Error()))
			return 0, false
		}
	case s.TypeFloat:
		if strings.EqualFold(strings.TrimSpace(valueString), "") {
			return 0, false
		}
		value, err := strconv.ParseFloat(valueString, 64)
		if err == nil {
			return value, true
		} else {
			beego.Error(fmt.Sprintf("Get field:%s varlue: %s as %s error:%s", field.Name, valueString, field.Type, err.Error()))
			return 0, false
		}
	default:
		beego.Error(fmt.Sprintf("Get field:%s varlue: %s as %s error, not support type", field.Name, valueString, field.Type))
		return 0, false
	}
}
func (field *Field) GetCheckedValue(input interface{}) (interface{}, bool) {
	switch field.Type {
	case s.TypeString:
		v, ok := input.(string)
		return v, ok
	case s.TypeFloat:
		v, ok := input.(float64)
		return v, ok
	case s.TypeInt:
		v, ok := input.(int64)
		return v, ok
	default:
		beego.Error("GetCheckedValue: type not support " + field.Type)
		return input.(string), false
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
			field.Default = int64(0)
		} else {
			field.Default = ""
		}
	}
}

var EntityDefMap = make(map[string]ItemDef)

func init() {
	beego.Info("Load itemDef from csv")
	lEntityDefMap := ReadDefFromCsv()
	for k, oItemDef := range lEntityDefMap {
		oItemDef.addDefaultFields()
		lEntityDefMap[k] = oItemDef
	}
	lEntityDefMap = LoadUiListFromCsv(lEntityDefMap)
	for k, oItemDef := range lEntityDefMap {
		oItemDef.initAccDate()
		EntityDefMap[k] = oItemDef
	}
	//	odefd, ok := EntityDefMap["product"]
	//	if ok {
	//		fmt.Println("ddd", odefd.Name, odefd.GetFieldMap())
	//	}
	//	fmt.Println(EntityDefMap)
}
