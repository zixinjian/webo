package routers

import (
	"erp/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/agreement/ui/add", &controllers.AgreementController{}, "*:UiAdd")
	beego.Router("/agreement/item/add", &controllers.AgreementController{}, "*:ItemAdd")
	beego.Router("/agreement/ui/update", &controllers.AgreementController{}, "*:UiUpdate")
	beego.Router("/agreement/ui/receiptlist", &controllers.AgreementController{}, "*:UiReceiptList")
	beego.Router("/agreement/ui/totalslist", &controllers.AgreementController{}, "*:UiTotalsList")
	beego.Router("/agreement/ui/recordinglist", &controllers.AgreementController{}, "*:UiRecordingList")
	beego.Router("/agreement/ui/look", &controllers.AgreementController{}, "*:UiLook")
	beego.Router("/agreement/ui/list", &controllers.AgreementController{}, "*:UiList")

	beego.Router("/agreement/ui/approval/list", &controllers.AgreementController{}, "*:UiApprovalList")
	beego.Router("/agreement/ui/npapproval/list", &controllers.AgreementController{}, "*:UiNpapprovalList")
	beego.Router("/agreement/ui/ag/list", &controllers.AgreementController{}, "*:UiAgList")
	beego.Router("/agreement/ui/print", &controllers.AgreementController{}, "*:UiPrintList")
	beego.Router("/agreement/ui/approvalprint", &controllers.AgreementController{}, "*:UiApprovalPrintList")
	beego.Router("/agreement/ui/npapprovalprint", &controllers.AgreementController{}, "*:UiNpapprovalPrintList")
	beego.Router("/agreement/ui/deploy/update", &controllers.DeployrecordController{}, "*:UiUpdate")
	beego.Router("/agreement/ui/atomiserpo/update", &controllers.AtomiserPoController{}, "*:UiAtomiserPoUpdate")
	beego.Router("/agreement/ui/contract/list", &controllers.AgreementController{}, "*:UiContractList")
	beego.Router("/agreement/item/ag/update", &controllers.AgreementController{}, "*:UpdateAgList")
	beego.Router("/agreement/item/list", &controllers.AgreementController{}, "*:ItemList")
	beego.Router("/agreement/item/create", &controllers.AgreementController{}, "*:CreateStatement")

	beego.Router("/agreement/ui/inbill/list", &controllers.AgreementController{}, "*:UiInBillList") // 对内对账单
	beego.Router("/agreement/ui/onbill/list", &controllers.AgreementController{}, "*:UiOnBillList") //对外对账单
}
