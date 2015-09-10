package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
)

type UiController struct {
	BaseController
}

func (this *UiController) List() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Ctx.WriteString(stat.ParamItemError)
		return
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Ctx.WriteString(stat.ItemNotDefine)
		return
	}
	this.Data["item"] = item
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "item/list.html"
}

func (this *UiController) Add() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Ctx.WriteString(stat.ParamItemError)
		return
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Ctx.WriteString(stat.ItemNotDefine)
		return
	}
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(oItemDef, u.TUId())
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "item/add.tpl"
}

func (this *UiController) Update() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Ctx.WriteString(stat.ParamItemError)
		return
	}
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Ctx.WriteString(stat.ItemNotDefine)
		return
	}
	sn := this.GetString(s.Sn)
	if sn == "" {
		beego.Error("ui.Update", stat.ParamSnIsNone)
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := t.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		//		fmt.Println("oldValue", oldValueMap)
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "item/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}
