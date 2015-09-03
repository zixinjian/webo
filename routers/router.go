package routers

import (
	"github.com/astaxie/beego"
	"webo/controllers"
)

func init() {
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/main", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/service", &controllers.ServiceController{})
	beego.Router("/item/get/?:id", &controllers.ItemController{}, "*:Get")
	beego.Router("/item/list/:hi:string", &controllers.ItemController{}, "*:List")
	beego.Router("/item/add/:hi:string", &controllers.ItemController{}, "*:Add")
	beego.Router("/item/update/:hi:string", &controllers.ItemController{}, "*:Update")
	beego.Router("/item/delete/:hi:string", &controllers.ItemController{}, "*:Delete")
	beego.Router("/ui/add/:hi:string", &controllers.UiController{}, "*:Add")
	beego.Router("/ui/list/:hi:string", &controllers.UiController{}, "*:List")
	beego.Router("/ui/update/:hi:string", &controllers.UiController{}, "*:Update")

	beego.Router("/ui/purchase/mycreate", &controllers.PurchaseController{}, "*:UiMyCreate")
	beego.Router("/ui/purchase/curlist", &controllers.PurchaseController{}, "*:UiCurList")
	beego.Router("/ui/purchase/history", &controllers.PurchaseController{}, "*:UiHistoryList")
	beego.Router("/travel", &controllers.MainController{}, "*:Travel")

	beego.Router("ui/user/list", &controllers.UserController{}, "*:UiList")
	beego.Router("/user/disable", &controllers.UserController{}, "*:Disable")
}
