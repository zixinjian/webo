package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/lang"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/t"
	"webo/models/userMgr"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) JsonError(status string, reason string){
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *BaseController) GetItemDefFromParamHi() (itemDef.ItemDef, string) {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemIsNone_, this.Ctx.Input.Params)
		return itemDef.ItemDef{}, stat.ParamItemError
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(stat.ItemNotDefine_, item)
		return itemDef.ItemDef{}, stat.ItemNotDefine
	}
	return oItemDef, stat.Success
}

func (this *BaseController) GetParams(oItemDef itemDef.ItemDef) (queryParams t.Params, limitParams t.LimitParams, orderByParams t.Params) {
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("PurchaseController.List requestMap: ", requestMap)

	limitParams = this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams = this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	queryParams = this.GetQueryParamFromJsonMap(requestMap, oItemDef)
	addParams := this.GetFormValues(oItemDef)
	for k, v := range addParams {
		queryParams[k] = v
	}
	return queryParams, limitParams, orderByParams
}

func (this *BaseController) GetFormValues(itemD itemDef.ItemDef) map[string]interface{} {
	retMap := make(map[string]interface{})
	formValues := this.Input()
	beego.Debug("BaseController.GetFormValues from values: ", formValues)
	for k, _ := range formValues {
		if field, ok := itemD.GetField(k); ok {
			if v, fok := field.GetFormValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = this.ReplaceSpecialValues(v)
			}
		}
	}
	return retMap
}
func (this *BaseController) GetExFormValues(itemD itemDef.ItemDef)(map[string]interface{}, map[string]interface{}){
	retMap := make(map[string]interface{})
	exMap := make(map[string]interface{})
	formValues := this.Input()
	beego.Debug("BaseController.GetFormValues from values: ", formValues)
	for k, v := range formValues {
		if field, ok := itemD.GetField(k); ok {
			if fv, fok := field.GetFormValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = this.ReplaceSpecialValues(fv)
			}
		}else{
			exMap[k] = v
		}
	}
	retMap[s.Creater] = this.GetCurUserSn()
	for _, field := range itemD.Fields {
		if strings.EqualFold(field.Model, s.Upload) {
			delete(retMap, field.Name)
		}
	}
	return retMap, exMap
}
func (this *BaseController) ReplaceSpecialValues(value interface{}) interface{} {
	if str, ok := value.(string); ok {
		rValue := strings.TrimSpace(str)
		switch rValue {
		case s.CurUser:
			return this.GetCurUserSn()
		default:
			return value
		}
	} else {
		return value
	}
}

func (this *BaseController) GetSessionString(sessionName string) string {
	if this.GetSession(sessionName) != nil {
		return this.GetSession(sessionName).(string)
	}
	// TODO
	return ""
}

func (this *BaseController) GetCurUserSn() string {
	return this.GetSessionString(SessionUserSn)
}
func (this *BaseController) GetCurUserUserName() string {
	return this.GetSessionString(SessionUserUserName)
}
func (this *BaseController) GetCurUserName() string {
	return this.GetSessionString(SessionUserName)
}
func (this *BaseController) GetCurRole() string {
	return this.GetSessionString(SessionUserRole)
}
func (this *BaseController) GetCurDepartment() string {
	return this.GetSessionString(SessionUserDepartment)
}

func (this *BaseController) GetQueryParamFromJsonMap(requestMap map[string]interface{}, oItemDef itemDef.ItemDef) map[string]interface{} {
	queryParams := make(t.Params, 0)
	fieldMap := oItemDef.GetFieldMap()
	for k, v := range requestMap {
		if field, ok := fieldMap[k]; ok {
			if fv, fok := field.GetCheckedValue(v); fok {
				queryParams[k] = fv
			} else {
				beego.Error(fmt.Sprintf("Check param[%s]value %v error", k, v))
			}
		} else {
			beego.Error(fmt.Sprintf("Check param[%s]value %v error no such field", k, v))
		}
	}
	return queryParams
}

