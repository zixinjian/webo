package routers

import (
	"erp/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/appliances/ui/list", &controllers.AppliancesController{}, "*:UiList")
	beego.Router("/appliances/item/list", &controllers.AppliancesController{}, "*:ItemList")
	beego.Router("/appliances/ui/improperlist", &controllers.AppliancesController{}, "*:UiImproperList")
	beego.Router("/appliances/ui/planlist", &controllers.AppliancesController{}, "*:UiPlanlist")
	beego.Router("/appliances/ui/add", &controllers.AppliancesController{}, "*:UiAdd")
	beego.Router("/appliances/ui/update", &controllers.AppliancesController{}, "*:UiUpdate")
	beego.Router("/appliances/ui/improperupdate", &controllers.AppliancesController{}, "*:UiImproperUpdate")
	beego.Router("/appliances/ui/improperadd", &controllers.AppliancesController{}, "*:UiImproperAdd")
	beego.Router("/appliances/ui/mobilescan", &controllers.AppliancesController{}, "*:UiMobilescan")
	beego.Router("/appliances/ui/modellist", &controllers.AppliancesController{}, "*:UiModelList")

	beego.Router("/appliances/ui/recordinlist", &controllers.AppliancesController{}, "*:UiRecordinList")
	beego.Router("/appliances/ui/recordoutlist", &controllers.AppliancesController{}, "*:UiRecordoutList")
	beego.Router("/appliances/ui/recordreturnlist", &controllers.AppliancesController{}, "*:UiRecordreturnList")

	beego.Router("/appliances/ui/imoutlist", &controllers.AppliancesController{}, "*:UiImoutList")
	beego.Router("/appliances/ui/imreturnlist", &controllers.AppliancesController{}, "*:UiImreturnList")
	beego.Router("/appliances/ui/imhistorylist", &controllers.AppliancesController{}, "*:UiImhistoryList")

	beego.Router("/appliances/ui/cometolist", &controllers.AppliancesController{}, "*:UiCometoList")
	beego.Router("/appliances/ui/imcomelist", &controllers.AppliancesController{}, "*:UiImComeList")

	beego.Router("/appliances/ui/outhistory", &controllers.AppliancesController{}, "*:UiOutHistory")
	beego.Router("/appliances/ui/return/history", &controllers.AppliancesController{}, "*:UiReturnHistory")
	beego.Router("/appliances/item/deploy/products", &controllers.AppliancesController{}, "*:DeployProducts")

	beego.Router("/page/ui/recordout/list", &controllers.AppliancesController{}, "*:UiPageRecordoutList")
}
