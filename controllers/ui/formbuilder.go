package ui

import (
	"fmt"
	"webo/models/itemDef"
)

type FormBuilder struct {
}

var textFormat = `    <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-6">
            <input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
        </div>
    </div>`
var passwordFormat = `    <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-6">
            <input type="password" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
        </div>
    </div>`
var selectFormat = `    <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-6">
            <select class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s">
            %s
            </select>
        </div>
    </div>`

func BuildAddForm(entity string) string {
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		fmt.Println("BuildAddForm none")
	}
	var form string
	for _, field := range oEntityDef.Fields {
		form = form + createFromGroup(field, field.Default)
	}
	return form
}

func BuildUpdatedForm(entity string, oldValueMap map[string]interface{}) string {
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		fmt.Println("BuildUpdatedForm none")
	}
	var form string
	for _, field := range oEntityDef.Fields {
		var oldValue interface{}
		//        fmt.Println("old", oldValueMap, field.Name)
		value, ok := oldValueMap[field.Name]
		if ok {
			//            fmt.Println("value", value)
			oldValue = value
		} else {
			oldValue = field.Default
		}
		form = form + createFromGroup(field, oldValue)
	}
	return form
}

func createFromGroup(field itemDef.Field, value interface{}) string {
	var fromGroup string
	switch field.Input {
	case "text":
		//        fmt.Println("text", value)
		fromGroup = fmt.Sprintf(textFormat, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "password":
		fromGroup = fmt.Sprintf(passwordFormat, field.Label, field.Require, field.Label, field.Name, field.Name, "*****")
	case "select":
		var options string
		for _, option := range field.Enum {
			if option == value {
				options = options + fmt.Sprintf(`<option value="%s" selected>%s</option>`, option, option)
			}
			options = options + fmt.Sprintf(`<option value="%s" ％s>%s</option>`, option, option)
		}
		fromGroup = fmt.Sprintf(selectFormat, field.Label, field.Require, field.Label, field.Name, field.Name, field.Default, options)
	default:
		fromGroup = ""
	}
	return fromGroup
}