func (this *BaseController) GetLimitParamFromJsonMap(requestMap map[string]interface{}) map[string]int64 {
	limitParams := make(map[string]int64, 0)
	if k, ok := requestMap["limit"]; ok {
		limitParams["limit"] = int64(k.(float64))
	}
	if k, ok := requestMap["offset"]; ok {
		limitParams["offset"] = int64(k.(float64))
	}
	return limitParams
}
func (this *BaseController) GetOrderParamFromJsonMap(requestMap map[string]interface{}) t.Params {
	orderByParams := make(t.Params, 0)
	if sort, ok := requestMap["sort"]; ok {
		sortStr := strings.TrimSpace(sort.(string))
		if sortStr != "" {
			order := "asc"
			if o, ok := requestMap["order"]; ok {
				if strings.TrimSpace(o.(string)) == "desc" {
					order = "desc"
				}
			}
			orderByParams[sortStr] = order
		}
	}
	return orderByParams
}

func TransAutocompleteList(resultMaps []map[string]interface{}, keyField string) []map[string]interface{} {
	if len(resultMaps) <= 0 {
		return resultMaps
	}
	if keyField == s.Keyword {
		return resultMaps
	}
	for idx, oldMap := range resultMaps {
		for k, v := range oldMap {
			if keyField == k {
				oldMap[s.Keyword] = v
			}
		}
		resultMaps[idx] = oldMap
	}
	return resultMaps
}

func transList(oItemDef itemDef.ItemDef, resultMaps []map[string]interface{}) []map[string]interface{} {
	if len(resultMaps) < 0 {
		return resultMaps
	}
	retList := make([]map[string]interface{}, len(resultMaps))
	neetTransMap := oItemDef.GetNeedTrans()
	for idx, oldMap := range resultMaps {
		var retMap = make(map[string]interface{}, len(oldMap))
		for key, value := range oldMap {
			if _, ok := neetTransMap[key]; ok {
				retMap[key] = lang.GetLabel(value.(string))
			} else {
				retMap[key] = value
			}
		}
		retList[idx] = retMap
	}
	return retList
}

func makeFields(oItemDef itemDef.ItemDef, names []string)itemDef.ItemDef{
	fields := make([]itemDef.Field, len(names))

	fieldMap := oItemDef.GetFieldMap()
	for idx, name := range names {
		if field, ok := fieldMap[name]; ok {
			fields[idx] = field
		} else {
			beego.Error("Field not found", name, idx)
		}
	}
	oItemDef.Fields = fields
	return oItemDef
}

func FillUserEnum(fieldName string, oItemDef itemDef.ItemDef, queryParams t.Params, orderParams t.Params) itemDef.ItemDef {
	for idx, field := range oItemDef.Fields {
		if strings.EqualFold(fieldName, field.Name) {
			field.Enum = userMgr.GetUserEnum(queryParams, orderParams, t.LimitParams{})
		}
		oItemDef.Fields[idx] = field
	}
	return oItemDef
}


func extendUserField(oItemDef itemDef.ItemDef, fieldName string)itemDef.ItemDef{
	userField, _ := oItemDef.GetField(fieldName)
	userField.Name = fieldName + s.Name
	userSnField, _ := oItemDef.GetField(fieldName)
	userSnField.Input = s.Hidden
	newFields := make([]itemDef.Field, 0)
	for _, field := range oItemDef.Fields{
		if field.Name == fieldName{
			newFields = append(newFields, userField)
			newFields = append(newFields, userSnField)
		}else{
			newFields = append(newFields, field)
		}
	}
	oItemDef.Fields = newFields
	return oItemDef
}

func extendUser2ItemMapList(resultMaps []map[string]interface{}, fields ...string)[]map[string]interface{}{
	for idx, travelMap := range resultMaps{
		resultMaps[idx]=extendUser2ItemMap(travelMap, fields...)
	}
	return resultMaps
}

func extendUser2ItemMap(travelMap map[string]interface{}, fields ...string)map[string]interface{}{
	for _, fieldName := range fields{
		if sn, ok := travelMap[fieldName]; ok{
			if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
				name, _ := userMap[s.Name]
				travelMap[fieldName + s.Name]= name
			}
		}
	}
	return travelMap
}