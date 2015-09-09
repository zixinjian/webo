package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"webo/models/itemDef"
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
	this.TplNames = "user/list.html"
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