package controllers

import (
	"github.com/astaxie/beego"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
	"webo/models/userMgr"
	"fmt"
)

type TravelController struct {
	BaseController
}
func (this *TravelController) Travel() {
	userName := this.GetCurUserUserName()
	userRole := this.GetCurRole()
	beego.Info(fmt.Sprintf("User:%s login as role:%s", userName, userRole))
	this.Data["userName"] = userName
	switch userRole {
	case "role_admin", "role_manager":
		navs := []ui.NavLiValue{
			{"active", "/travel/ui/myCreate", "我的申请"},
			{"", "/travel/ui/list", "待我审批的申请"},
			{"", "/travel/ui/list", "出差申请列表"},
		}
		this.Data["travelNav"] = ui.BuildNavs(navs)
		this.Data["userMgr"] = userMgrHtml
	default:
		navs := []ui.NavLiValue{
			{"active", "/travel/ui/list", "我的申请"},
		}
		this.Data["travelNav"] = ui.BuildNavs(navs)
		this.Data["userMgr"] = ""
	}
	this.TplNames = "travel.html"
}

func (this *TravelController) UiList() {
	if this.GetCurRole() == s.RoleUser {
		this.Ctx.WriteString("没有权限")
		return
	}
	item := s.Travel
	this.Data["item"] = item
	this.Data["listUrl"] = "/travel/item/list"
	this.Data["addUrl"] = "/travel/ui/add"
	this.Data["updateUrl"] = "/travel/ui/update"
	this.TplNames = "travel/list.tpl"
}
func (this *TravelController) MyCreate() {
	item := s.Travel
	this.Data["item"] = item
	this.Data["listUrl"] = "/travel/item/list"
	this.Data["addUrl"] = "/travel/ui/add"
	this.Data["updateUrl"] = "/travel/ui/update"
	this.TplNames = "travel/myCreate.html"
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
		s.ApproverName: this.GetCurUserName(),
		s.Approver:     this.GetCurUserSn(),
	}

	this.Data["Form"] = ui.BuildUpdatedFormWithStatus(extendUserField(oItemDef, s.Approver), oldValueMap, make(map[string]string))
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
		this.Data["Form"] = ui.BuildUpdatedForm(extendUserField(oItemDef, s.Approver), expandTravelMapForUpdate(oldValueMap))
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "travel/update.tpl"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *TravelController) List() {
	item := s.Travel
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := extendUser2ItemMapList(resultMaps, s.Approver, s.Traveler)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}
//func expandTravelList(resultMaps []map[string]interface{})[]map[string]interface{}{
//	for idx, travelMap := range resultMaps{
//		resultMaps[idx]=expandTravelMap(travelMap)
//	}
//	return resultMaps
//}
//func expandTravelMap(travelMap map[string]interface{})map[string]interface{}{
//	if sn, ok := travelMap[s.Approver]; ok{
//		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
//			name, _ := userMap[s.Name]
//			travelMap[s.ApproverSn]= travelMap[s.Approver]
//			travelMap[s.Approver] = name
//		}
//	}
//	if sn, ok := travelMap[s.Traveler];ok{
//		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
//			name, _ := userMap[s.Name]
//			travelMap[s.TravelerSn] = travelMap[s.Traveler]
//			travelMap[s.Traveler] = name
//		}
//	}
//	return travelMap
//}
func expandTravelMapForUpdate(travelMap map[string]interface{})map[string]interface{}{
	if sn, ok := travelMap[s.Approver]; ok{
		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
			name, _ := userMap[s.Name]
			travelMap[s.ApproverName] = name
		}
	}
	if sn, ok := travelMap[s.Traveler]; ok{
		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
			name, _ := userMap[s.Name]
			userName, _:= userMap[s.UserName]
			travelMap["traveler_key"] = userName
			travelMap["traveler_name"] = name
		}
	}
	beego.Debug("expandTravelMapForUpdate", travelMap)
	return travelMap
}
