package wborm
import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strconv"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/stat"
)


func RawCount(sql string, values []interface{})int64{
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(sql, values...).Values(&maps); err == nil {
		beego.Debug("GetCount :", maps)
		if len(maps) <= 0 {
			beego.Error("GetCount error: len(maps) < 0")
			return 0
		}
		if total, ok := maps[0]["count"]; ok {
			total64, err := strconv.ParseInt(total.(string), 10, 64)
			if err != nil {
				beego.Error("GetCount error: ", err)
				return 0
			}
			return total64
		}
	}else{
		beego.Error("GetCount error: ", err)
		return 0
	}
	beego.Error("GetCount error:")
	return 0
}

func QueryCount(queryParams t.Params){

}

func QueryListValues(queryParams t.Params, limitParams t.LimitParams, orderBy t.Params, group string) (string, int64, []map[string]interface{}){

	sql := `SELECT supplier, count(id) AS c, ROUND(count(CASE WHEN godowndate != "" AND godowndate <= requireddate THEN "ontime" END)*100/CAST(count(id) AS FLOAT), 2) AS rat FROM purchase WHERE supplier != ""`
	sqlBuilderCount := svc.NewSqlBuilder()
	sqlBuilderCount.Filters(queryParams)
	cqlCount := sqlBuilderCount.GetCustomerSql(sql) + " GROUP BY " + group
	valuesCount := sqlBuilderCount.GetValues()
	beego.Debug("GetPurchaseTotal: ", cqlCount, ":", valuesCount)
	total := RawCount(cqlCount, valuesCount)
	if total == -1 {
		return stat.Failed, 0, make([]map[string]interface{}, 0)
	}

	sqlBuilder := svc.GetSqlBuilder(queryParams, limitParams, orderBy)
	query := sqlBuilder.GetCustomerSql(sql) + " GROUP BY supplier"
	beego.Debug("GetSupplierTimelyList sql is: ", sql)
	var resultMaps []orm.Params
	retList := make([]map[string]interface{}, 0)
	o := orm.NewOrm()
	values := sqlBuilder.GetValues()
	_, err := o.Raw(query, values...).Values(&resultMaps)
	if err == nil {
		retList = make([]map[string]interface{}, len(resultMaps))
//		for idx, oldMap := range resultMaps {
//			retList[idx] = transPurchaseMap(oldMap)
//		}
		return stat.Success, total, retList
	}
	beego.Error("GetPurchaseList Query error:", err)
	return stat.Failed, total, retList
}