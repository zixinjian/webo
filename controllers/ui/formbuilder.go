package ui

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/u"
)

type FormBuilder struct {
}

var textFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`

const staticFormat = `<div class="form-group">
    <label class="col-sm-3 control-label">%s</label>
    <div class="col-sm-6">
      <p class="form-control-static">%s</p>
    </div>
  </div>
`

var moneyFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<div class="input-group"><span class="input-group-addon">￥</span><input type="text" class="input-block-level form-control" data-validate="{required: %s, number:true, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/></div>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var percentFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<div class="input-group"><input type="text" class="input-block-level form-control" data-validate="{required: %s, number:true, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/><span class="input-group-addon">％</span></div>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var floatFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control" data-validate="{required: %s, number:true, messages:{required:'请输入正确的%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var textareaFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<textarea class="form-control" rows="3" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" %s>%s</textarea>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var datetimeFormat = `    <div class="form-group">
        	<label class="col-sm-3 control-label">%s</label>
        	<div class="col-sm-6">
            	<input type="text" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/>
        		<span class="help-block" id="%sHelpBlock"></span>
        	</div>
    	</div>
    	`
var dateFormate = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="text" class="input-block-level form-control datetimepicker" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var passwordFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<input type="password" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s" %s/>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var selectFormat = `    <div class="form-group">
			<label class="col-sm-3 control-label">%s</label>
			<div class="col-sm-6">
				<select class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s" %s>
				%s
				</select>
				<span class="help-block" id="%sHelpBlock"></span>
			</div>
		</div>
    	`
var uploadFormat = `    <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-6">
            <input type="file" name="%sUpload" id="%s_upload" %s/>
        </div>
    </div>
`
var autocompleteFormat = `    <div class="form-group">
            <label class="col-sm-3 control-label">%s关键字</label>
            <div class="col-sm-6">
                <input type="text" class="input-block-level form-control" id="%s_key" value="%s" %s/>
                <label>%s名称</label><input type="text" class="input-block-level form-control" readonly="true" id="%s_name" name="%s_name" data-validate="{required: %s, messages:{required:'请输入正确的%s!'}}" value="%s" placeholder="自动联想">
                <input type="hidden" id="%s" name="%s" value="%s">
            </div>
        </div>
