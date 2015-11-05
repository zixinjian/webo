package controllers

import (
	"webo/controllers/ui"
	"webo/models/s"
	"webo/models/itemDef"
	"webo/models/u"
	"webo/models/svc"
)

type AccountController struct {
	BaseController
}
func (this *AccountController) Account() {
	this.Data["userName"] = this.GetCurUserName()
	this.TplNames = "account.html"
}


func (this *AccountController) UiAdd() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Account
	oItemDef, _ := itemDef.EntityDefMap[item]
	oldValueMap := map[string]interface{}{
		s.Sn:           u.TUId(),
		s.PayerName: 	this.GetCurUserName(),
		s.Payer:     	this.GetCurUserSn(),
	}
	this.Data["Form"] = ui.BuildUpdatedFormWithStatus(extendUserField(oItemDef, s.Payer), oldValueMap, make(map[string]string))
	this.Data["Service"] = "/item/add/" + item
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "item/add.tpl"
}
func (this *AccountController) UiList() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Account
	this.Data["item"] = item
	this.Data["listUrl"] = "/account/item/list"
	this.Data["addUrl"] = "/account/ui/add"
	this.Data["updateUrl"] = "/ui/update/account"
	this.TplNames = "account/list.tpl"
}

func (this *AccountController) List() {
	item := s.Account
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := extendUser2ItemMapList(resultMaps, s.Payer)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}