package svc

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/t"
	"webo/models/u"
	"webo/models/wblog"
)
const DbMark = "'"

func GetItems(item string, queryParams t.Params, orderBy t.Params, limitParams t.LimitParams) (string, []map[string]interface{}) {
	code, retMaps := Query(item, queryParams, limitParams, orderBy)
	return code, retMaps
}
func Query(entity string, queryParams t.Params, limitParams map[string]int64, orderBy t.Params) (string, []map[string]interface{}) {
	sqlBuilder := NewSqlBuilder()
	sqlBuilder.QueryTable(entity)
	for k, v := range queryParams {
		sqlBuilder.Filter(k, v)
	}
	if limit, ok := limitParams[s.Limit]; ok {
		sqlBuilder.Limit(limit)
	}
	if offset, ok := limitParams[s.Offset]; ok {
		sqlBuilder.Offset(offset)
	}
	for k, v := range orderBy {
		sqlBuilder.OrderBy(k, v)
	}
	query := sqlBuilder.GetSql()

	values := sqlBuilder.GetValues()
	//fmt.Println("buildsql: ", query)
	o := orm.NewOrm()
	var resultMaps []orm.Params
	retList := make([]map[string]interface{}, 0)
	_, err := o.Raw(query, values...).Values(&resultMaps)
	if err == nil {
		//		fmt.Println("res", res, resultMaps)
		retList = make([]map[string]interface{}, len(resultMaps))
		//		fmt.Println("old", resultMaps)
		for idx, oldMap := range resultMaps {
			var retMap = make(map[string]interface{}, len(oldMap))
			for key, value := range oldMap {
				retMap[strings.ToLower(key)] = value
			}
			retList[idx] = retMap
		}
		return stat.Success, retList
	} else {
		beego.Error(fmt.Sprintf("Query error:%s for sql:%s", err.Error(), query))
	}
	return stat.Failed, retList
}
func List(entity string, queryParams t.Params, limitParams t.LimitParams, orderBy t.Params) (string, int64, []map[string]interface{}) {
	total := Count(entity, queryParams)
	code, retMaps := Query(entity, queryParams, limitParams, orderBy)
	return code, total, retMaps
}
func Count(entity string, params t.Params) int64 {
	sqlBuilder := NewSqlBuilder()
	sqlBuilder.QueryTable(entity)
	for k, v := range params {
		sqlBuilder.Filter(k, v)
	}
	query := sqlBuilder.GetCountSql()
	values := sqlBuilder.GetValues()
	//	fmt.Println("buildsqlcount: ", query)
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(query, values...).Values(&maps); err == nil {
		//		fmt.Println("res", res, maps)
		if total, ok := maps[0]["COUNT(id)"]; ok {
			total64, err := strconv.ParseInt(total.(string), 10, 64)
			if err != nil {
				panic(err)
			}
			return total64
		}
	}
	return -1
}

func Get(entity string, params t.Params) (string, map[string]interface{}) {
	_, retList := Query(entity, params, map[string]int64{}, t.Params{})
	if len(retList) > 0 {
		return stat.Success, retList[0]
	}
	return stat.ItemNotFound, nil
}

func Delete(item string, sn string, deleter string)string{
	status, oItem := Get(item, t.Params{"sn":sn})
	if creater, ok := oItem[s.Creater]; ok{
		if creater != deleter{
			return stat.PermissionDenied
		}
	}
	if status == stat.Success{
		wblog.LogItemDelete(oItem, deleter)
		return DeleteItem(item, sn)
	}else{
		return status
	}
}

func Add(entity string, params t.Params) (string, string) {

	Q := "'"
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	beego.Error(oEntityDef.Fields)
	if !ok {
		return stat.ItemNotDefine, ""
	}
	fields := make([]string, 0)
	marks := make([]string, 0)
	values := make([]interface{}, 0)
	for _, field := range oEntityDef.Fields {
		if strings.EqualFold(field.Model, s.Upload) {
			continue
		}
		fields = append(fields, field.Name)
		marks = append(marks, "?")
		value, ok := params[field.Name]
		if ok {
			values = append(values, value)
			continue
		}
		if field.Model == s.Sn {
			values = append(values, u.TUId())
			continue
		}
		if field.Model == s.CurTime {
			values = append(values, time.Now().Unix())
			continue
		}
		values = append(values, field.Default)
	}

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(fields, sep)

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, entity, Q, Q, columns, Q, qmarks)
	beego.Debug("Add Item", query, values)
	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		//		b, c := res.LastInsertId()
		//		fmt.Println("e", b, c)
		if i, e := res.LastInsertId(); e == nil && i > 0 {
			return stat.Success, ""
		} else {
			beego.Error(e, i)
			return ParseSqlError(err, oEntityDef)
		}
	} else {
		beego.Error("Add error", err)
		return ParseSqlError(err, oEntityDef)
	}
	return stat.UnKnownFailed, ""
}

func Update(entity string, params t.Params) (string, string) {
	Q := "'"
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		return stat.ItemNotDefine, ""
	}

	id, ok := params[s.Sn]
	if !ok {
		return stat.SnNotFound, ""
	}
	var names []string
	var values []interface{}
	for _, field := range oEntityDef.Fields {
		if field.Name == s.Sn {
			continue
		}
		if value, ok := params[field.Name]; ok {
			values = append(values, value)
			names = append(names, field.Name)
		}
	}
	values = append(values, id)

	sep := fmt.Sprintf("%s = ?, %s", Q, Q)
	setColumns := strings.Join(names, sep)

	query := fmt.Sprintf("UPDATE %s%s%s SET %s%s%s = ? WHERE %s = ?", Q, entity, Q, Q, setColumns, Q, s.Sn)
	//	fmt.Println("sql", query, values)
	beego.Debug("Update sql: %s", query)
	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		if i, e := res.RowsAffected(); e == nil && i > 0 {
			return stat.Success, ""
		}
	} else {
		beego.Error("Update error", err)
		return ParseSqlError(err, oEntityDef)
	}
	return stat.UnKnownFailed, ""
}

func DeleteItem(item string, sn string) string {
	if _, ok := itemDef.EntityDefMap[item];!ok {
		return stat.ItemNotDefine
	}

	Q := DbMark
	query := fmt.Sprintf("DELETE FROM %s%s%s WHERE sn = ?", Q, item, Q)

	values := make([]interface{}, 1)
	values[0] = sn

	beego.Debug("Delete sql: %s : sn: %s", query, sn)
	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		if i, e := res.RowsAffected(); e == nil && i > 0 {
			return stat.Success
		}else {
			beego.Error("Delete failed", err)
			return stat.UnKnownFailed
		}
	} else {
		beego.Error("Delete error", err)
		return stat.UnKnownFailed
	}
	return stat.UnKnownFailed
}

func ParseSqlError(err error, oEntityDef itemDef.ItemDef) (string, string) {
	errStr := err.Error()
	if strings.HasPrefix(errStr, SqlErrUniqueConstraint) {
		itemAndField := strings.TrimPrefix(errStr, SqlErrUniqueConstraint)
		lstStr := strings.Split(itemAndField, ".")
		if len(lstStr) < 2 {
			return stat.DuplicatedValue, itemAndField
		}
		field := strings.TrimSpace(lstStr[1])
		if v, ok := oEntityDef.GetField(field); ok {
			return stat.DuplicatedValue, v.Label
		}
		return stat.DuplicatedValue, itemAndField
	}
	beego.Error("ParseSqlError unknown error", errStr)
	return stat.UnKnownFailed, ""
}
