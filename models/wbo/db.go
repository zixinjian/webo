package wbo
import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strconv"
	"webo/models/t"
	"webo/models/stat"
	"fmt"
	"strings"
)

const SqlErrUniqueConstraint = "UNIQUE constraint failed: "

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
	sqlBuilder := NewSqlBuilder()
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
	sqlBuilder := CreateSqlBuilder(queryParams, limitParams, orderBys, groupBy)
	query := sqlBuilder.GetCustomerSql(sql)
	values := sqlBuilder.GetValues()
	return DbQueryValues(query, values)
}
//func QueryListValues(sql string, queryParams t.Params, limitParams t.LimitParams, orderBys t.Params, groupBy string) (string, int64, []map[string]interface{}){
//	_, total := QueryCount(sql, queryParams, groupBy)
//	if total == -1 {
//		return stat.Failed, 0, make([]map[string]interface{}, 0)
//	}
//	status, retMaps:= QueryValues(sql, queryParams, limitParams, orderBys, groupBy)
//	return status, total, retMaps
//}

func DbQueryValues(query string, values ...interface{})(string, []map[string]interface{}){
	beego.Debug(fmt.Sprintf("RawQueryMaps sql: %s, values: %s", query, values))
	var resultMaps []orm.Params
	o := orm.NewOrm()
	_, err := o.Raw(query, values).Values(&resultMaps)
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

func DbDelete(query string, values ...interface{}) string{
	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		if i, e := res.RowsAffected(); e == nil && i > 0 {
			return stat.Success
		}else {
			beego.Error("RawDelete failed", err)
			return stat.UnKnownFailed
		}
	} else {
		beego.Error("RawDelete error", err)
		return stat.UnKnownFailed
	}
	return stat.UnKnownFailed
}
func DbUpdate(sql string, values ...interface{})string{
	o := orm.NewOrm()
	if res, err := o.Raw(sql, values...).Exec(); err == nil {
		if i, e := res.RowsAffected(); e == nil && i > 0 {
			return stat.Success
		}
	} else {
		beego.Error("Update error", err)
		return stat.Failed
	}
	return stat.UnKnownFailed
}
//func AddValue(itemName string, params orm.Params)(string, string){
//	sn := u.TUId()
//	fields := make([]string, 0)
//	marks := make([]string, 0)
//	values := make([]interface{}, 0)
//	for _, field := range oEntityDef.Fields {
//		if strings.EqualFold(field.Model, s.Upload) {
//			continue
//		}
//		fields = append(fields, field.Name)
//		marks = append(marks, "?")
//		value, ok := params[field.Name]
//		if ok {
//			values = append(values, value)
//			continue
//		}
//		if field.Model == s.Sn {
//			values = append(values, sn)
//			continue
//		}
//		if field.Model == s.CurTime {
//			values = append(values, time.Now().Unix())
//			continue
//		}
//		values = append(values, field.Default)
//	}
//
//	sep := fmt.Sprintf("%s, %s", Q, Q)
//	qmarks := strings.Join(marks, ", ")
//	columns := strings.Join(fields, sep)
//
//	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, entity, Q, Q, columns, Q, qmarks)
//	beego.Debug("Add Item", query, values)
//	o := orm.NewOrm()
//	if res, err := o.Raw(query, values...).Exec(); err == nil {
//		//		b, c := res.LastInsertId()
//		//		fmt.Println("e", b, c)
//		if i, e := res.LastInsertId(); e == nil && i > 0 {
//			return stat.Success, sn
//		} else {
//			beego.Error(e, i)
//			return ParseSqlError(err, oEntityDef)
//		}
//	} else {
//		beego.Error("Add error", err)
//		return ParseSqlError(err, oEntityDef)
//	}
//	return stat.UnKnownFailed, ""
//}
//
//func parseSqlError(err error){
//	errStr := err.Error()
//	if strings.HasPrefix(errStr, SqlErrUniqueConstraint) {
//		itemAndField := strings.TrimPrefix(errStr, SqlErrUniqueConstraint)
//		lstStr := strings.Split(itemAndField, ".")
//		if len(lstStr) < 2 {
//			return stat.DuplicatedValue, itemAndField
//		}
//		field := strings.TrimSpace(lstStr[1])
//		if v, ok := oEntityDef.GetField(field); ok {
//			return stat.DuplicatedValue, v.Label
//		}
//		return stat.DuplicatedValue, itemAndField
//	}
//	beego.Error("ParseSqlError unknown error", errStr)
//	return stat.UnKnownFailed, ""
//}