package uibuilder

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"wb/cc"
	"wb/cs"
	"wb/ii"
	"wb/lg"
	"wb/ut"
)

var textareaFormat = `            <div class="form-group">
            <label class="col-sm-3 control-label">%s</label>
            <div class="col-sm-7">
                <textarea class="form-control" rows="3" class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s!'}}" name="%s" id="%s" autocomplete="off" %s>%s</textarea>
                <span class="help-block" id="%sHelpBlock"></span>
            </div>
        </div>
        `
var selectFormat = `            <div class="form-group">
            <label class="col-md-3 col-sm-3 control-label">%s</label>
            <div class="col-sm-7">
                <select class="input-block-level form-control" data-validate="{required: %s, messages:{required:'请输入%s'}}" name="%s" id="%s" autocomplete="off" value="%s" %s>
                %s
                </select>
                <span class="help-block" id="%sHelpBlock"></span>
            </div>
        </div>
        `
var uploadFormat = `            <div class="form-group">
        <label class="col-sm-3 control-label">%s</label>
        <div class="col-sm-7">
            <input type="file" multiple name="%sUpload" id="%s_upload" %s  ui-jq="fileinput" data-upload-url="%s" />
        </div>
    </div>
`
var autocompleteFormat = `            <div class="form-group">
            <label class="col-sm-3 control-label">%s关键字</label>
               <div class="col-sm-7">
                  <input type="text" class="input-block-level form-control" id="%s_key" value="%s" %s/>
                  <label>%s名称</label><input type="text" class="input-block-level form-control" readonly="true" id="%s_name" name="%s_name" data-validate="{required: %s, messages:{required:'请输入正确的%s!'}}" value="%s" placeholder="自动联想">
                  <input type="hidden" id="%s" name="%s" value="%s">
               </div>
            </div>
`
var hiddenFormat = `<input type="hidden" id="%s" name="%s" value="%s">
`
var initDatePickerFormat = `
$("#%s").datetimepicker({%sformat:'Y-m-d',scrollMonth:false, lang:'zh'%s})
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
            'swf'      : '../../lib/uploadify/uploadify/uploadify.swf',
            'uploader' : '/item/upload/%s?sn=' + $("#sn").val(),
            'cancelImg': '../../lib/uploadify/uploadify/uploadify-cancel.png',
            'fileObjName':'uploadFile'
        });
`

func BuildAddOnLoadJs(oItemInfo ii.ItemInfo) string {
	OnLoadJs := ""
	for _, field := range oItemInfo.Fields {
		switch field.Input {
		case ii.IDate, ii.IDateTime:
			defaultDate := ""
			if field.Default == "" || strings.EqualFold(field.Default, "curtime") {
				defaultDate = ",value:new Date()"
			}
			typeStr := ""
			if field.Input == ii.IDate {
				typeStr = "timepicker:false,"
			}
			if field.Input == ii.ITime {
				typeStr = "datepicker:false,"
			}
			OnLoadJs = OnLoadJs + fmt.Sprintf(initDatePickerFormat, field.Name, typeStr, defaultDate)
		case ii.IAutocomplete:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initAutocompleteFormat, field.Name, field.Src,
				field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name)
		case ii.IUpload:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initFileUploadJs, field.Name, oItemInfo.Name)
		}
	}
	return "<script>$(function(){" + OnLoadJs + "});</script>\n"
}

func BuildUpdateOnLoadJs(oItemInfo ii.ItemInfo) string {
	OnLoadJs := ""
	for _, field := range oItemInfo.Fields {
		switch field.Input {
		case ii.IAutocomplete:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initAutocompleteFormat, field.Name, field.Src,
				field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name, field.Name)
		case ii.IUpload:
			OnLoadJs = OnLoadJs + fmt.Sprintf(initFileUploadJs, field.Name, oItemInfo.Name)
		}
	}
	return "<script>$(function(){" + OnLoadJs + "});</script>\n"
}

func BuildUpdatedForm(oItemInfo ii.ItemInfo, oldValueMap map[string]interface{}) string {
	return BuildUpdatedFormWithStatus(oItemInfo, oldValueMap, make(map[string]string))
}

const rowStart = `      <div class="row">
`
const rowEnd = `      </div>
`
const helfColStart = `        <div class="col-md-6">
`
const helfColEnd = `        </div>
`

