package controllers

import (
	"encoding/json"
	"fmt"
	"webo/models/itemDef"
	"webo/models/lang"
	"webo/models/status"
	"webo/models/svc"
	"webo/models/s"
	"github.com/astaxie/beego"
)

type ItemController struct {
	BaseController
}
func (this *ItemController) ListWithQuery(oItemDef itemDef.ItemDef, addQueryParam svc.Params) {
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("requestMap", requestMap)

	limitParams := this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams := this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	queryParams := this.GetQueryParamFromJsonMap(requestMap, oItemDef)
	for k, v := range addQueryParam{
		queryParams[k]=v
	}

	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := transList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

func (this *ItemController) List() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		this.Data["json"] = TableResult{"false", 0, ""}
		this.ServeJson()
		return
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		this.Data["json"] = TableResult{"false", 0, ""}
		this.ServeJson()
		return
	}
	addParams := svc.Params{}
	creater := this.GetString(s.Creater)
	if creater == s.CurUser {
		sn := this.GetCurUserSn()
		addParams[s.Creater]= sn
	}
	this.ListWithQuery(oItemDef, addParams)
}


func transList(oItemDef itemDef.ItemDef, resultMaps []map[string]interface{}) []map[string]interface{} {
	if len(resultMaps) < 0 {
		return resultMaps
	}
	retList := make([]map[string]interface{}, len(resultMaps))
	neetTransMap := oItemDef.GetNeedTrans()
	//fmt.Println("neetTransMap", neetTransMap)
	for idx, oldMap := range resultMaps {
		var retMap = make(map[string]interface{}, len(oldMap))
		for key, value := range oldMap {
			if _, ok := neetTransMap[key]; ok {
				//fmt.Println("need", key, value, lang.GetLabel(value.(string)))
				retMap[key] = lang.GetLabel(value.(string))
			} else {
				retMap[key] = value
			}
		}
		//fmt.Println("retMap", retMap)
		retList[idx] = retMap
	}
	return retList
}
func (this *ItemController) Add() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		fmt.Println("no_")
	}
	curUserSn := this.GetSessionString(SessionUserSn)
	svcParams := this.GetFormValues(oEntityDef)
	svcParams[s.Creater] = curUserSn
	status, reason := svc.Add(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *ItemController) Update() {
	beego.Debug("Update requestBody: ", this.Ctx.Input.RequestBody)
	beego.Debug("Update params:", this.Ctx.Input.Params)
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		fmt.Println(status.ItemNotDefine)
	}
	svcParams := this.GetFormValues(oEntityDef)
	status, reason := svc.Update(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}
