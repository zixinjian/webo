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
	"webo/models/util"
)
func GetItems(item string, queryParams Params, orderBy Params)(string, []map[string]interface{}){
	code, retMaps := Query(item, queryParams, LimitParams{}, orderBy)
	return code, retMaps
}
func Query(entity string, queryParams Params, limitParams map[string]int64, orderBy Params) (string, []map[string]interface{}) {
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
//				fmt.Println(value.(string))
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
func List(entity string, queryParams Params, limitParams LimitParams, orderBy Params) (string, int64, []map[string]interface{}) {
	total := Count(entity, queryParams)
	code, retMaps := Query(entity, queryParams, limitParams, orderBy)
	return code, total, retMaps
}
func Count(entity string, params Params) int64 {
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

func Get(entity string, params Params) (string, map[string]interface{}) {
	_, retList := Query(entity, params, map[string]int64{}, Params{})
	if len(retList) > 0 {
		return stat.Success, retList[0]
	}
	return stat.ItemNotFound, nil
}

func Add(entity string, params Params) (string, string) {

	Q := "'"
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		return stat.ItemNotDefine, ""
	}
	nFieldLen := len(oEntityDef.Fields)
	fields := make([]string, nFieldLen)
	marks := make([]string, nFieldLen)
	values := make([]interface{}, nFieldLen)
	for idx, field := range oEntityDef.Fields {
		fields[idx] = field.Name
		marks[idx] = "?"
		value, ok := params[field.Name]
		if ok {
			values[idx] = value
			continue
		}
		if field.Model == s.Sn {
			values[idx] = util.TUId()
			continue
		}
		if field.Model == s.CurTime {
			now := time.Now().Unix()
			values[idx] = now
			continue
		}

		values[idx] = field.Default
	}
	//	fmt.Println("values", values)
	//	fmt.Println(marks)

	sep := fmt.Sprintf("%s, %s", Q, Q)
	qmarks := strings.Join(marks, ", ")
	columns := strings.Join(fields, sep)

	query := fmt.Sprintf("INSERT INTO %s%s%s (%s%s%s) VALUES (%s)", Q, entity, Q, Q, columns, Q, qmarks)
	//
	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		//		b, c := res.LastInsertId()
		//		fmt.Println("e", b, c)
		if i, e := res.LastInsertId(); e == nil && i > 0 {
			return stat.Success, ""
		}else{
			fmt.Println("add,error", e, i)
//			beego.Error(e, i)
		}
	} else {
		beego.Error("Add error", err)
		return ParseSqlError(err, oEntityDef)
	}
	return stat.UnKnownFailed, ""
}

func Update(entity string, params Params) (string, string) {
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

