package controllers

import (
	"fmt"
	"webo/models/s"
)

type TravelController struct {
	BaseController
}

func (this *TravelController) UiList() {
	item := s.Travel
	this.Data["item"] = item
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.TplNames = "travel/list.html"
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
