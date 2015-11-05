package wborm
import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strconv"
	"webo/models/svc"
	"webo/models/t"
	"webo/models/stat"
	"fmt"
	"strings"
)


func QueryRawCount(sql string, values []interface{})int64{
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(sql, values...).Values(&maps); err == nil {
		if len(maps) <= 0 {
			beego.Error("GetCount error: len(maps) < 0")
			return -1
		}
		if total, ok := maps[0]["count"]; ok {
			if total64, err := strconv.ParseInt(total.(string), 10, 64); err != nil {
				beego.Error("GetCount Parseint error: ", err)
				return -1
			}else{
				return total64
			}
		}
	}else{
		beego.Error("GetCount error: ", err)
		return -1
	}

	beego.Error("GetCount unknown error: ", maps)
	return 0
}


func QueryCount(sql string, queryParams t.Params, groupBy string)(string, int64){
	sqlBuilder := svc.NewSqlBuilder()
	sqlBuilder.Filters(queryParams)
	sqlBuilder.GroupBy(groupBy)
	query := sqlBuilder.GetCustomerSql(sql)
	values := sqlBuilder.GetValues()
	beego.Debug("QueryCount: ", query, ":", values)
	count := QueryRawCount(query, values)
	if count == -1 {
		return stat.Failed, -1
	}
	return stat.Failed, count
}
func QueryValues(sql string, queryParams t.Params, limitParams t.LimitParams, orderBys t.Params, groupBy string)(string, []map[string]interface{}){
	sqlBuilder := svc.CreateSqlBuilder(queryParams, limitParams, orderBys, groupBy)
	query := sqlBuilder.GetCustomerSql(sql)
	values := sqlBuilder.GetValues()
	beego.Debug(fmt.Sprintf("QueryItems sql: %s, values: %s", query, values))
	var resultMaps []orm.Params
	o := orm.NewOrm()
	_, err := o.Raw(query, values...).Values(&resultMaps)
	if err == nil {
		retList := make([]map[string]interface{}, len(resultMaps))
		for idx, maps := range resultMaps{
			for idx, oldMap := range resultMaps {
				var retMap = make(map[string]interface{}, len(oldMap))
				for key, value := range oldMap {
					retMap[strings.ToLower(key)] = value
				}
				retList[idx] = retMap
			}
			retList[idx] = maps
		}
		return stat.Success, retList
	}
	beego.Error("QueryItems Query error:", err)
	return stat.Failed, make([]map[string]interface{}, 0)
}
//func QueryListValues(sql string, queryParams t.Params, limitParams t.LimitParams, orderBys t.Params, groupBy string) (string, int64, []map[string]interface{}){
//	_, total := QueryCount(sql, queryParams, groupBy)
//	if total == -1 {
//		return stat.Failed, 0, make([]map[string]interface{}, 0)
//	}
//	status, retMaps:= QueryValues(sql, queryParams, limitParams, orderBys, groupBy)
//	return status, total, retMaps
//}

