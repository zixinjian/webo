package controllers

import (
	"encoding/json"
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
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("ListWithQuery requestMap: ", requestMap)

	limitParams := this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams := this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	queryParams := this.GetQueryParamFromJsonMap(requestMap, oItemDef)
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
	fmt.Println("params", this.Ctx.Input.Params, this.Ctx.Input)
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
	//	creater := this.GetString(s.Creater)
	//	if creater == s.CurUser {
	//		sn := this.GetCurUserSn()
	//		addParams[s.Creater]= sn
	//	}
	this.ListWithQuery(oItemDef, addParams)
}

func (this *ItemController) Add() {
	beego.Debug("BaseController.GetFormValues form values: ", this.Input())
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		fmt.Println(stat.ItemNotDefine)
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
		fmt.Println("hi", item)
	}
	oEntityDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		fmt.Println(stat.ItemNotDefine)
	}
	svcParams := this.GetFormValues(oEntityDef)
	status, reason := svc.Update(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *ItemController) Upload() {
	item, ok := this.Ctx.Input.Params[":hi"]
	if !ok {
		fmt.Println("hi", item)
	}
	_, vok := itemDef.EntityDefMap[item]
	if !vok {
		fmt.Println(stat.ItemNotDefine)
	}
	sn := this.GetString("sn")
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
	case s.Supplier, s.Product:
		keyword = s.Keyword
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
		s.Keyword: s.Asc,
	}

	queryParams := t.Params{
		"%" + s.Keyword: term,
	}
	_, _, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := TransAutocompleteList(resultMaps, keyword)
	this.Data["json"] = &retList
	this.ServeJson()
}
