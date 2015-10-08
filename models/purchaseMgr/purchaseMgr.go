package purchaseMgr

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
	"webo/models/productMgr"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/supplierMgr"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/u"
	"webo/models/wborm"
	"fmt"
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

func GetSupplierTimelyList(queryParams t.Params, limitParams t.LimitParams, orderBys t.Params)(string, int64, []map[string]interface{}){
	countSql := `SELECT count(s.sn) as count FROM
			(SELECT supplier, count(id) FROM purchase GROUP BY supplier) as p, supplier as s
			WHERE p.supplier = s.sn `
	_, total := wborm.QueryCount(countSql, queryParams, "")
	if total == -1 {
		return stat.Failed, 0, make([]map[string]interface{}, 0)
	}
	sql := `SELECT s.name as supplier, p.total as total, p.intime as intime, ROUND(p.intime*100/CAST(p.total AS FLOAT), 2) AS rat FROM
			(SELECT supplier, count(id) AS total, count(CASE WHEN godowndate != "" AND godowndate <= requireddate THEN "intime" END) AS intime FROM purchase GROUP BY supplier) as p, supplier as s
			WHERE p.supplier = s.sn `
	status, retMaps := wborm.QueryValues(sql, queryParams, limitParams, orderBys, "")
	beego.Debug("GetSupplierTimelyList : ", retMaps)
	return status, total, retMaps
}

func GetBuyerTimelyList(queryParams t.Params, limitParams t.LimitParams, orderBys t.Params) (string, int64, []map[string]interface{}){
	countSql := `SELECT count(s.sn) as count FROM
			(SELECT buyer, count(id) FROM purchase GROUP BY buyer) as p, user as s
			WHERE p.buyer = s.sn AND s.department = "department_purchase" `
	_, total := wborm.QueryCount(countSql, queryParams, "")
	if total == -1 {
		return stat.Failed, 0, make([]map[string]interface{}, 0)
	}
	sql := `SELECT s.name as buyer, p.total as total, p.intime as intime, ROUND(p.intime*100/CAST(p.total AS FLOAT), 2) AS rat FROM
			(SELECT buyer, count(id) AS total, count(CASE WHEN godowndate != "" AND godowndate <= requireddate THEN "intime" END) AS intime FROM purchase GROUP BY buyer) as p, user as s
			WHERE p.buyer = s.sn AND s.department = "department_purchase" `
	status, retMaps := wborm.QueryValues(sql, queryParams, limitParams, orderBys, "")
//	beego.Debug("GetBuyerTimelyList : ", retMaps)
	return status, total, retMaps
//	status, num, userMaps := svc.List(s.User, queryParams, limitParams, orderBy)
//	if status != stat.Success{
//		return status, num, userMaps
//	}
//	retMaps := make([]map[string]interface{}, len(userMaps))
//	for idx, userMap := range userMaps{
//		sn := u.GetStringValue(userMap, s.Sn)
//		name := u.GetStringValue(userMap, s.Name)
//		retMap := make(map[string]interface{}, 5)
//		noDelay, total, rat := getBuyerTimely(sn)
//		retMap["delay"] = total - noDelay
//		retMap["total"] = total
//		retMap["rat"] = rat
//		retMap[s.Sn] = sn
//		retMap[s.Name] = name
//		retMaps[idx] = retMap
//	}
//	return status, num, retMaps
}
func getBuyerTimely(sn string)(noDelay int64, total int64, rat string){
	queryParams := t.Params{
		s.Buyer:sn,
	}
	noDelay, total, rat = calcTimely(queryParams)
	return
}
func calcTimely(queryParams t.Params)(noDelay int64, total int64, rat string){
	sqlBuilder := svc.NewSqlBuilder()
	for k, v := range queryParams {
		sqlBuilder.Filter(k, v)
	}
	where := sqlBuilder.GetConditonSql()
	noDelaySql := "SELECT count(id) as count FROM purchase WHERE godowndate != '' AND godowndate < requireddate " + "AND " + where
	noDelay = wborm.QueryRawCount(noDelaySql, sqlBuilder.GetValues())
	totalSql := "SELECT count(id) as count FROM purchase WHERE " + where + " AND (godowndate != '' OR requireddate < '?')"
	total = wborm.QueryRawCount(totalSql, append(sqlBuilder.GetValues(), u.GetToday()))
	if total == 0 || noDelay == 0{
		rat = "0%"
		return
	}
	rat = fmt.Sprintf("%.2f", float64(noDelay) * 100/float64(total)) + "%"
	return
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
	return wborm.QueryRawCount(query, values)
}

func CalcProductTimely(queryParams t.Params)map[string]interface{}{
	retMap := make(map[string]interface{}, 3)
	noDelay, total, rat := calcTimely(queryParams)
	retMap["delay"] = total - noDelay
	retMap["total"] = total
	retMap["rat"] = rat
	return retMap
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
			retMap[s.Supplier + s.Name] = u.GetStringValue(supplierMap, s.Name)
			retMap[s.Supplier + s.Key] = u.GetStringValue(supplierMap, s.Keyword)
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

