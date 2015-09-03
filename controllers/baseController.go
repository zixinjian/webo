package controllers

import (
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
	"strings"
	"fmt"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) GetFormValues(itemD itemDef.ItemDef) map[string]interface{} {
	//fmt.Println("def", itemD)
	var retMap map[string]interface{}
	retMap = make(map[string]interface{})
	formValues := this.Input()
	for _, field := range itemD.Fields {
		if _, ok := formValues[field.Name]; ok {
			if v, fok := field.GetFormValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = v
			}
		}
	}
	return retMap
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


