package controllers

import (
"fmt"
"webo/controllers/ui"
//	"webo/models/itemDef"
"github.com/astaxie/beego"
"webo/models/itemDef"
	"webo/models/svc"
	"webo/models/s"
	"webo/models/status"
)

type UserController struct {
	BaseController
}

func (this *UserController) UiList() {
	item := "user"
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "user/list.html"
}

func (this *UserController) Disable() {
	role := this.GetSessionString(SessionUserRole)
	if role != s.RoleAdmin{
		this.Data["json"] = &JsonResult{status.PermissionDenied, status.PermissionDenied}
		this.ServeJson()
		return
	}
	sn := this.GetStrings(s.Sn)
	beego.Info("Disable user sn:", sn)
	svcParams := svc.Params{
		s.Sn : sn,
		s.Flag : "flag_disable",
	}
	status, reason := svc.Update("user", svcParams)
	this.Data["json"] = &JsonResult{status, reason}
	this.ServeJson()
}