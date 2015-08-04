package svc

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
	"webo/models/itemDef"
	"webo/models/util"
)

func Query(entity string, queryParams Params, limitParams map[string]int64, orderBy Params) (string, []map[string]interface{}) {
	sqlBuilder := NewQueryBUilder()
	sqlBuilder.QueryTable(entity)
	for k, v := range queryParams {
		sqlBuilder.Filter(k, v)
	}
	//	fmt.Println("order", orderBy)
	if limit, ok := limitParams["limit"]; ok {
		sqlBuilder.Limit(limit)
	}
	if offset, ok := limitParams["offset"]; ok {
		sqlBuilder.Offset(offset)
	}
	for k, v := range orderBy {
		sqlBuilder.OrderBy(k, v)
	}
	query := sqlBuilder.GetSql()
	values := sqlBuilder.GetValues()
	//	fmt.Println("buildsql: ", query)
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
		return "success", retList
	} else {
		fmt.Println("res", err)
	}
	return "faild", retList
}

func List(entity string, queryParams Params, limitParams map[string]int64, orderBy Params) (string, int64, []map[string]interface{}) {
	total := Count(entity, queryParams)
	code, retMaps := Query(entity, queryParams, limitParams, orderBy)
	return code, total, retMaps
}
func Count(entity string, params Params) int64 {
	sqlBuilder := NewQueryBUilder()
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

			fmt.Println("total", total64)
			return total64
		}
	}
	return -1
}

func Get(entity string, params Params) (string, map[string]interface{}) {
	_, _, retList := List(entity, params, map[string]int64{}, Params{})
	if len(retList) > 0 {
		return "success", retList[0]
	}
	return "not_found", nil
}

func Add(entity string, params Params) string {

	Q := "'"
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		return "entity_no_define"
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
		if field.Model == "sn" {
			values[idx] = util.TUId()
			continue
		}
		if field.Model == "curtime" {
			now := time.Now().Unix()
			values[idx] = now
			//			fmt.Println("time", time.Unix(now , 0).String())
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
	if res, err := o.Raw(query, values...).Exec(); err != nil {
		fmt.Println(err)
		fmt.Println("res", res)
	}
	return "success"
}

func Update(entity string, params Params) string {
	Q := "'"
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		return "entity_no_define"
	}

	id, ok := params["id"]
	if !ok {
		return "no_id"
	}
	var names []string
	var values []interface{}
	for _, field := range oEntityDef.Fields {
		if field.Name == "id" {
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

	query := fmt.Sprintf("UPDATE %s%s%s SET %s%s%s = ? WHERE %s%s%s = ?", Q, entity, Q, Q, setColumns, Q, Q, "id", Q)

	o := orm.NewOrm()
	if res, err := o.Raw(query, values...).Exec(); err == nil {
		fmt.Println("res", res)
	}
	return "success"
}