`
var hiddenFormat = `<input type="hidden" id="%s" name="%s" value="%s">
`
var initDatePickerFormat = `
$("#%s").datetimepicker({%sformat:'Y.m.d',scrollMonth:false, lang:'zh'%s})
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
		minLength: 1,
		select: function( event, ui) {
			$( "#%s_key" ).val(ui.item.keyword);
			$( "#%s_name" ).val(ui.item.name);
			$( "#%s" ).val(ui.item.sn);
			return false;
		},
		change: function( event, ui ) {
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
            'swf'      : '../../lib/3rd/uploadify/uploadify.swf',
            'uploader' : '/item/upload/%s?sn=' + $("#sn").val(),
            'cancelImg': '../../lib/3rd/uploadify/uploadify-cancel.png',
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
				defaultDate = ",value:new Date()"
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

func BuildUpdateOnLoadJs(oItemDef itemDef.ItemDef) string {
	OnLoadJs := ""
	for _, field := range oItemDef.Fields {
		switch field.Input {
		case "date":
			defaultDate := ""
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
	return BuildUpdatedFormWithStatus(oItemDef, oldValueMap, make(map[string]string))
}

func BuildUpdatedFormWithStatus(oItemDef itemDef.ItemDef, oldValueMap map[string]interface{}, statusMap map[string]string) string {
	sn, ok := oldValueMap[s.Sn]
	if !ok {
		beego.Error("BuildUPdatedFrom: param sn is none")
		return ""
	}
	form := fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn)
	for _, field := range oItemDef.Fields {
		if _, ok := oldValueMap[field.Name]; !ok {
			oldValueMap[field.Name] = field.Default
		}
	}
	for _, field := range oItemDef.Fields {
		form = form + createFromGroup(field, oldValueMap, statusMap)
	}
	return form
}

func BuildFormElement(oItemDef itemDef.ItemDef, oldValueMap map[string]interface{}, statusMap map[string]string) map[string]string {
	retMap:=make(map[string]string)
	sn, ok := oldValueMap[s.Sn]
	if !ok {
		beego.Error("BuildUPdatedFrom: param sn is none")
		return retMap
	}
	for _, field := range oItemDef.Fields {
		if _, ok := oldValueMap[field.Name]; !ok {
			oldValueMap[field.Name] = field.Default
		}
		fmt.Println(fmt.Sprintf("{{str2html .Form_%s}}", field.Name))
	}
	for _, field := range oItemDef.Fields {
		retMap[field.Name] = createFromGroup(field, oldValueMap, statusMap)
	}
	retMap[s.Sn] = fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn)
	return retMap
}

func BuildAddForm(oItemDef itemDef.ItemDef, sn string) string {
	return BuildAddFormWithStatus(oItemDef, sn, make(map[string]string))
}

//func BuildAddFormElement(oItemDef itemDef.ItemDef, sn string) map[string]string{
//	oldValueMap := make(map[string]interface{})
//	oldValueMap[s.Sn] = sn
//	return BuildFormElement(oItemDef, oldValueMap, make(map[string]string))
//}

func BuildAddFormWithStatus(oItemDef itemDef.ItemDef, sn string, statusMap map[string]string) string {
	oldValueMap := make(map[string]interface{})
	oldValueMap[s.Sn] = sn
	return BuildUpdatedFormWithStatus(oItemDef, oldValueMap, statusMap)
}

func createFromGroup(field itemDef.Field, valueMap map[string]interface{}, statusMap map[string]string) string {
	value, ok := valueMap[field.Name]
	if !ok {
		return ""
	}
	status, sok := statusMap[field.Name]
	if !sok {
		status = ""
	}
	var fromGroup string
	switch field.Input {
	case "textarea":
		fromGroup = fmt.Sprintf(textareaFormat, field.Label, field.Require, field.Label, field.Name, field.Name, status, value, field.Name)
	case "text":
		fromGroup = fmt.Sprintf(textFormat, field.Label, field.Require, field.Label, field.Name, field.Name, u.ToStr(value), status, field.Name)
	case "static":
		fromGroup = fmt.Sprintf(staticFormat, field.Label, value)
	case "money":
		fromGroup = fmt.Sprintf(moneyFormat, field.Label, field.Require, field.Label, field.Name, field.Name, u.ToStr(value), status, field.Name)
	case s.TFloat:
		fromGroup = fmt.Sprintf(floatFormat, field.Label, field.Require, field.Label, field.Name, field.Name, u.ToStr(value), status, field.Name)
	case s.Percent:
		fromGroup = fmt.Sprintf(percentFormat, field.Label, field.Require, field.Label, field.Name, field.Name, u.ToStr(value), status, field.Name)
	case "date", "datetime":
		//		fmt.Println("date", field.Name, value)
		fromGroup = fmt.Sprintf(dateFormate, field.Label, field.Require, field.Label, field.Name, field.Name, value, status, field.Name)
	case s.Password:
		fromGroup = fmt.Sprintf(passwordFormat, field.Label, field.Require, field.Label, field.Name, field.Name, "*****", status, field.Name)
	case s.Hidden:
		fromGroup = fmt.Sprintf(hiddenFormat, field.Name, field.Name, value)
	case "select":
		var options string
		for _, option := range field.Enum {
			if option.Sn == value {
				options = options + fmt.Sprintf(`<option value="%s" selected>%s</option>`, option.Sn, option.Label)
				continue
			}
			options = options + fmt.Sprintf(`<option value="%s">%s</option>`, option.Sn, option.Label)
		}
		fromGroup = fmt.Sprintf(selectFormat, field.Label, field.Require, field.Label, field.Name, field.Name, field.Default, status, options, field.Name)
	case s.Autocomplete:
		key, kok := valueMap[field.Name+s.EKey]
		if !kok {
			key = ""
		}
		name, nok := valueMap[field.Name+s.EName]
		if !nok {
			name = ""
		}
		fromGroup = fmt.Sprintf(autocompleteFormat, field.Label, field.Name, key, status,
			field.Label, field.Name, field.Name, field.Require, field.Label, name,
			field.Name, field.Name, value)
	case s.Upload:
		fromGroup = fmt.Sprintf(uploadFormat, field.Label, field.Name, field.Name, status)
	case "none":
		fromGroup = ""
	default:
		beego.Error(fmt.Sprintf("FromBuilder.createFormGroup input %s type: %s not support ", field.Name, field.Input))
		fromGroup = ""
	}
	return fromGroup
}

//func BuildSelectElement(name, label, require, status string, valueMaps []map[string]interface{}, defaultValue interface{}, valueField string, labelField string)string{
//	options := BuildSelectOptions(valueMaps, defaultValue , valueField, labelField)
//	return fmt.Sprintf(selectFormat, label, require, label, name, name, u.ToStr(defaultValue), status, options)
//}

func BuildSelectOptions(valueMaps []map[string]interface{}, defaultValue interface{}, valueField string, labelField string, addFields ...string)string{
	options :=""
	for _, valueMap := range valueMaps {
		optionValue, ok := valueMap[valueField]
		if !ok {
			continue
		}
		optionValueStr := u.ToStr(optionValue)
		optionLabel, _:= valueMap[labelField]

		optionDatas := ""
		for _, addField := range addFields{
			addData := u.GetStringValue(valueMap, addField)
			optionDatas = optionDatas + fmt.Sprintf(`data-wb-a-%s = "%s" `, addField, addData)
		}

		if optionValue == defaultValue {
			options = options + fmt.Sprintf(`<option value="%s" %s selected>%s</option>`, optionValueStr, optionDatas, optionLabel)
			continue
		}
		options = options + fmt.Sprintf(`<option value="%s" %s>%s</option>`, optionValueStr, optionDatas, optionLabel)
	}
	return options
}