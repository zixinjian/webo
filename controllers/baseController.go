package controllers

import (
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
	"strings"
	"fmt"
	"webo/models/stat"
	"webo/models/s"
	"webo/models/lang"
)

type BaseController struct {
	beego.Controller
}


func (this *BaseController) GetItemDefFromParamHi() (itemDef.ItemDef, string){
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

func (this *BaseController) GetFormValues(itemD itemDef.ItemDef) map[string]interface{} {
	var retMap map[string]interface{}
	retMap = make(map[string]interface{})
	formValues := this.Input()
	for _, field := range itemD.Fields {
		if _, ok := formValues[field.Name]; ok {
			if v, fok := field.GetFormValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = this.ReplaceSpecialValues(v)
			}
		}
	}
	return retMap
}

func (this *BaseController) ReplaceSpecialValues(value interface{}) interface{}{
	if str, ok := value.(string); ok{
		rValue := strings.TrimSpace(str)
		switch rValue {
		case s.CurUser:
			return this.GetCurUserSn()
		default:
			return value
		}
	}else{
		return value
	}
}

//func (this *BaseController) GetRequestParams(itemD itemDef.ItemDef)svc.Params{
//	queryParams := this.GetFormValues(itemD)
//	return queryParams
//}

func (this *BaseController) GetSessionString(sessionName string) string {
	if this.GetSession(sessionName) != nil{
		return this.GetSession(sessionName).(string)
	}
	// TODO
	return "d"
}

func (this *BaseController) GetCurUserSn() string {
	return this.GetSessionString(SessionUserSn)
}
func (this *BaseController) GetCurUser() string {
	return this.GetSessionString(SessionUserName)
}
func (this *BaseController) GetCurRole() string{
	return this.GetSessionString(SessionUserRole)
}
func (this *BaseController) GetCurDepartment() string{
	return this.GetSessionString(SessionUserDepartment)
}

func (this *BaseController)GetQueryParamFromJsonMap(requestMap map[string]interface{}, oItemDef itemDef.ItemDef) map[string]interface{}{
	queryParams := make(svc.Params, 0)
	fieldMap := oItemDef.GetFieldMap()
	for k, v := range requestMap{
		if field, ok := fieldMap[k];ok{
			if fv, fok := field.GetCheckedValue(v);fok{
				queryParams[k] = fv
			}else{
				beego.Error(fmt.Sprintf("Check param[%s]value %v error", k, v))
			}
		}else{
			beego.Error(fmt.Sprintf("Check param[%s]value %v error no such field", k, v))
		}
	}
	return queryParams
}

func (this *BaseController)GetLimitParamFromJsonMap(requestMap map[string]interface{}) map[string]int64{
	limitParams := make(map[string]int64, 0)
	if k, ok := requestMap["limit"]; ok {
		limitParams["limit"] = int64(k.(float64))
	}
	if k, ok := requestMap["offset"]; ok {
		limitParams["offset"] = int64(k.(float64))
	}
	return limitParams
}
func (this *BaseController)GetOrderParamFromJsonMap(requestMap map[string]interface{}) svc.Params{
	orderByParams := make(svc.Params, 0)
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
	retList := make([]map[string]interface{}, len(resultMaps))
	if len(resultMaps) <= 0 {
		return retList
	}

	for idx, oldMap := range resultMaps {
		var retMap = make(map[string]interface{}, 3)
		if sn, sok := oldMap[s.Sn]; sok{
			retMap[s.Sn] = sn
			if name, nok := oldMap[s.Name]; nok{
				retMap[s.Name] = name
			}else {
				retMap[s.Name] = ""
			}
			if keyword, kok := oldMap[s.Keyword]; kok{
				retMap[s.Keyword] = keyword
			}else {
				retMap[s.Keyword] = ""
			}
			retList[idx] = retMap
		}
	}
	return retList
}

func TransList(oItemDef itemDef.ItemDef, resultMaps []map[string]interface{}) []map[string]interface{} {
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