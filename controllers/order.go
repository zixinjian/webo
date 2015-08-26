package controllers

import (
	"fmt"
	"webo/controllers/ui"
//	"webo/models/svc"
	//	"webo/models/itemDef"
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"encoding/json"
	"webo/models/svc"
)

type OrderController struct {
	BaseController
}

//func (this *OrderController) Add() {
//	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
//	//	fmt.Println("params", this.Ctx.Input.Params)
//	//	fmt.Println("requestBosy", this.Input()["id"])
//	item, ok := this.Ctx.Input.Params[":hi"]
//	if !ok {
//		beego.Error("Item param is none")
//	}
//	oItemDef, ok := itemDef.EntityDefMap[item]
//	if !ok {
//		beego.Error(fmt.Sprintf("Item %s not support", item))
//	}
//	this.Data["Service"] = "/item/add/" + item
//	this.Data["Form"] = ui.BuildAddForm(oItemDef)
//	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
//	this.TplNames = "item/add.tpl"
//}

func (this *OrderController) UiMyCreate() {
	beego.Error("UiMyCreate")
	item:="order"
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s?create", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "item/list.html"
}

func (this *OrderController) MyCreate() {
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
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
	queryParams := svc.Params{
		"":"",
	}
	limitParams := getLimitParamFromRequestMap(requestMap)
	orderByParams := getOrderParamFromRequestMap(requestMap)
	result, total, resultMaps := svc.List(item, queryParams, limitParams, orderByParams)
	retList := transList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

//func (this *OrderController) Update() {
//	item, ok := this.Ctx.Input.Params[":hi"]
//	if !ok {
//		beego.Error("Item param is none")
//	}
//	oItemDef, ok := itemDef.EntityDefMap[item]
//	if !ok {
//		beego.Error(fmt.Sprintf("Item %s not support", item))
//	}
//	sn := this.GetString("sn")
//	if sn == "" {
//		fmt.Println("sn is none")
//	}
//	fmt.Println("sn", sn)
//	params := svc.Params{"sn": sn}
//	code, oldValueMap := svc.Get(item, params)
//	if code == "success" {
//		fmt.Println("oldValue", oldValueMap)
//		this.Data["Service"] = "/item/update/" + item
//		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
//		this.TplNames = "item/update.html"
//	} else {
//		this.Ctx.WriteString("Id not found")
//	}
//}
