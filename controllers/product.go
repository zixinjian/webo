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
	item := s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["Service"] = "/item/add/" + oItemDef.Name
	this.FillFormElement(ui.BuildFormElement(oItemDef, t.Params{s.Sn:u.TUId()}, map[string]string{}))
	_, categorys := svc.GetAll(s.Category)
	this.Data["CategoryOptions"] = ui.BuildSelectOptions(categorys, "", s.Key, s.Name, s.Flag)
	this.TplNames = "product/add.tpl"
}
func (this *ProductController) UiSetting() {
	this.TplNames = "product/catagorySetting.tpl"
}

func (this *ProductController) UiList() {
	item:=s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/item/product/list"
	this.Data["addUrl"] = "ui/product/add"
	this.Data["updateUrl"] = "ui/product/update"
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "product/list.tpl"
}

func (this *ProductController) UiUpdate() {
	item := s.Product
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
		oldValueMap = transProductMap(oldValueMap)
		this.FillFormElement(ui.BuildFormElement(oItemDef, oldValueMap,map[string]string{}))
		_, suppliers := supplierMgr.GetSupplierListByProductSn(sn)
		this.Data["supplierList"] = suppliers
		_, categorys := svc.GetAll(s.Category)
		this.Data["CategoryOptions"] = ui.BuildSelectOptions(categorys, "", s.Key, s.Name, s.Flag)
		this.TplNames = "product/update.tpl"
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
func (this *ProductController) Add() {
	item := s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	svcParams, exParams := this.GetExFormValues(oItemDef)
	beego.Debug("ProductController.Add:", svcParams)
	suppliers, ok := exParams["supplierlist"]
	if !ok{
		this.JsonError(stat.Failed, "添加至少一个供应商！")
		return
	}
	sn := u.TUId()
	supplierMgr.AddSupplierRels(sn, suppliers.([]string))
	svcParams[s.Sn]= sn
	status, reason := svc.Add(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}

func (this *ProductController) Update() {
	item := s.Product
	oItemDef, _ := itemDef.EntityDefMap[item]
	svcParams, exParams := this.GetExFormValues(oItemDef)
	beego.Debug("ProductController.Update:", svcParams)
	suppliers, ok := exParams["supplierlist"]
	if !ok{
		this.JsonError(stat.Failed, "添加至少一个供应商！")
		return
	}
	sn := this.GetString(s.Sn)
	supplierMgr.UpdateSupplierRels(sn, suppliers.([]string))
	status, reason := svc.Update(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
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
	var retMap = make(map[string]interface{}, len(oldMap) + 1)
	for key, value := range oldMap {
		switch key {
		case s.Category:
			retMap[key] = lang.GetLabel(value.(string))
		default:
			retMap[key] = value
		}
	}
	productSn, _ := retMap[s.Sn]
	retMap[s.Supplier + s.Ext + s.Enum] = getSupplierList(productSn.(string))

	return retMap
}
func getSupplierList(productSn string)[]map[string]interface{}{
	if status, supplierMaps := supplierMgr.GetSupplierListByProductSn(productSn); strings.EqualFold(status, stat.Success) {
		return supplierMaps
	}
	return make([]map[string]interface{}, 0)
}
//func expandProductMap(oldMap map[string]interface{}) map[string]interface{} {
//	var retMap = make(map[string]interface{}, len(oldMap)+3)
//	for key, value := range oldMap {
//		switch key {
//		case s.Category:
//			retMap[key] = lang.GetLabel(value.(string))
//		case s.Supplier:
//			if !strings.EqualFold(value.(string), "") {
//				if supplierMap, sok := supplierMgr.Get(value.(string)); sok {
//					supplierKey, _ := supplierMap[s.Keyword]
//					supplierName, _ := supplierMap[s.Name]
//					retMap[s.Supplier+s.EKey] = supplierKey.(string)
//					retMap[s.Supplier+s.EName] = supplierName.(string)
//					retMap[s.Supplier] = value
//				}
//			}
//		default:
//			retMap[key] = value
//		}
//	}
//	productSn, _ := retMap[s.Sn]
//	retMap[s.Supplier + s.Ext + s.List] = getSupplierList(productSn.(string))
//	return retMap
//}
