package routers

import (
	"github.com/astaxie/beego"
	"webo/controllers"
)

func init() {
	//框架服务
	beego.Router("/", &controllers.LoginController{})
	beego.Router("/main", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/service", &controllers.ServiceController{})

	//基本资源服务
	beego.Router("/item/get/?:id", &controllers.ItemController{}, "*:Get")
	beego.Router("/item/list/:hi:string", &controllers.ItemController{}, "*:List")
	beego.Router("/item/add/:hi:string", &controllers.ItemController{}, "*:Add")
	beego.Router("/item/update/:hi:string", &controllers.ItemController{}, "*:Update")
	beego.Router("/item/delete/:hi:string", &controllers.ItemController{}, "*:Delete")
	beego.Router("/item/upload/:hi:string", &controllers.ItemController{}, "*:Upload")
	beego.Router("/item/autocomplete/:hi:string", &controllers.ItemController{}, "*:Autocomplete")

	//基本资源页面
	beego.Router("/ui/add/:hi:string", &controllers.UiController{}, "*:Add")
	beego.Router("/ui/list/:hi:string", &controllers.UiController{}, "*:List")
	beego.Router("/ui/update/:hi:string", &controllers.UiController{}, "*:Update")

	//用户管理
	beego.Router("/ui/user/list", &controllers.UserController{}, "*:UiList")
	beego.Router("/item/update/user", &controllers.UserController{}, "*:Update")

	//供应商
	beego.Router("/ui/add/supplier", &controllers.SupplierController{}, "*:UiAdd")
	beego.Router("/ui/update/supplier", &controllers.SupplierController{}, "*:UiUpdate")

	//产品管理
	beego.Router("/ui/product/add", &controllers.ProductController{}, "*:UiAdd")
	beego.Router("/ui/product/setting", &controllers.ProductController{}, "*:UiSetting")
	beego.Router("/ui/product/list", &controllers.ProductController{}, "*:UiList")
	beego.Router("/ui/product/update", &controllers.ProductController{}, "*:UiUpdate")
	beego.Router("/item/product/list", &controllers.ProductController{}, "*:List")
	beego.Router("/item/product/add", &controllers.ProductController{}, "*:Add")
	beego.Router("/item/product/update", &controllers.ProductController{}, "*:Update")
	beego.Router("/item/add/catagory", &controllers.ConfigController{}, "*:List")
	beego.Router("/item/update/catagory", &controllers.ConfigController{}, "*:Update")

	//采购管理
	beego.Router("/ui/purchase/mycreate", &controllers.PurchaseController{}, "*:UiMyCreate")
	beego.Router("/ui/purchase/curlist", &controllers.PurchaseController{}, "*:UiCurList")
	beego.Router("/ui/purchase/history", &controllers.PurchaseController{}, "*:UiHistoryList")
	beego.Router("/ui/purchase/add", &controllers.PurchaseController{}, "*:UiAdd")
	beego.Router("/ui/purchase/update", &controllers.PurchaseController{}, "*:UiUpdate")
	beego.Router("/ui/purchase/userupdate", &controllers.PurchaseController{}, "*:UiUserUpdate")
	beego.Router("/ui/purchase/show", &controllers.PurchaseController{}, "*:UiHistoryUpdate")
	beego.Router("/item/list/purchase", &controllers.PurchaseController{}, "*:List")


	//分析
	beego.Router("/ui/purchase/priceAnalyze", &controllers.PurchaseController{}, "*:PriceAnalyze")
	beego.Router("/ui/purchase/buyertimely", &controllers.PurchaseController{}, "*:BuyerTimely")
	beego.Router("/ui/purchase/producttimely", &controllers.PurchaseController{}, "*:ProductTimely")
	beego.Router("/ui/purchase/suppliertimely", &controllers.PurchaseController{}, "*:SupplierTimely")
	beego.Router("/purchase/list/buyertimely", &controllers.PurchaseController{}, "*:BuyerTimelyList")
	beego.Router("/purchase/calc/producttimely", &controllers.PurchaseController{}, "*:CalcProductTimely")
	beego.Router("/purchase/list/suppliertimely", &controllers.PurchaseController{}, "*:SupplierTimelyList")

//	beego.Router("/ui/setting", &controllers.SettingController{}, "*:List")

	//账目
	beego.Router("/ui/expense/list", &controllers.PurchaseController{}, "*:ExpenseList")
	beego.Router("/ui/expense/accountcurrentlist", &controllers.PurchaseController{}, "*:AccountCurrentList")

	//出差管理
	beego.Router("/travel", &controllers.TravelController{}, "*:Travel")

	beego.Router("/travel/ui/myCreate", &controllers.TravelController{}, "*:MyCreate")
	beego.Router("/travel/ui/list", &controllers.TravelController{}, "*:UiList")
	beego.Router("/travel/ui/update", &controllers.TravelController{}, "*:UiUpdate")
	beego.Router("/travel/ui/add", &controllers.TravelController{}, "*:UiAdd")
	beego.Router("/travel/item/list", &controllers.TravelController{}, "*:List")

	beego.Router("/account", &controllers.AccountController{}, "*:Account")
	beego.Router("/account/ui/add", &controllers.AccountController{}, "*:UiAdd")
	beego.Router("/account/ui/list", &controllers.AccountController{}, "*:UiList")
	beego.Router("/account/item/list", &controllers.AccountController{}, "*:List")
}
