package controllers

import (
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/util"
	"github.com/astaxie/beego"
	"webo/models/svc"
	"webo/models/stat"
	"webo/models/s"
	"encoding/json"
	"webo/models/lang"
	"strings"
	"webo/models/supplierMgr"
	"fmt"
)

type ProductController struct {
	ItemController
}

func (this *ProductController) UiAdd() {
	item := "product"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["Service"] = "/item/add/" + oItemDef.Name
	this.Data["Form"] = ui.BuildAddForm(oItemDef, util.TUId())
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "product/add.tpl"
}

func (this *ProductController) UiList() {
	item := "product"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/product/list"
	this.Data["addUrl"] = "ui/product/add"
	this.Data["updateUrl"] = "ui/product/update"
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "product/list.html"
}

func (this *ProductController) UiUpdate() {
	item := "product"
	oItemDef, _ := itemDef.EntityDefMap[item]
	sn := this.GetString(s.Sn)
	if sn == "" {
		beego.Error("UiUpdate error: ", stat.ParamSnIsNone)
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := svc.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
		this.TplNames = "product/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *ProductController)List(){
	item := s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("ProductController.list: ", requestMap)

	limitParams := this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams := this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	queryParams := this.GetQueryParamFromJsonMap(requestMap, oItemDef)
	addQueryParam := this.GetFormValues(oItemDef)
	for k, v := range addQueryParam{
		queryParams[k]=v
	}

	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := TransProductList(oItemDef, resultMaps)

	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

func TransProductList(oItemDef itemDef.ItemDef, resultMaps []map[string]interface{}) []map[string]interface{} {
	if len(resultMaps) < 0 {
		return resultMaps
	}
	retList := make([]map[string]interface{}, len(resultMaps))
	for idx, oldMap := range resultMaps {
		var retMap = make(map[string]interface{}, len(oldMap))
		for key, value := range oldMap {
			switch key {
			case s.Category:
				retMap[key] = lang.GetLabel(value.(string))
			case s.Supplier:
				if !strings.EqualFold(value.(string), ""){
					if supplierMap, sok := supplierMgr.Get(value.(string));sok{
						supplier, _ := supplierMap[s.Name]
						retMap[s.Supplier] = supplier.(string)
					}
				}
			default:
				retMap[key] = value
			}
		}
		retList[idx] = retMap
	}
	return retList
}