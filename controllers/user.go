package controllers

import (
	"fmt"
	"strings"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/svc"
)

type UserController struct {
	BaseController
}

func (this *UserController) UiList() {
	item := "user"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "user/list.tpl"
}

func (this *UserController) Update() {
	item := s.User
	oEntityDef, _ := itemDef.EntityDefMap[item]
	svcParams := this.GetFormValues(oEntityDef)
	if pwd, ok := svcParams[s.Password]; ok {
		if strings.EqualFold(pwd.(string), "*****") {
			delete(svcParams, s.Password)
		}
	}
	status, reason := svc.Update(item, svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}
//func (this *UserController) Disable() {
//	role := this.GetSessionString(SessionUserRole)
//	if role != s.RoleAdmin{
//		this.Data["json"] = &JsonResult{status.PermissionDenied, status.PermissionDenied}
//		this.ServeJson()
//		return
//	}
//	sn := this.GetStrings(s.Sn)
//	beego.Info("Disable user sn:", sn)
//	svcParams := svc.Params{
//		s.Sn : sn,
//		s.Flag : "flag_disable",
//	}
//	status, reason := svc.Update("user", svcParams)
//	this.Data["json"] = &JsonResult{status, reason}
//	this.ServeJson()
//}
