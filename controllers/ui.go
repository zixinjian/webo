package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"webo/models/svc"
	//	"webo/models/itemDef"
)

type UiController struct {
	BaseController
}

func (this *UiController) Add() {
	fmt.Println("requestBosy", this.Ctx.Input.RequestBody)
	fmt.Println("params", this.Ctx.Input.Params)
	fmt.Println("requestBosy", this.Input()["id"])
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	//	oEntityDef, ok := itemDef.EntityDefMap[item]
	//	fmt.Println("form", this.GetFormValues(oEntityDef))
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(item)
	this.TplNames = "add.html"
}

func (this *UiController) Update() {
	fmt.Println("Update", this.Ctx.Input.RequestBody)
	fmt.Println("params", this.Ctx.Input.Params)
	fmt.Println("requestBosy", this.Input()["id"])
	id := this.GetString("id")
	fmt.Println("id", id)
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	params := svc.Params{"id": id}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		//		fmt.Println("oldValue", oldValueMap)
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(item, oldValueMap)
		this.TplNames = "add.html"
	} else {
		this.Ctx.WriteString("Id not found")
	}
}
