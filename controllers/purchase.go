package controllers

import (
	"fmt"
	"webo/controllers/ui"
	"github.com/astaxie/beego"
	"webo/models/itemDef"
	"webo/models/svc"
	"webo/models/u"
	"webo/models/stat"
	"webo/models/s"
	"strings"
	"webo/models/purchaseMgr"
	"webo/models/t"
	"webo/models/supplierMgr"
	"webo/models/productMgr"
)

type PurchaseController struct {
	ItemController
}

func (this *PurchaseController) UiMyCreate() {
	item:="purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/item/list/purchase?creater=curuser&godowndate"
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
<input data-model="buyers" type="radio" name = "buyers" id="%s" value="%s" %s> %s
</label>
`
const AdminUserFormat =`<label class="radio-inline">
                <input data-model="buyers" type="radio" name = "buyers" id="all" value="all" %s> 全部
            </label>
`
func (this *PurchaseController) UiCurList() {
	beego.Info("UiCurList")
	item:="purchase"
	oItemDef, ok := itemDef.EntityDefMap[item]
	if !ok {
		beego.Error(fmt.Sprintf("Item %s not support", item))
	}
	queryParam := t.Params{
		"department":"department_purchase",
	}
	_, _, retMaps := svc.List(s.User, queryParam, t.LimitParams{}, t.Params{})

	var allCheked string
	role := this.GetCurRole()
	switch role {
	case s.RoleManager, s.RoleAdmin:
		allCheked = "checked"
	default:
		allCheked = ""
	}
	userHtml := fmt.Sprintf(AdminUserFormat, allCheked)
	curUserSn := this.GetCurUserSn()
	for _, userMap := range retMaps{
		sn, sok := userMap["sn"]
		name, nok := userMap["name"]
		if !(sok && nok){
			continue
		}
		checked := ""
		if u.IsNullStr(allCheked) && strings.EqualFold(curUserSn, sn.(string)){
			checked = "checked"
		}
		userHtml = userHtml + fmt.Sprintf(buyerHtmlFormat, sn, sn, checked, name)
	}
	this.Data["buyers"] = userHtml
	this.Data["queryParams"] = CurListQueryParamsJs
	this.Data["listUrl"] = "/item/list/purchase?godowndate"
	this.Data["addUrl"] = ""
	this.Data["updateUrl"] = "/ui/purchase/userupdate"
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
	this.Data["listUrl"] = "/item/list/purchase"
	this.Data["updateUrl"] = "/ui/purchase/show"
	this.Data["queryParams"] = HistoryListQueryParamsJs
}

func (this *PurchaseController) UiAdd() {
	item:="purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	addItemDef := fillBuyerEnum(getAddPurchaseDef(oItemDef))
	this.Data["Service"] = "/item/add/" + item
	this.Data["Form"] = ui.BuildAddForm(addItemDef, u.TUId())
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
	queryParams := t.Params{
		s.Department:"department_purchase",
	}
	orderParams := t.Params{
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
	statusMap := map[string]string{}
	this.UiUpdateWithStatus(statusMap)
}
func (this *PurchaseController) UiUserUpdate() {
	statusMap := map[string]string{
		s.Category:s.Disabled,
		s.Product:s.Disabled,
		s.Model:s.Disabled,
		s.PlaceDate:s.Disabled,
		s.Requireddate:s.Disabled,
		s.Requireddepartment:s.Disabled,
	}
	this.UiUpdateWithStatus(statusMap)
}
func (this *PurchaseController) UiHistoryUpdate() {
	item:=s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	statusMap := make(map[string]string, len(oItemDef.Fields))
	for _, field := range oItemDef.Fields{
		statusMap[field.Name] = s.Disabled
	}
	this.UiUpdateWithStatus(statusMap)
}

func (this *PurchaseController) UiUpdateWithStatus(statusMap map[string]string) {
	item:=s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	sn := this.GetString(s.Sn)
	if sn == "" {
		this.Ctx.WriteString(stat.ParamSnIsNone)
		return
	}
	params := t.Params{s.Sn: sn}
	code, oldValueMap := svc.Get(item, params)
	oldValueMap = expandPurchaseMap(oldValueMap)
	if code == stat.Success {
		this.Data["Service"] = "/item/update/" + item
		oItemDef = fillBuyerEnum(oItemDef)
		this.Data["Form"] = ui.BuildUpdatedFormWithStatus(oItemDef, oldValueMap, statusMap)
		this.Data["Onload"] = ui.BuildUpdateOnLoadJs(oItemDef)
		this.TplNames = "purchase/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *PurchaseController) List() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := purchaseMgr.GetPurchases(queryParams, limitParams, orderByParams)
	this.Data["json"] = &TableResult{result, int64(total), resultMaps}
	this.ServeJson()
}

func expandPurchaseMap(oldMap t.ItemMap)t.ItemMap{
	var retMap = make(t.ItemMap, 0)
	for key, value := range oldMap {
		retMap[strings.ToLower(key)] = value
	}
	if userName, ok := oldMap["user_name"]; ok {
		retMap["buyer"] = userName
	}
	if supplierSn, ok := retMap[s.Supplier];ok && !u.IsNullStr(supplierSn){
		if supplierMap, sok := supplierMgr.Get(supplierSn.(string)); sok {
			supplierKey, _ := supplierMap[s.Keyword]
			supplierName, _:= supplierMap[s.Name]
			retMap[s.Supplier+ s.EKey] = supplierKey.(string)
			retMap[s.Supplier+ s.EName] = supplierName.(string)
			retMap[s.Supplier] = supplierSn
		}
	}
	if productSn, ok := retMap[s.Product];ok && !u.IsNullStr(productSn){
		if productMap, sok := productMgr.Get(productSn.(string));sok{
			productKey, _ := productMap[s.Keyword]
			productName, _:= productMap[s.Name]
			retMap[s.Product+ s.EKey] = productKey.(string)
			retMap[s.Product+ s.EName] = productName.(string)
			retMap[s.Product] = productSn
		}
	}
	return retMap
}