package ui

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/util"
)

type FormBuilder struct {
}

var textFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
			</div>
		</div>
    	`
var moneyFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, number:true, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s"/>
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
var uploadFormat = `    <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-6">
            <input type="file" name="%sUpload" id="%s_upload" />
        </div>
    </div>
`
var autocompleteFormat = `    <div class="form-group">
            <label class="col-sm-3 control-label">%s关键字</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" id="%s_key" value="%s"/>
                <label>%s名称</label><input type="text" class="input-block-level form-control" readonly="true" id="%s_name" name="%s_name" data-validate="{required: %s, messages:{required:'请输入正确的%s!'}}" value="%s" readonly="true">
                <input type="hidden" id="%s" name="%s" value="%s">
            </div>
        </div>
`

func BuildAddForm(oItemDef itemDef.ItemDef, sn string) string {
	form := fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn)
	for _, field := range oItemDef.Fields {
		form = form + createFromGroup(field, field.Default, "", "")
	}
	return form
}

var initDatePickerFormat = `
$("#%s").datetimepicker({%sformat:'Y.m.d',lang:'zh',%s})
`
var initAutocompleteFormat = `
	$("#%s_key").autocomplete({
		source: "%s",
		autoFocus:true,
		focus: function( event, ui ) {
			$( "#%s_key" ).val(ui.item.keyword);
			$( "#%s_name" ).val(ui.item.name);
			$( "#%s" ).val(ui.item.sn);
			return false;
		},
		minLength: 2,
		select: function( event, ui) {
			$( "#%s_key" ).val(ui.item.keyword);
			$( "#%s_name" ).val(ui.item.name);
			$( "#%s" ).val(ui.item.sn);
			return false;
		},
		change: function( event, ui ) {
			console.log("ui", ui.item)
			if(!ui.item){
				$( "#%s_name" ).val("");
				$( "#%s" ).val("");
			}
		}
	})
	.autocomplete( "instance" )._renderItem = function( ul, item ) {
		return $( "<li>" )
				.append(item.keyword + "(" + item.name + ")")
				.appendTo( ul );
	};
`
var initFileUploadJs = `$('#%s_upload').uploadify({
            'swf'      : '../../asserts/3rd/uploadify/uploadify.swf',
            'uploader' : '/item/upload/%s?sn=' + $("#sn").val(),
            'cancelImg': '../../asserts/3rd/uploadify/uploadify-cancel.png',
            'fileObjName':'uploadFile'
        });
`

func BuildAddOnLoadJs(oItemDef itemDef.ItemDef) string {
	OnLoadJs := ""
	for _, field := range oItemDef.Fields {
		switch field.Input {
		case "date":
			defaultDate := ""
			if strings.EqualFold(field.Default.(string), "curtime") {
				defaultDate = "value:new Date()"
			}
			OnLoadJs = OnLoadJs + fmt.Sprintf(initDatePickerFormat, field.Name, "timepicker:false,", defaultDate)
		case s.Autocomplete:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initAutocompleteFormat, field.Name, field.Range,
				field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name)
		case s.Upload:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initFileUploadJs, field.Name, oItemDef.Name)
		}
	}
	return "<script>$(function(){" + OnLoadJs + "});</script>\n"
}

func BuildUpdatedForm(oItemDef itemDef.ItemDef, oldValueMap map[string]interface{}) string {
	sn, ok := oldValueMap["sn"]
	if !ok {
		beego.Error("BuildUPdatedFrom: param sn is none")
	}
	form := fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn)
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
		form = form + createFromGroup(field, oldValue, "", "")
	}
	return form
}

func createFromGroup(field itemDef.Field, value interface{}, key string, vlabel string) string {
	var fromGroup string
	switch field.Input {
	case "textarea":
		fromGroup = fmt.Sprintf(textareaFormat, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "text":
		fromGroup = fmt.Sprintf(textFormat, field.Label, field.Require, field.Label, field.Name, field.Name, util.ToStr(value))
	case "money":
		fromGroup = fmt.Sprintf(moneyFormat, field.Label, field.Require, field.Label, field.Name, field.Name, util.ToStr(value))
	case "date", "datetime":
		//fmt.Println("text", value)
		fromGroup = fmt.Sprintf(dateFormate, field.Label, field.Require, field.Label, field.Name, field.Name, value)
	case "password":
		fromGroup = fmt.Sprintf(passwordFormat, field.Label, field.Require, field.Label, field.Name, field.Name, "*****")
	case "select":
		var options string
		for _, option := range field.Enum {
			if option.Sn == value {
				options = options + fmt.Sprintf(`<option value="%s" selected>%s</option>`, option.Sn, option.Label)
				continue
			}
			options = options + fmt.Sprintf(`<option value="%s">%s</option>`, option.Sn, option.Label)
		}
		fromGroup = fmt.Sprintf(selectFormat, field.Label, field.Require, field.Label, field.Name, field.Name, field.Default, options)
	case s.Autocomplete:
		fromGroup = fmt.Sprintf(autocompleteFormat, field.Label, field.Name, key,
			field.Label, field.Name, field.Name, field.Require, field.Label, vlabel,
			field.Name, field.Name, value)
	case s.Upload:
		fromGroup = fmt.Sprintf(uploadFormat, field.Label, field.Name, field.Name)
	case "none":
		fromGroup = ""
	default:
		panic(fmt.Sprintf("createFormGroup input %s type: %s not support ", field.Name, field.Input))
		fromGroup = ""
	}
	return fromGroup
}
