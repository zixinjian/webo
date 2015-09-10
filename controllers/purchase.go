package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/controllers/ui"
	"webo/models/itemDef"
	"webo/models/productMgr"
	"webo/models/purchaseMgr"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/supplierMgr"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
)

type PurchaseController struct {
	ItemController
}

func (this *PurchaseController) UiMyCreate() {
	item := "purchase"
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
const buyerHtmlFormat = `<label class="radio-inline">
<input data-model="buyers" type="radio" name = "buyers" id="%s" value="%s" %s> %s
</label>
`
const AdminUserFormat = `<label class="radio-inline">
                <input data-model="buyers" type="radio" name = "buyers" id="all" value="all" %s> 全部
            </label>
`

//待处理的订单列表
func (this *PurchaseController) UiCurList() {
	beego.Info("UiCurList")
	item := s.Purchase
	this.Data["buyers"] = this.createBuyerList()
	this.Data["queryParams"] = CurListQueryParamsJs
	this.Data["listUrl"] = "/item/list/purchase?godowndate"
	this.Data["addUrl"] = ""
	this.Data["updateUrl"] = "/ui/purchase/userupdate"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.Data["sortOrder"] = s.Asc
	this.TplNames = "purchase/list.html"
}

const HistoryListQueryParamsJs = `<script>
    function queryParams(params){
        return params
    }
</script>
`

func (this *PurchaseController) UiAdd() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	addItemDef := fillBuyerEnum(getAddPurchaseDef(oItemDef))
	this.Data["Service"] = "/item/add/" + item
	statusMap := map[string]string{
		s.ProductPrice: s.ReadOnly,
	}
	this.Data["Form"] = ui.BuildAddFormWithStatus(addItemDef, u.TUId(), statusMap)
	this.Data["Onload"] = ui.BuildAddOnLoadJs(addItemDef)
	this.TplNames = "purchase/add.tpl"
}


func fillBuyerEnum(oItemDef itemDef.ItemDef) itemDef.ItemDef {
	for idx, field := range oItemDef.Fields {
		if strings.EqualFold("buyer", field.Name) {
			field.Enum = getBuyerEnum()
		}
		oItemDef.Fields[idx] = field
	}
	return oItemDef
}

func getBuyerEnum() []itemDef.EnumValue {
	queryParams := t.Params{
		s.Department: "department_purchase",
	}
	orderParams := t.Params{
		s.Name: s.Asc,
	}
	if code, userMaps := svc.GetItems(s.User, queryParams, orderParams); strings.EqualFold(code, stat.Success) {
		EnumList := make([]itemDef.EnumValue, len(userMaps))
		for idx, user := range userMaps {
			v, _ := user[s.Sn]
			u, _ := user[s.UserName]
			l, _ := user[s.Name]
			EnumList[idx] = itemDef.EnumValue{v.(string), u.(string), l.(string)}
		}
		return EnumList
	} else {
		return make([]itemDef.EnumValue, 0)
	}
}

func (this *PurchaseController) UiUpdate() {
	statusMap := map[string]string{}
	this.UiUpdateWithStatus(statusMap)
}
func (this *PurchaseController) UiUserUpdate() {
	statusMap := map[string]string{
		s.Sn:                 s.Disabled,
		s.Category:           s.Disabled,
		s.Product:            s.Disabled,
		s.Model:              s.Disabled,
		s.PlaceDate:          s.Disabled,
		s.Requireddate:       s.Disabled,
		s.Requireddepartment: s.Disabled,
		s.ProductPrice:       s.Disabled,
	}
	this.UiUpdateWithStatus(statusMap)
}
func (this *PurchaseController) UiHistoryUpdate() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	statusMap := make(map[string]string, len(oItemDef.Fields))
	for _, field := range oItemDef.Fields {
		statusMap[field.Name] = s.Disabled
	}
	this.UiUpdateWithStatus(statusMap)
}

