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
	"encoding/json"
)

type PurchaseController struct {
	ItemController
}

const buyerHtmlFormat = `<label class="radio-inline">
<input data-model="buyers" type="radio" name = "buyers" id="%s" value="%s" %s> %s
</label>
`
const AdminUserFormat = `<label class="radio-inline">
                <input data-model="buyers" type="radio" name = "buyers" id="all" value="all" %s> 全部
            </label>
`
const CurListQueryParamsJs = `<script>
    function queryParams(params){
        params["godowndate"]=""
        return params
    }
</script>
`
const HistoryListQueryParamsJs = `<script>
    function queryParams(params){
        return params
    }
</script>
`

// 我创建的订单列表
func (this *PurchaseController) UiMyCreate() {
	item := "purchase"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["listUrl"] = "/item/list/purchase?creater=curuser&godowndate"
	this.Data["addUrl"] = "/ui/purchase/add"
	this.Data["updateUrl"] = "/ui/purchase/update"
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.TplNames = "purchase/mycreates.html"
}

//待处理的订单列表
func (this *PurchaseController) UiCurList() {
	beego.Info("UiCurList")
	item := s.Purchase
	this.Data["buyers"] = this.createBuyerList()
	this.Data["queryParams"] = CurListQueryParamsJs
	role := this.GetCurRole()
	department := this.GetCurDepartment()
	if role == s.RoleManager && department == "department_purchase"{
		this.Data["listUrl"] = "/item/list/purchase?godowndate"
	}else {
		this.Data["listUrl"] = "/item/list/purchase?buyer=curuser&godowndate"
	}
	this.Data["addUrl"] = ""
	this.Data["updateUrl"] = "/ui/purchase/userupdate"
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["thlist"] = ui.BuildListThs(oItemDef)
	this.Data["sortOrder"] = s.Asc
	this.TplNames = "purchase/list.html"
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
	this.Data["updateUrl"] = "/ui/purchase/show"
	this.TplNames = "purchase/priceanalyze.tpl"
}

// 添加
func (this *PurchaseController) UiAdd() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	addItemDef := fillBuyerEnum(getAddPurchaseDef(oItemDef))
	this.Data["Service"] = "/item/add/" + item
	statusMap := map[string]string{
		s.ProductPrice: s.ReadOnly,
		s.Power:s.Disabled,
	}
	this.Data["Form"] = ui.BuildAddFormWithStatus(addItemDef, u.TUId(), statusMap)
	this.Data["Onload"] = ui.BuildAddOnLoadJs(addItemDef)
	this.TplNames = "purchase/add.tpl"
}

//管理者修改
func (this *PurchaseController) UiUpdate() {
	statusMap := map[string]string{
		s.PaymentAmount:	  s.ReadOnly,
		s.PaymentDate:	  	  s.ReadOnly,
	}
	this.UiUpdateWithStatus(statusMap)
}

//用户修改
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
		s.Power:       		  s.Disabled,
		s.PaymentAmount:	  s.ReadOnly,
		s.PaymentDate:	  	  s.ReadOnly,
	}
	this.UiUpdateWithStatus(statusMap)
}

//历史中修改
func (this *PurchaseController) UiHistoryUpdate() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	statusMap := make(map[string]string, len(oItemDef.Fields))
	for _, field := range oItemDef.Fields {
		statusMap[field.Name] = s.Disabled
	}
	this.UiUpdateWithStatus(statusMap)
}

const oldValueFormat = `<script>
    var oldValue = %s
</script>
`
//修改基本方法
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
		if purchaseStr, err := json.Marshal(oldValueMap);err == nil{
			this.Data["PuchaseItem"] = fmt.Sprintf(oldValueFormat, string(purchaseStr))
		}else{
			this.Data["PuchaseItem"] = fmt.Sprintf(oldValueFormat, "{}")
		}
		this.TplNames = "purchase/update.html"
	} else {
		this.Ctx.WriteString(stat.ItemNotFound)
	}
}

func (this *PurchaseController) ExpenseList() {
	this.UiCurList()
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["thlist"] = ui.BuildListThs(getExpandListDef(oItemDef))
	this.Data["updateUrl"] = "/ui/purchase/show"
	this.TplNames = "purchase/expand.html"
}

func (this *PurchaseController) AccountCurrentList() {
	this.UiCurList()
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	this.Data["thlist"] = ui.BuildListThs(getExpandListDef(oItemDef))
	this.Data["updateUrl"] = "/ui/purchase/show"
	this.TplNames = "purchase/accountcurrent.html"
}

func (this *PurchaseController) createBuyerList() string {
	role := this.GetCurRole()
	department := this.GetCurDepartment()
	if role != s.RoleManager || department != "department_purchase"{
		return ""
	}
	queryParam := t.Params{
		"department": "department_purchase",
	}
	_, _, retMaps := svc.List(s.User, queryParam, t.LimitParams{}, t.Params{})

	allCheked := s.Checked
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
func (this *PurchaseController) BuyerTimely() {
	this.TplNames = "purchase/buyertimely.tpl"
}
func (this *PurchaseController) ProductTimely() {
	this.TplNames = "purchase/producttimely.tpl"
}
func (this *PurchaseController) SupplierTimely() {
	this.TplNames = "purchase/suppliertimely.tpl"
}
func (this *PurchaseController) BuyerTimelyList() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := purchaseMgr.GetBuyerTimelyList(queryParams, limitParams, orderByParams)

	this.Data["json"] = &TableResult{result, int64(total), resultMaps}
	this.ServeJson()
}
func (this *PurchaseController) CalcProductTimely() {
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, _, _ := this.GetParams(oItemDef)
	resultMaps := purchaseMgr.CalcProductTimely(queryParams)
	this.Data["json"] = &resultMaps
	this.ServeJson()
}
func (this *PurchaseController) SupplierTimelyList(){
	item := s.Purchase
	oItemDef, _ := itemDef.EntityDefMap[item]
	queryParams, limitParams, orderByParams := this.GetParams(oItemDef)
	result, total, resultMaps := purchaseMgr.GetSupplierTimelyList(queryParams, limitParams, orderByParams)
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

func fillBuyerEnum(oItemDef itemDef.ItemDef) itemDef.ItemDef {
	queryParams := t.Params{
		s.Department: "department_purchase",
	}
	orderParams := t.Params{
		s.Name: s.Asc,
	}
	return FillUserEnum(s.Buyer, oItemDef, queryParams, orderParams)
}
func getAddPurchaseDef(oItemDef itemDef.ItemDef) itemDef.ItemDef {
	names := []string{s.Sn, s.Category, s.Product, s.Model, s.Power, s.ProductPrice, s.Buyer, s.Num, s.PlaceDate, s.Requireddate, s.Requireddepartment, s.Mark}
	return makeFields(oItemDef, names)
}

func getExpandListDef(oItemDef itemDef.ItemDef) itemDef.ItemDef {
	names := []string{s.Sn, s.Category, s.Product, s.Model, s.Power, s.Num, s.UintPrice, s.ProductPrice, s.TotalPrice, s.Buyer, s.Mark}
	return makeFields(oItemDef, names)
}
