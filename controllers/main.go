package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"webo/models/s"
)

type MainController struct {
	BaseController
}

const userMgrHtml = `<ul class="nav nav-sidebar">
	<li><a href="/ui/user/list" target="frame-content">用户管理</a></li>
</ul>
`
const managerNavHtml = `<li class="active"><a href="/ui/purchase/mycreate" target="frame-content">创建订单<span class="sr-only"></span></a></li>
<li><a href="/ui/purchase/curlist" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
const userNavHtml = `<li class="active"><a href="/ui/purchase/curlist" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
const activeUrlFormat = `<iframe name = "frame-content" src="%s" layout-auto-height="-20" style="width:100%%;border:none"></iframe>
`

func (this *MainController) Get() {
	this.SetSession(SessionUserName, "admin")
	this.SetSession(SessionUserRole, s.RoleAdmin)
	this.SetSession(SessionUserSn, "20150729203140000")
	this.SetSession(SessionUserDepartment, "department")
	userName := this.GetCurUser()
	userRole := this.GetCurRole()
	beego.Info(fmt.Sprintf("User:%s login as role:%s", userName, userRole))
	this.Data["userName"] = userName
	switch userRole {
	case "role_admin", "role_manager":
		beego.Debug("Show role", userRole)
		this.Data["orderNav"] = managerNavHtml
		this.Data["userMgr"] = userMgrHtml
		this.Data["activeUrl"] = fmt.Sprintf(activeUrlFormat, "/ui/purchase/mycreate")

	default:
		this.Data["orderNav"] = userNavHtml
		this.Data["userMgr"] = ""
		this.Data["activeUrl"] = fmt.Sprintf(activeUrlFormat, "/ui/purchase/mycreate")
	}
	this.TplNames = "main.html"
}

func (this *MainController) Travel() {
	this.SetSession(SessionUserName, "admin")
	this.SetSession(SessionUserRole, "role_admin")
	this.SetSession(SessionUserSn, "snlsnsldn")
	this.SetSession(SessionUserDepartment, "department")
	userName := this.GetCurUser()
	userRole := this.GetCurRole()
	beego.Info(fmt.Sprintf("User:%s login as role:%s", userName, userRole))
	this.Data["userName"] = userName
	switch userRole {
	case "role_admin", "role_manager":
		beego.Debug("Show role", userRole)
		this.Data["orderNav"] = managerNavHtml
		this.Data["userMgr"] = userMgrHtml

	default:
		this.Data["orderNav"] = userNavHtml
		this.Data["userMgr"] = ""
	}
	this.TplNames = "travel.html"
}
