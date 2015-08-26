package controllers

import (
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
	"strings"
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
			if v, fok := field.GetValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = v
			}
		}
	}
	return retMap
}

func (this *BaseController) GetRequestParams(itemD itemDef.ItemDef)svc.Params{
	queryParams := this.GetFormValues(itemD)
//	if()
	return queryParams
}

func (this *BaseController) GetSessionString(sessionName string) string {
	return this.GetSession(sessionName).(string)
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

func getLimitParamFromRequestMap(requestMap map[string]interface{}) map[string]int64{
	limitParams := make(map[string]int64, 0)
//	beego.Debug("List Item user", requestMap)
	if k, ok := requestMap["limit"]; ok {
		limitParams["limit"] = int64(k.(float64))
	}
	if k, ok := requestMap["offset"]; ok {
		limitParams["offset"] = int64(k.(float64))
	}
	return limitParams
}
func getOrderParamFromRequestMap(requestMap map[string]interface{}) svc.Params{
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


