package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"strings"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/svc"
	"webo/models/t"
)

type ItemController struct {
	BaseController
}

func (this *ItemController) ListWithQuery(oItemDef itemDef.ItemDef, addQueryParam t.Params) {
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	for k, v := range addQueryParam {
		queryParams[k] = v
	}
	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := transList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

func (this *ItemController) List() {
	item, ok := this.Ctx.Input.Params[":hi"]
	beego.Debug("params", this.Ctx.Input.Params, this.Ctx.Input)
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
	addParams := this.GetFormValues(oItemDef)
	this.ListWithQuery(oItemDef, addParams)
}

func (this *ItemController) Add() {
	item, ok := this.Ctx.Input.Params[":hi"]
//	beego.Debug("params", this.Ctx.Input.Params, this.Ctx.Input)
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Data["json"] = JsonResult{stat.ParamItemError, stat.ParamItemError}
		this.ServeJson()
		return
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Data["json"] = JsonResult{stat.ItemNotDefine, stat.ItemNotDefine}
		this.ServeJson()
		return
	}
	curUserSn := this.GetSessionString(SessionUserSn)
	svcParams := this.GetFormValues(oEntityDef)
	svcParams[s.Creater] = curUserSn
	for _, field := range oEntityDef.Fields {
		if strings.EqualFold(field.Model, s.Upload) {
			delete(svcParams, field.Name)
		}
	}
	beego.Debug("ItemController.Add:", svcParams)
	status, reason := svc.Add(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *ItemController) Update() {
	beego.Debug("Update requestBody: ", this.Ctx.Input.RequestBody)
	beego.Debug("Update params:", this.Ctx.Input.Params)
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Data["json"] = JsonResult{stat.ParamItemError, stat.ParamItemError}
		this.ServeJson()
		return
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Data["json"] = JsonResult{stat.ItemNotDefine, stat.ItemNotDefine}
		this.ServeJson()
		return
	}
	svcParams := this.GetFormValues(oEntityDef)
	status, reason := svc.Update(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *ItemController) Delete() {
	beego.Debug("Update requestBody: ", this.Ctx.Input.RequestBody)
	beego.Debug("Update params:", this.Ctx.Input.Params)
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Data["json"] = JsonResult{stat.ParamItemError, stat.ParamItemError}
		this.ServeJson()
		return
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Data["json"] = JsonResult{stat.ItemNotDefine, stat.ItemNotDefine}
		this.ServeJson()
		return
	}
	svcParams := this.GetFormValues(oEntityDef)
	sn, ok := svcParams[s.Sn]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s Sn not define", item))
		this.Data["json"] = JsonResult{stat.ItemNotDefine, stat.ItemNotDefine}
		this.ServeJson()
	}
	status := svc.Delete(item, sn.(string), this.GetCurUserSn())
	this.Data["json"] = &JsonResult{status, ""}
	this.ServeJson()
}

func (this *ItemController) Upload() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error(stat.ParamItemError)
		this.Ctx.WriteString(stat.ParamItemError)
		return
	}
	if _, ok := itemDef.EntityDefMap[item]; !ok {
		beego.Error(fmt.Sprintf("Item %s not define", item))
		this.Ctx.WriteString(stat.ItemNotDefine)
		return
	}
	sn := this.GetString(s.Sn)
	if sn == "" {
		beego.Error("ItemController.Upload error: ", stat.SnNotFound)
		this.Ctx.WriteString(stat.SnNotFound)
	}
	f, h, e := this.GetFile("uploadFile")
	fmt.Println(f, h, e)
	if e != nil {
		beego.Error("Upload error", e.Error())
		return
	}
	f.Close()
	saveDir := fmt.Sprintf("static/files/%s/%s/", item, sn)
	err := os.MkdirAll(saveDir, 0777)
	if err != nil {
		beego.Error("ItemController.Upload error: ", stat.UploadErrorCreateDir)
		this.Ctx.WriteString(stat.UploadErrorCreateDir)
		return
	}
	this.SaveToFile("uploadFile", saveDir+h.Filename)
	this.Ctx.WriteString("ok")
}

func (this *ItemController) Autocomplete() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		beego.Error("ItemController.Autocomplete: ", stat.ParamItemError)
		this.Data["json"] = "[]"
		this.ServeJson()
		return
	}
	_, vok := itemDef.EntityDefMap[item]
	if !vok {
		beego.Error("ItemController.Autocomplete: ", stat.ItemNotDefine)
		this.Data["json"] = "[]"
		this.ServeJson()
		return
	}
	var keyword string
	switch item {
	case s.Product:
		keyword = s.Model
	case s.Supplier:
		keyword = s.Keyword
	case s.User:
		keyword = s.UserName
	default:
		keyword = s.Name
	}
	this.BaseAutocomplete(item, keyword)
}

func (this *ItemController) BaseAutocomplete(item string, keyword string) {
	oItemDef, _ := itemDef.EntityDefMap[item]
	term := this.GetString(s.Term)
	if strings.EqualFold(term, "") {
		this.Data["json"] = "[]"
		this.ServeJson()
		return
	}
	limitParams := t.LimitParams{
		s.Limit: t.LimitDefault,
	}

	orderByParams := t.Params{
		keyword: s.Asc,
	}
	addParams := this.GetFormValues(oItemDef)
	addParams["%" + keyword] = term
	_, _, resultMaps := svc.List(oItemDef.Name, addParams, limitParams, orderByParams)
	retList := TransAutocompleteList(resultMaps, keyword)
	this.Data["json"] = &retList
	this.ServeJson()
}
