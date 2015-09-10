package controllers

import (
	"github.com/astaxie/beego"
	"strings"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/lang"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/supplierMgr"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
)

type ProductController struct {
	ItemController
}

func (this *ProductController) UiAdd() {
	item := "product"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["Service"] = "/item/add/" + oItemDef.Name
	this.Data["Form"] = ui.BuildAddForm(oItemDef, u.TUId())
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "product/add.tpl"
}

func (this *ProductController) UiList() {
	item := "product"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/item/product/list"
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
	params := t.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		oldValueMap = expandProductMap(oldValueMap)
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "product/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *ProductController) List() {
	item := s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := transProductList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

func transProductList(oItemDef itemDef.ItemDef, resultMaps []map[string]interface{}) []map[string]interface{} {
	if len(resultMaps) < 0 {
		return resultMaps
	}
	retList := make([]map[string]interface{}, len(resultMaps))
	for idx, oldMap := range resultMaps {
		retList[idx] = transProductMap(oldMap)
	}
	return retList
}

func transProductMap(oldMap map[string]interface{}) map[string]interface{} {
	var retMap = make(map[string]interface{}, len(oldMap))
	for key, value := range oldMap {
		switch key {
		case s.Category:
			retMap[key] = lang.GetLabel(value.(string))
		case s.Supplier:
			if !strings.EqualFold(value.(string), "") {
				if supplierMap, sok := supplierMgr.Get(value.(string)); sok {
					supplier, _ := supplierMap[s.Name]
					retMap[s.Supplier] = supplier.(string)
				}
			}
		default:
			retMap[key] = value
		}
	}
	return retMap
}

func expandProductMap(oldMap map[string]interface{}) map[string]interface{} {
	var retMap = make(map[string]interface{}, len(oldMap)+2)
	for key, value := range oldMap {
		switch key {
		case s.Category:
			retMap[key] = lang.GetLabel(value.(string))
		case s.Supplier:
			if !strings.EqualFold(value.(string), "") {
				if supplierMap, sok := supplierMgr.Get(value.(string)); sok {
					supplierKey, _ := supplierMap[s.Keyword]
					supplierName, _ := supplierMap[s.Name]
					retMap[s.Supplier+s.EKey] = supplierKey.(string)
					retMap[s.Supplier+s.EName] = supplierName.(string)
					retMap[s.Supplier] = value
				}
			}
		default:
			retMap[key] = value
		}
	}
	return retMap
}
