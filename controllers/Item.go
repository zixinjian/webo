package controllers

import (
	"fmt"
	"webo/models/itemDef"
	"webo/models/svc"
	"encoding/json"
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
	}
//	oEntityDef, ok := itemDef.EntityDefMap[item]
	queryParams :=make(svc.Params, 0)
	limitParams :=make(map[string]int64, 0)
	if k, ok := requestMap["limit"]; ok{
		limitParams["limit"]=int64(k.(float64))
	}
	if k, ok:=requestMap["offset"];ok{
		limitParams["offset"]=int64(k.(float64))
	}
	orderByParams :=make(svc.Params, 0)
	result, total, retList := svc.List(item, queryParams, limitParams, orderByParams)
	fmt.Println(result, total, retList)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}
func (this *ItemController) Get() {
	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	fmt.Println("params", this.Ctx.Input.Params)
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
	svcParams := this.GetFormValues(oEntityDef)
	ret := svc.Add(item, svcParams)
	this.Data["json"] = &JsonResult{ret, ""}
	this.ServeJson()
}
func (this *ItemController) Update() {
	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	fmt.Println("params", this.Ctx.Input.Params)
	tr := new(TableResult)
	tr.Rows = []map[string]string{{"id": "1", "user": "user1", "name": "a", "department": "dep1", "role": "admin", "flat": ""}}
	tr.Total = 1
	this.Data["json"] = tr
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
