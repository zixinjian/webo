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
	this.Data["listUrl"] = "/travel/item/list"
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
		s.ApproverName: this.GetCurUserName(),
		s.Approver:     this.GetCurUserSn(),
	}

	this.Data["Form"] = ui.BuildUpdatedFormWithStatus(expandTraveDef(oItemDef), oldValueMap, make(map[string]string))
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
		this.Data["Form"] = ui.BuildUpdatedForm(expandTraveDef(oItemDef), expandTravelMapForUpdate(oldValueMap))
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "travel/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *TravelController) List() {
	item := s.Travel
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := expandTravelList(resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}
func expandTravelList(resultMaps []map[string]interface{})[]map[string]interface{}{
	for idx, travelMap := range resultMaps{
		resultMaps[idx]=expandTravelMap(travelMap)
	}
	return resultMaps
}
func expandTravelMap(travelMap map[string]interface{})map[string]interface{}{
	if sn, ok := travelMap[s.Approver]; ok{
		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
			name, _ := userMap[s.Name]
			travelMap[s.ApproverSn]= travelMap[s.Approver]
			travelMap[s.Approver] = name
		}
	}
	if sn, ok := travelMap[s.Traveler];ok{
		if status, userMap := userMgr.Get(sn.(string)); status == stat.Success{
			name, _ := userMap[s.Name]
			travelMap[s.TravelerSn] = travelMap[s.Traveler]
			travelMap[s.Traveler] = name
		}
	}
	return travelMap
}
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

func expandTraveDef(oItemDef itemDef.ItemDef)itemDef.ItemDef{
	approverField, _ := oItemDef.GetField(s.Approver)
	approverField.Name = s.ApproverName
	approverSn, _ := oItemDef.GetField(s.Approver)
	approverSn.Input = s.Hidden
	newFields := make([]itemDef.Field, 0)
	for _, field := range oItemDef.Fields[:len(oItemDef.Fields)-3]{
		newFields = append(newFields, field)
	}
	newFields = append(newFields, approverField)
	newFields = append(newFields, approverSn)
	oItemDef.Fields = newFields
	return oItemDef
}