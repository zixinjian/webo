package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"webo/models/s"
)

type MainController struct {
	BaseController
}

const userMgrHtml =`<li>
	<a href="/ui/user/list" target="main" class="auto">
		<i class="icon-users icon text-success-lter"></i>
		<span>用户</span>
	</a>
</li>
`
const managerNavHtml = `<li class="active">
	<a href="/ui/purchase/mycreate" target="main" class="auto">
		<i class="icon-note text-primary-dker"></i>
		<span>我创建订单</span>
	</a>
</li>
<li>
	<a href="/ui/purchase/curlist" target="main" class="auto">
		<i class="icon icon-basket-loaded"></i>
		<span>待处理的订单</span>
	</a>
</li>
`
const userNavHtml = `<li class="active">
	<a href="/ui/purchase/curlist" target="main" class="auto">
		<i class="glyphicon glyphicon-align-justify icon text-primary-dker"></i>
		<span>待处理的订单</span>
	</a>
</li>
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
		this.Data["activeUrl"] = "/ui/purchase/mycreate"
	}else {
		this.Data["orderNav"] = userNavHtml
		this.Data["userMgr"] = ""
		this.Data["activeUrl"] = "/ui/purchase/curlist"
	}
	this.Data["userName"] = userName
	this.TplNames = "main.html"
}


