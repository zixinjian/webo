package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
)

type PurchaseController struct {
	ItemController
}

func (this *PurchaseController) UiMyCreate() {
	item:="purchase"
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s?creater=curuser", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "item/list.html"
}
const CurListQueryParamsJs = `<script>
    function queryParams(params){
        params["godowndate"]=""
        return params
    }
</script>
`
const buyerHtmlFormat =`<label class="radio-inline">
<input data-model="buyers" type="radio" name = "buyers" id="%s" value="%s"> %s
</label>
`
func (this *PurchaseController) UiCurList() {
	beego.Info("UiCurList")
	item:="purchase"
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	queryParam := svc.Params{
		"department":"department_purchase",
	}
	_, _, retMaps := svc.List("user", queryParam, svc.LimitParams{}, svc.Params{})
	userHtml := ""
	for _, userMap := range retMaps{
		sn, sok := userMap["sn"]
		name, nok := userMap["name"]
		if !(sok && nok){
			continue
		}
		userHtml = userHtml + fmt.Sprintf(buyerHtmlFormat, sn, sn, name)
	}
	this.Data["buyers"] = userHtml
	this.Data["queryParams"] = CurListQueryParamsJs
	this.Data["listUrl"] = fmt.Sprintf("/item/list/%s", item)
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = fmt.Sprintf("/ui/update/%s", item)
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "purchase/list.html"
}

const HistoryListQueryParamsJs = `<script>
    function queryParams(params){
        return params
    }
</script>
`

func (this *PurchaseController) UiHistoryList() {
	this.UiCurList()
	this.Data["queryParams"] = HistoryListQueryParamsJs
}
