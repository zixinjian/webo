package routers

import (
	"erp/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/account/ui/add", &controllers.AccountController{}, "*:UiAdd")       // 收付款 > 账户 > 添加
	beego.Router("/account/ui/update", &controllers.AccountController{}, "*:UiUpdate") // 收付款 > 账户 > 更新
	beego.Router("/account/ui/list", &controllers.AccountController{}, "*:UiList")     // 收付款 > 账户 > 列表
	beego.Router("/account/item/list", &controllers.AccountController{}, "*:ItemList")
}