func BuildUpdatedFormWithStatus(oItemDef ii.ItemInfo, oldValueMap map[string]interface{}, statusMap map[string]string) string {
	sn, ok := oldValueMap[cc.Col.Sn]
	if !ok {
		lg.Error("BuildUPdatedFrom: param sn is none")
		return ""
	}
	form := fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn) + "\n"
	for _, field := range oItemDef.Fields {
		if _, ok := oldValueMap[field.Name]; !ok {
			oldValueMap[field.Name] = field.Default
		}
	}
	i := 0
	for _, field := range oItemDef.Fields {
		input := createFromGroup(oItemDef.Name, field, oldValueMap, statusMap)
		if input == "" {
			continue
		}
		if i%2 == 0 {
			form = form + rowStart
		} else if field.Input == ii.ITextArea {
			form = form + rowEnd + rowStart
			i = i + 1
		}
		form = form + helfColStart + input + helfColEnd
		if field.Input == ii.ITextArea {
			i = i + 2
		} else {
			i = i + 1
		}
		if i%2 == 0 {
			form = form + rowEnd
		}
	}
	return form
}

func BuildFormElement(oItemInfo ii.ItemInfo, oldValueMap map[string]interface{}, statusMap map[string]string) map[string]string {
	retMap := make(map[string]string)
	sn, ok := oldValueMap[cc.Col.Sn]
	//lg.Debug("BuildFormElement oldvalue", oldValueMap)
	if ok {
		retMap[cc.Col.Sn] = fmt.Sprintf(`<input type="hidden" id="sn" name="sn" value="%s">`, sn)
	}
	for _, field := range oItemInfo.Fields {
		if _, ok := oldValueMap[field.Name]; !ok {
			oldValueMap[field.Name] = field.Default
		}
	}
	for _, field := range oItemInfo.Fields {
		if field.Name == cc.Col.Sn {
			continue
		}
		retMap[field.Name] = createFromGroup(oItemInfo.Name, field, oldValueMap, statusMap)
	}
	return retMap
}

func BuildAddForm(oItemDef ii.ItemInfo, sn string) string {
	return BuildAddFormWithStatus(oItemDef, sn, make(map[string]string))
}

func BuildAddFormWithStatus(oItemDef ii.ItemInfo, sn string, statusMap map[string]string) string {
	oldValueMap := make(map[string]interface{})
	oldValueMap[cc.Col.Sn] = sn
	return BuildUpdatedFormWithStatus(oItemDef, oldValueMap, statusMap)
}

const (
	StatusStatic = "static"
)

func createFromGroup(item string, field ii.FieldInfo, valueMap map[string]interface{}, statusMap map[string]string) string {
	value, ok := valueMap[field.Name]
	if !ok {
		return ""
	}
	status, sok := statusMap[field.Name]
	if !sok {
		status = ""
	}
	if status == StatusStatic {
		return createStaticInput(field, value, status)
	}
	var fromGroup string
	switch field.Input {
	case ii.ITextArea:
		fromGroup = fmt.Sprintf(textareaFormat, field.Label, field.Require, field.Label, field.Name, field.Name, status, value, field.Name)
	case ii.IText, ii.IFloat, ii.IInt, ii.IPercent, ii.IDate, ii.IDateTime, ii.IMoney, ii.IPassword:
		fromGroup = createTextInput(field, value, status)
	case ii.IStatic:
		fromGroup = createStaticInput(field, value, status)
	case ii.IHidden:
		fromGroup = fmt.Sprintf(hiddenFormat, field.Name, field.Name, value)
	case ii.ISelect:
		var options string
		for _, option := range field.Enums {
			if option.Key == value {
				options = options + fmt.Sprintf(`<option value="%s" selected>%s</option>`, option.Key, option.Label)
				continue
			}
			options = options + fmt.Sprintf(`<option value="%s">%s</option>`, option.Key, option.Label)
		}
		fromGroup = fmt.Sprintf(selectFormat, field.Label, field.Require, field.Label, field.Name, field.Name, field.Default, status, options, field.Name)
	case ii.IAutocomplete:
		key, kok := valueMap[field.Name+cc.EKey]
		if !kok {
			key = ""
		}
		name, nok := valueMap[field.Name+cc.EName]
		if !nok {
			name = ""
		}
		fromGroup = fmt.Sprintf(autocompleteFormat, field.Label, field.Name, key, status,
			field.Label, field.Name, field.Name, field.Require, field.Label, name,
			field.Name, field.Name, value)
	case ii.IUpload:
		sn := ut.Maps.GetString(valueMap, "sn")
		uploadUrl := fmt.Sprintf("/item/upload/%s?sn=%s", item, sn)
		fromGroup = fmt.Sprintf(uploadFormat, field.Label, field.Name, field.Name, status, uploadUrl)
	case ii.INone:
		fromGroup = ""
	default:
		lg.Error(fmt.Sprintf("FromBuilder.createFormGroup input %s type: %s not support ", field.Name, field.Input))
		fromGroup = ""
	}
	return fromGroup
}

const staticFormat = `          <div class="form-group">
    <label class="col-sm-3 control-label">%s</label>
    <div class="col-sm-8">
      <p class="form-control-static">%s</p>
    </div>
  </div>
`

func createStaticInput(field ii.FieldInfo, value interface{}, status string) string {
	return fmt.Sprintf(staticFormat, field.Label, ut.ToStr(value)+field.Unit)
}

