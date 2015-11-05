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
const managerNavHtml = `<li class="active"><a href="/ui/purchase/mycreate" target="frame-content">我创建订单<span class="sr-only"></span></a></li>
<li><a href="/ui/purchase/curlist" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
const userNavHtml = `<li class="active"><a href="/ui/purchase/curlist" target="frame-content">待处理的订单<span class="sr-only"></span></a></li>
`
const activeUrlFormat = `<iframe name = "frame-content" src="%s" layout-auto-height="-20" style="width:100%%;border:none"></iframe>
`

func (this *MainController) Get() {
	userName := this.GetCurUserUserName()
	userRole := this.GetCurRole()
	beego.Info(fmt.Sprintf("User:%s login as role:%s", userName, userRole))
	department := this.GetCurDepartment()
	if userRole == s.RoleManager && department == "department_purchase"{
		this.Data["orderNav"] = managerNavHtml
		this.Data["userMgr"] = userMgrHtml
		this.Data["activeUrl"] = fmt.Sprintf(activeUrlFormat, "/ui/purchase/mycreate")
	}else {
		this.Data["orderNav"] = userNavHtml
		this.Data["userMgr"] = ""
		this.Data["activeUrl"] = fmt.Sprintf(activeUrlFormat, "/ui/purchase/curlist")
	}
	this.Data["userName"] = userName
	this.TplNames = "main.html"
}


