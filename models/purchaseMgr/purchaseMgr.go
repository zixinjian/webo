package purchaseMgr

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"webo/models/productMgr"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/supplierMgr"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
)

const purchaseListSql = "SELECT purchase.*, user.name as user_name, user.username as user_username FROM purchase, user WHERE user.sn = purchase.buyer"
const purchaseCountSql = "SELECT COUNT(purchase.id) as count FROM purchase, user WHERE user.sn = purchase.buyer"

func GetPurchases(queryParams t.Params, limitParams t.LimitParams, orderBy t.Params) (string, int64, []map[string]interface{}) {
	beego.Debug("purchase.GetPurchases:", queryParams, limitParams, orderBy)
	surface := s.Purchase
	sqlBuilder := svc.NewSqlBuilder()
	for k, v := range queryParams {
		sqlBuilder.Filter(surface+"."+k, v)
	}
	count := GetPurchaseTotal(sqlBuilder)
	if count == -1 {
		return stat.Failed, 0, make([]map[string]interface{}, 0)
	}
	if limit, ok := limitParams[s.Limit]; ok {
		sqlBuilder.Limit(limit)
	}
	if offset, ok := limitParams[s.Offset]; ok {
		sqlBuilder.Offset(offset)
	}
	for k, v := range orderBy {
		sqlBuilder.OrderBy(surface+"."+k, v)
	}
	if code, retMaps := GetPurchaseList(sqlBuilder); strings.EqualFold(code, stat.Success) {
		return stat.Success, count, retMaps
	} else {
		return code, 0, make([]map[string]interface{}, 0)
	}
}

func GetPurchaseList(sqlBuilder *svc.SqlBuilder) (string, []map[string]interface{}) {
	query := sqlBuilder.GetCustomerSql(purchaseListSql)
	values := sqlBuilder.GetValues()
	beego.Debug("GetPurchaseList: ", query, values)

	o := orm.NewOrm()
	var resultMaps []orm.Params
	retList := make([]map[string]interface{}, 0)
	_, err := o.Raw(query, values...).Values(&resultMaps)
	if err == nil {
		retList = make([]map[string]interface{}, len(resultMaps))
		for idx, oldMap := range resultMaps {
			retList[idx] = transPurchaseMap(oldMap)
		}
		return stat.Success, retList
	}
	beego.Error("GetPurchaseList Query error:", err)
	return stat.Failed, retList
}

func GetPurchaseTotal(sqlBuilder *svc.SqlBuilder) int64 {
	query := sqlBuilder.GetCustomerSql(purchaseCountSql)
	values := sqlBuilder.GetValues()
	beego.Debug("GetPurchaseTotal: ", query, ":", values)
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(query, values...).Values(&maps); err == nil {
		if len(maps) <= 0 {
			return 0
		}
		if total, ok := maps[0]["count"]; ok {
			total64, err := strconv.ParseInt(total.(string), 10, 64)
			if err != nil {
				beego.Error("GetPurchaseTotal error: ", err)
				return 0
			}
			return total64
		}
	}
	return 0
}

func transPurchaseMap(oldMap orm.Params) t.ItemMap {
	var retMap = make(t.ItemMap, 0)
	for key, value := range oldMap {
		retMap[strings.ToLower(key)] = value
	}
	if userName, ok := oldMap["user_name"]; ok {
		retMap["buyer"] = userName
	}
	if supplierSn, ok := retMap[s.Supplier]; ok && !u.IsNullStr(supplierSn) {
		if supplierMap, sok := supplierMgr.Get(supplierSn.(string)); sok {
			retMap[s.Supplier] = u.GetStringValue(supplierMap, s.Name)
		}
	}
	if productSn, ok := retMap[s.Product]; ok && !u.IsNullStr(productSn) {
		if productMap, sok := productMgr.Get(productSn.(string)); sok {
			retMap[s.Product + s.Name] = u.GetStringValue(productMap, s.Name)
			retMap[s.Product + s.Brand] = u.GetStringValue(productMap, s.Brand)
			retMap[s.Product + s.Model] = u.GetStringValue(productMap, s.Model)
		}
	}
	return retMap
}