const addonFormat = `<span class="input-group-addon">%s</span>`
const uiJqDateFormat = `ui-jq="datetimepicker" ui-options = "{timepicker:false,format:'Y-m-d',lang:'zh',%s scrollMonth:false}"`

const inputGroupStart = `<div class="input-group">
`
const inputGroupEnd = `
</div>`

func createTextInput(field ii.FieldInfo, value interface{}, status string) string {
	attrMap := make(map[string]string)
	attrMap["Name"] = field.Name
	attrMap["Label"] = field.Label
	if field.Require == "true" {
		attrMap["RequireFlag"] = "<span class=\"wb-require-star\">*</span>"
	} else {
		attrMap["RequireFlag"] = ""
	}
	attrMap["Value"] = ut.ToStr(field.FormatValue(value))
	attrMap["Require"] = field.Require
	attrMap["Desc"] = ""
	attrMap["Status"] = status
	attrMap["Class"] = ""
	attrMap["Type"] = "text"
	attrMap["Validate"] = ""
	attrMap["Uijq"] = ""
	unit := field.Unit
	prefix := ""
	postfix := ""
	switch field.Input {
	case ii.IInt:
		attrMap["Validate"] = "digits:true, "
	case ii.IFloat:
		attrMap["Validate"] = "number:true, "
	case ii.IPercent:
		attrMap["Validate"] = "number:true, "
		unit = "%"
	case ii.IMoney:
		attrMap["Validate"] = "number:true, "
		prefix = `<span class="input-group-addon">￥</span>`
	case ii.IDate, ii.IDateTime:
		curTime := ""
		if value == "" {
			if field.Default == "curtime" {
				value = "curtime"
			}
			if field.Default == "curdate" {
				value = "curtime"
			}
		}
		if value == "curtime" || value == "curdate" {
			curTime = "value:new Date(), "
		}
		attrMap["Uijq"] = fmt.Sprintf(uiJqDateFormat, curTime)
	case ii.IPassword:
		attrMap["Type"] = "password"
	}
	if unit != "" {
		postfix = fmt.Sprintf(addonFormat, unit)
	}

	attrMap["Postfix"] = postfix
	attrMap["Prefix"] = prefix
	if prefix != "" || postfix != "" {
		attrMap["InputGroupStart"] = inputGroupStart
		attrMap["InputGroupEnd"] = inputGroupEnd
	} else {
		attrMap["InputGroupStart"] = ""
		attrMap["InputGroupEnd"] = ""
	}
	b := new(bytes.Buffer)
	eErr := textTpl.Execute(b, attrMap)
	if eErr != nil {
		lg.Critical("createTextInput error", eErr)
		return ""
	}
	return b.String()
}

const tableView = `

<table class="table table-striped m-b-none">

</table>`

func createViewPanel(vMap map[string]interface{}) {

}

//func BuildSelectElement(name, label, require, status string, valueMaps []map[string]interface{}, defaultValue interface{}, valueField string, labelField string)string{
//	options := BuildSelectOptions(valueMaps, defaultValue , valueField, labelField)
//	return fmt.Sprintf(selectFormat, label, require, label, name, name, u.ToStr(defaultValue), status, options)
//}

func BuildSelectOptions(valueMaps []cs.MObject, defaultValue interface{}, valueField string, labelField string, addFields ...string) string {
	options := ""
	for _, valueMap := range valueMaps {
		optionValue, ok := valueMap[valueField]
		if !ok {
			lg.Error("BuildSelectOptions: no such value field", valueField)
			continue
		}
		optionValueStr := ut.ToStr(optionValue)
		optionLabel, _ := valueMap[labelField]

		optionDatas := ""
		for _, addField := range addFields {
			addData := ut.Maps.GetStringDeepInMap(valueMap, addField)
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

var textTplFormat = `        <div class="form-group">
            <label class="col-sm-3 control-label">{{.RequireFlag}}{{.Label}}</label>
              <div class="col-sm-7">
              {{.InputGroupStart}} {{.Prefix}}
                <input type="{{.Type}}" class="input-block-level form-control {{.Class}}" name="{{.Name}}" id="{{.Name}}" autocomplete="off" value="{{.Value}}"
                  data-validate="{required: {{.Require}}, {{.Validate}} messages:{required:'请输入正确的{{.Label}}!'}}"
                  {{.Uijq}} {{.Status}}/>
                  {{.Postfix}}{{.InputGroupEnd}}
                <span class="help-block" id="{{.Name}}sHelpBlock">{{.Desc}}</span>
            </div>
        </div>
`
var textTpl *template.Template

func init() {
	textTpl = template.New("uibuilder_textTpl")
	_, err := textTpl.Parse(textTplFormat)
	if err != nil {
		lg.Critical("uibuilder init textTpl error", err)
		return
	}
}
