package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
	"webo/models/util"
	"webo/models/stat"
	"webo/models/s"
	"strings"
	"encoding/json"
	"webo/models/purchase"
)

type PurchaseController struct {
	ItemController
}

func (this *PurchaseController) UiMyCreate() {
	item:="purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/ui/purchase/list?creater=curuser"
	this.Data["addUrl"] = "/ui/purchase/add"
	this.Data["updateUrl"] = "/ui/purchase/update"
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "purchase/mycreates.html"
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
	this.Data["listUrl"] = "/ui/purchase/list"
	this.Data["addUrl"] = fmt.Sprintf("/ui/add/%s", item)
	this.Data["updateUrl"] = "/ui/purchase/update"
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

func (this *PurchaseController) UiAdd() {
	item:="purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	addItemDef := fillBuyerEnum(getAddPurchaseDef(oItemDef))
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(addItemDef, util.TUId())
	this.Data["Onload"] = ui.BuildAddOnLoadJs(addItemDef)
	this.TplNames = "purchase/add.tpl"
}

func getAddPurchaseDef(oItemDef itemDef.ItemDef)itemDef.ItemDef{
	names := []string{"category", s.Name,  s.Model, "buyer", s.Num, "placedate", "requireddate", "requireddepartment", s.Mark}
	fields := make([]itemDef.Field, len(names))

	fieldMap := oItemDef.GetFieldMap()
	for idx, name := range names{
		if field, ok := fieldMap[name];ok {
			fields[idx] = field
		}else{
			beego.Error("Field not found", name)
		}
	}
	oItemDef.Fields = fields
	return oItemDef
}

func fillBuyerEnum(oItemDef itemDef.ItemDef) itemDef.ItemDef{
	for idx, field := range oItemDef.Fields{
		if strings.EqualFold("buyer", field.Name){
			field.Enum = getBuyerEnum()
		}
		oItemDef.Fields[idx] = field
	}
	return oItemDef
}

func getBuyerEnum()[]itemDef.EnumValue{
	queryParams := svc.Params{
		s.Department:"department_purchase",
	}
	orderParams := svc.Params{
		s.Name : s.Asc,
	}
	if code, userMaps := svc.GetItems(s.User, queryParams, orderParams); strings.EqualFold(code, stat.Success){
		EnumList := make([]itemDef.EnumValue, len(userMaps))
		for idx, user := range userMaps{
			v, _ := user[s.Sn]
			u, _ := user[s.UserName]
			l, _ := user[s.Name]
			EnumList[idx]=itemDef.EnumValue{v.(string), u.(string), l.(string)}
		}
		return EnumList
	}else{
		return make([]itemDef.EnumValue, 0)
	}
}

func (this *PurchaseController) UiUpdate() {
	item:="purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	sn := this.GetString("sn")
	if sn == "" {
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := svc.Params{"sn": sn}
	code, oldValueMap := svc.Get(item, params)
	if code == "success" {
		this.Data["Service"] = "/item/update/" + item
		oItemDef = fillBuyerEnum(oItemDef)
		this.Data["Form"] = ui.BuildUpdatedForm(oItemDef, oldValueMap)
		this.Data["Onload"] = ui.BuildAddOnLoadJs(oItemDef)
		this.TplNames = "purchase/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *PurchaseController) ListWithQuery(oItemDef itemDef.ItemDef, addQueryParam svc.Params) {
	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("ListWithQuery requestMap: ", requestMap)

	limitParams := this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams := this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	queryParams := this.GetQueryParamFromJsonMap(requestMap, oItemDef)
	for k, v := range addQueryParam{
		queryParams[k]=v
	}

	result, total, resultMaps := svc.List(oItemDef.Name, queryParams, limitParams, orderByParams)
	retList := TransList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}

func (this *PurchaseController) List() {
	item := s.Purchase

	requestBody := this.Ctx.Input.RequestBody
	var requestMap map[string]interface{}
	json.Unmarshal(requestBody, &requestMap)
	beego.Debug("PurchaseController.List requestMap: ", requestMap)

	limitParams := this.GetLimitParamFromJsonMap(requestMap)
	delete(requestMap, s.Limit)
	delete(requestMap, s.Offset)

	orderByParams := this.GetOrderParamFromJsonMap(requestMap)
	delete(requestMap, s.Order)
	delete(requestMap, s.Sort)

	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams := this.GetQueryParamFromJsonMap(requestMap, oItemDef)
	addParams := this.GetFormValues(oItemDef)
	for k, v := range addParams{
		queryParams[k]=v
	}

	result, total, resultMaps := purchase.GetPurchases(queryParams, limitParams, orderByParams)
	retList := TransList(oItemDef, resultMaps)
	this.Data["json"] = &TableResult{result, int64(total), retList}
	this.ServeJson()
}