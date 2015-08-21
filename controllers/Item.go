package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/lang"
	"webo/models/status"
	"webo/models/svc"
)

type ItemController struct {
	BaseController
}

func (this *ItemController) List() {
	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	//	fmt.Println("params", this.Ctx.Input.Params)
	//	fmt.Println("requestBosy", this.Input()["id"])
	//	fmt.Println("ge", this.GetString("xx"))
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	//	fmt.Println("requestMap", requestMap)
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
	queryParams := make(svc.Params, 0)
	limitParams := make(map[string]int64, 0)
	beego.Debug("List Item user", requestMap)
	if k, ok := requestMap["limit"]; ok {
		limitParams["limit"] = int64(k.(float64))
	}
	if k, ok := requestMap["offset"]; ok {
		limitParams["offset"] = int64(k.(float64))
	}
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
	result, total, resultMaps := svc.List(item, queryParams, limitParams, orderByParams)
	//fmt.Println(result, total, retList)
	retList := transList(oItemDef, resultMaps)
	//fmt.Println("retList", retList)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
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
func (this *ItemController) Get() {
	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	//	fmt.Println("params", this.Ctx.Input.Params)
	tr := new(TableResult)
	tr.Rows = []map[string]string{{"id": "1", "user": "user1", "name": "a", "department": "dep1", "role": "admin", "flat": ""}}
	tr.Total = 1
	this.Data["json"] = tr
	this.ServeJson()
}
func (this *ItemController) Add() {
	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	//	fmt.Println("params", this.Ctx.Input.Params)
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		fmt.Println("no_")
	}
	svcParams := this.GetFormValues(oEntityDef)
	status, reason := svc.Add(item, svcParams)
	//fmt.Println("addservice", ret)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}
func (this *ItemController) Update() {
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
	//fmt.Println("addservice", ret)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
	this.ServeJson()
}
func (this *ItemController) Delete() {
	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	fmt.Println("params", this.Ctx.Input.Params)
	tr := new(TableResult)
	tr.Rows = []map[string]string{{"id": "1", "user": "user1", "name": "a", "department": "dep1", "role": "admin", "flat": ""}}
	tr.Total = 1
	this.Data["json"] = tr
	this.ServeJson()
}