func (this *PurchaseController) UiUpdateWithStatus(statusMap map[string]string) {
	item := s.Purchase
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

//历史订单列表
func (this *PurchaseController) UiHistoryList() {
	this.UiCurList()
	this.Data["listUrl"] = "/item/list/purchase"
	this.Data["updateUrl"] = "/ui/purchase/show"
	this.Data["sortOrder"] = s.Desc
	this.Data["queryParams"] = HistoryListQueryParamsJs
}

//历史价格分析
func (this *PurchaseController) PriceAnalyze() {
	this.TplNames = "purchase/priceanalyze.tpl"
}

func (this *PurchaseController) ExpenseList() {
	beego.Info("ExpenseList")
	this.UiCurList()
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]

	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "purchase/expand.html"
}

func (this *PurchaseController) createBuyerList() string {
	queryParam := t.Params{
		"department": "department_purchase",
	}
	_, _, retMaps := svc.List(s.User, queryParam, t.LimitParams{}, t.Params{})

	var allCheked string
	role := this.GetCurRole()
	switch role {
	case s.RoleManager, s.RoleAdmin:
		allCheked = s.Checked
	default:
		allCheked = ""
	}
	userHtml := fmt.Sprintf(AdminUserFormat, allCheked)
	curUserSn := this.GetCurUserSn()
	for _, userMap := range retMaps {
		sn, sok := userMap["sn"]
		name, nok := userMap["name"]
		if !(sok && nok) {
			continue
		}
		checked := ""
		if u.IsNullStr(allCheked) && strings.EqualFold(curUserSn, sn.(string)) {
			checked = "checked"
		}
		userHtml = userHtml + fmt.Sprintf(buyerHtmlFormat, sn, sn, checked, name)
	}
	return userHtml
}

func (this *PurchaseController) List() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := purchaseMgr.GetPurchases(queryParams, limitParams, orderByParams)
	this.Data["json"] = &TableResult{result, int64(total), resultMaps}
	this.ServeJson()
}

func expandPurchaseMap(oldMap t.ItemMap) t.ItemMap {
	var retMap = make(t.ItemMap, 0)
	for key, value := range oldMap {
		retMap[strings.ToLower(key)] = value
	}
	if userName, ok := oldMap["user_name"]; ok {
		retMap[s.Buyer] = userName
	}
	if supplierSn, ok := retMap[s.Supplier]; ok && !u.IsNullStr(supplierSn) {
		if supplierMap, sok := supplierMgr.Get(supplierSn.(string)); sok {
			supplierKey, _ := supplierMap[s.Keyword]
			supplierName, _ := supplierMap[s.Name]
			retMap[s.Supplier+s.EKey] = supplierKey.(string)
			retMap[s.Supplier+s.EName] = supplierName.(string)
			retMap[s.Supplier] = supplierSn
		}
	}
	if productSn, ok := retMap[s.Product]; ok && !u.IsNullStr(productSn) {
		if productMap, sok := productMgr.Get(productSn.(string)); sok {
			productKey, _ := productMap[s.Keyword]
			productName, _ := productMap[s.Name]
			ProductPrice, _ := productMap[s.Price]
			retMap[s.Product+s.EKey] = productKey.(string)
			retMap[s.Product+s.EName] = productName.(string)
			retMap[s.Product] = productSn
			retMap[s.ProductPrice] = ProductPrice
		}
	}
	return retMap
}

func getAddPurchaseDef(oItemDef itemDef.ItemDef) itemDef.ItemDef {
	names := []string{s.Sn, s.Category, s.Product, s.Model, s.ProductPrice, s.Buyer, s.Num, s.PlaceDate, s.Requireddate, s.Requireddepartment, s.Mark}
	return makeFields(oItemDef, names)
}

func getExpandListDef(oItemDef itemDef.ItemDef) itemDef.ItemDef{
	names := []string{s.Sn, s.Category, s.Product, s.Model, s.ProductPrice, s.Buyer, s.Num, s.PlaceDate, s.Requireddate, s.Requireddepartment, s.Mark}
	return makeFields(oItemDef, names)
}

