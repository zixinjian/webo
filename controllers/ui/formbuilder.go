package ui

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/lang"
)

type FormBuilder struct {
}

var textFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
			</div>
		</div>
    	`
var moneyFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
			</div>
		</div>
    	`
var textareaFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<textarea class="form-control" rows="3" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off">%s</textarea>
			</div>
		</div>
    	`
var datetimeFormat = `    <div class="form-group">
        	<label class="col-sm-3 control-label">%s</label>
        	<div class="col-sm-6">
            	<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
        	</div>
    	</div>
    	`
var dateFormate = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control datetimepicker" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
			</div>
		</div>
    	`
var passwordFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="password" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
			</div>
		</div>
    	`
var selectFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<select class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s">
				%s
				</select>
			</div>
		</div>
    	`

func BuildAddForm(oItemDef itemDef.ItemDef) string {
	var form string
	for _, field := range oItemDef.Fields {
		form = form + createFromGroup(field, field.Default)
	}
	return form
}

var initDatePickerFormat = `$(function(){
        $("#%s").datetimepicker({%sformat:'Y.m.d',lang:'zh',%s})
    });`

func BuildAddOnLoadJs(oItemDef itemDef.ItemDef) string {
	OnLoadJs := ""
	for _, field := range oItemDef.Fields {
		switch field.Input {
		case "date":
			defaultDate := ""
			if strings.EqualFold(field.Default.(string), "curtime") {
				defaultDate = "value:new Date()"
			}
			OnLoadJs = OnLoadJs + fmt.Sprintf(initDatePickerFormat, "timepicker:false,", field.Name, defaultDate)
		}
	}
	return "<script>" + OnLoadJs + "</script>"
}

func BuildUpdatedForm(oItemDef itemDef.ItemDef, oldValueMap map[string]interface{}) string {
	sn, ok := oldValueMap["sn"]
	if !ok {
		beego.Error("BuildUPdatedFrom: param sn is none")
	}
	form := fmt.Sprintf(`<input type="hidden" name="sn" value="%s"></input>`, sn)
	for _, field := range oItemDef.Fields {
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
	case "textarea":
		fromGroup = fmt.Sprintf(textareaFormat, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "text":
		//fmt.Println("text", value)
		fromGroup = fmt.Sprintf(textFormat, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "date", "datetime":
		//fmt.Println("text", value)
		fromGroup = fmt.Sprintf(dateFormate, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "password":
		fromGroup = fmt.Sprintf(passwordFormat, field.Label, field.Require, field.Label, field.Name, field.Name, "*****")
	case "select":
		var options string
		for _, option := range field.Enum {
			if option == value {
				options = options + fmt.Sprintf(`<option value="%s" selected>%s</option>`, option, lang.GetLabel(option))
				continue
			}
			options = options + fmt.Sprintf(`<option value="%s" ％s>%s</option>`, option, lang.GetLabel(option))
		}
		fromGroup = fmt.Sprintf(selectFormat, field.Label, field.Require, field.Label, field.Name, field.Name, field.Default, options)
	case "none":
		fromGroup = ""
	default:
		panic(fmt.Sprintf("createFormGroup input type: %s not support ", field.Input))
		fromGroup = ""
	}
	return fromGroup
}
