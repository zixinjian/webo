package controllers
import (
	"webo/models/u"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/s"
	"github.com/astaxie/beego"
	"webo/models/stat"
	"webo/models/t"
	"webo/models/svc"
)


type SupplierController struct {
	UiController
}

func (this *SupplierController) UiAdd() {
	item := s.Supplier
	oItemDef, _ := itemDef.EntityDefMap[item]
	oItemDef.FillEnum(s.Category, getCategoryEnumList())
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(oItemDef, u.TUId())
	this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
	this.TplNames = "item/add.tpl"
}

func (this *SupplierController) UiUpdate() {
	item := s.Supplier
	oItemDef, _ := itemDef.EntityDefMap[item]
	sn := this.GetString(s.Sn)
	if sn == "" {
		beego.Error("ui.Update", stat.ParamSnIsNone)
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := t.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		this.Data["Service"] = "/item/update/" + item
		oItemDef.FillEnum(s.Category, getCategoryEnumList())
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "item/update.tpl"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}