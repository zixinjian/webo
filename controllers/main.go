package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
)

type MainController struct {
	BaseController
}
const userMgrHtml  = `<ul class="nav nav-sidebar">
	<li><a href="/ui/list/user" target="frame-content">用户管理</a></li>
</ul>
`
const managerNavHtml = `<li class="active"><a href="../static/html/supplierList.html" target="frame-content">创建订单<span class="sr-only"></span></a></li>
<li><a href="../static/html/supplierList.html" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
const userNavHtml = `<li class="active"><a href="../static/html/supplierList.html" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
func (this *MainController) Get() {
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