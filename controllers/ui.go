package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"webo/models/svc"
	//	"webo/models/itemDef"
	"github.com/astaxie/beego"
	"webo/models/itemDef"
)

type UiController struct {
	BaseController
}

func (this *UiController) Add() {
	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	//	fmt.Println("params", this.Ctx.Input.Params)
	//	fmt.Println("requestBosy", this.Input()["id"])
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error("Item param is none")
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(oItemDef)
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "item/add.tpl"
}

func (this *UiController) List() {
	//	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	//	fmt.Println("params", this.Ctx.Input.Params)
	//	fmt.Println("requestBosy", this.Input()["id"])
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error("Item param is none")
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	//	oEntityDef, ok := itemDef.EntityDefMap[item]
	//	fmt.Println("form", this.GetFormValues(oEntityDef))
	this.Data["item"] = item
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "item/list.html"
}

func (this *UiController) Update() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error("Item param is none")
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	sn := this.GetString("sn")
	if sn == "" {
		fmt.Println("sn is none")
	}
	fmt.Println("sn", sn)
	params := svc.Params{"sn": sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		fmt.Println("oldValue", oldValueMap)
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.TplNames = "item/update.html"
	} else {
		this.Ctx.WriteString("Id not found")
	}
}
