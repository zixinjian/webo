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

type TravelController struct {
	BaseController
}

func (this *TravelController) UiList() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Travel
	this.Data["item"] = item
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = "/travel/ui/add"
	this.Data["updateUrl"] = "/travel/ui/update"
	this.TplNames = "travel/list.html"
}

func (this *TravelController) UiAdd() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Travel
	oItemDef, _ := itemDef.EntityDefMap[item]
	oldValueMap := map[string]interface{}{
		s.Sn:           u.TUId(),
		s.ApproverName: this.GetCurUser(),
		s.Approver:     this.GetCurUserSn(),
	}
	approverField, _ := oItemDef.GetField(s.Approver)
	approverField.Name = "approvername"
	approverSn, _ := oItemDef.GetField(s.Approver)
	approverSn.Input = s.InputHidden
	oItemDef.Fields = append(oItemDef.Fields[:len(oItemDef.Fields)])
	this.Data["Form"] = ui.BuildUpdatedFormWithStatus(oItemDef, oldValueMap, make(map[string]string))
	this.Data["Service"] = "/item/add/" + item
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "travel/add.tpl"
}
func (this *TravelController) UiUpdate() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Travel
	oItemDef, _ := itemDef.EntityDefMap[item]
	sn := this.GetString(s.Sn)
	if sn == "" {
		beego.Error("TravelController.UiUpdate", stat.ParamSnIsNone)
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := t.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		this.Data["Service"] = "/item/update/" + item
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "travel/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

//func (this *TravelController) Update() {
//	item := s.Travel
//	oEntityDef, _ := itemDef.EntityDefMap[item]
//	svcParams := this.GetFormValues(oEntityDef)
//	if pwd, ok := svcParams[s.Password]; ok {
//		if strings.EqualFold(pwd.(string), "*****") {
//			delete(svcParams, s.Password)
//		}
//	}
//	status, reason := svc.Update(item, svcParams)
//	this.Data["json"] = &JsonResult{status, reason}
//	this.ServeJson()
//}
