package wbo

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"webo/models/itemDef"
	"webo/models/s"
	"webo/models/t"
	"webo/models/u"
)

type SqlBuilder struct {
	oEntityDef itemDef.ItemDef
	table      []string
	limit      int64
	offset     int64
	orders     []string
	conditions []condition
	groupBy	   string
	//	relations       []relation
}

//type relation struct {
//	Table 			string
//	Field			string
//	RTable 			string
//	RField 			string
//	Opt				string
//}

type condition struct {
	FieldName string
	Value     interface{}
	Opt       string
}

func (this *SqlBuilder) QueryTable(table string) {
	oEntityDef, ok := itemDef.EntityDefMap[table]
	if !ok {
		beego.Error(fmt.Errorf("<queryBuilder.QueryTable>no such entry table define: %v", table))
	}
	this.table = append(this.table, table)
	this.oEntityDef = oEntityDef
}

func (this *SqlBuilder) QueryTables(tables ...string) {
	for _, table := range tables {
		this.QueryTable(table)
	}
}

func (this *SqlBuilder) Filter(key string, value interface{}) {
	switch key[:1] {
	case "-":
		if v, ok := value.(string); ok {
			this.addCondition(key[1:], v, s.Like)
		} else {
			beego.Error("Add filter error startswith not string")
		}
	case "%":
		if v, ok := value.(string); ok {
			this.addCondition(key[1:], v+"%", s.Like)
		} else {
			beego.Error("Add filter error startswith not string")
		}
	default:
		this.addCondition(key, value, "=")
	}
}

func (this *SqlBuilder) Filters(queryParam t.Params){
	for k, v := range queryParam{
		this.Filter(k, v)
	}
}

func (this *SqlBuilder) GroupBy(groupBy string){
	this.groupBy = groupBy
}

func (this *SqlBuilder) addCondition(fieldName string, value interface{}, opt string) {
	//	beego.Debug("AddCondition:", fieldName, value, opt)
	this.conditions = append(this.conditions, condition{fieldName, value, opt})
	//	}
}
func (this *SqlBuilder) Limit(limit int64) {
	if limit > 0 {
		this.limit = limit
	}
}
func (this *SqlBuilder) Offset(offset int64) {
	if offset > 0 {
		this.offset = offset
	}
}
func (this *SqlBuilder) Limits(limitParams t.LimitParams) {
	if limit, ok := limitParams[s.Limit]; ok {
		this.Limit(limit)
	}
	if offset, ok := limitParams[s.Offset]; ok {
		this.Offset(offset)
	}
}
func (this *SqlBuilder) OrderBy(fieldName string, value interface{}) {
	if strings.EqualFold(value.(string), "ASC") {
		this.orders = append(this.orders, fieldName+" ASC")
		return
	}
	if strings.EqualFold(value.(string), "DESC") {
		this.orders = append(this.orders, fieldName+" DESC")
	}
}
func (this *SqlBuilder) OrderBys(orderBy t.Params){
	for k, v := range orderBy {
		this.OrderBy(k, v)
	}
}
func (this *SqlBuilder) GetConditonSql() string {
	sql := ""
	for idx, cond := range this.conditions {
		if idx > 0 {
			sql = sql + "AND "
		}
		sql = sql + cond.FieldName + fmt.Sprintf(" %s ? ", cond.Opt)
	}
	return sql
}

func (this *SqlBuilder) GetCountSql() string {
	sql := fmt.Sprintf("SELECT COUNT(id) FROM %s ", this.table)
	if len(this.conditions) > 0 {
		sql = sql + "WHERE " + this.GetConditonSql()
	}
	if !u.IsNullStr(this.groupBy){
		sql = sql + " GROUP BY " + this.groupBy
	}
	return sql
}
func (this *SqlBuilder) GetSql() string {
	tableStr := strings.Join(this.table, ",")
	sql := fmt.Sprintf("SELECT * FROM %s ", tableStr)

	return this.GetCustomerSql(sql)
}

func (this *SqlBuilder) GetCustomerSql(sql string) string {
	sql = sql + " "
	if len(this.conditions) > 0 {
		if strings.Contains(strings.ToUpper(sql), "WHERE") {
			sql = sql + "AND "
		} else {
			sql = sql + "WHERE "
		}
		sql = sql + this.GetConditonSql()
	}
	if !u.IsNullStr(this.groupBy){
		sql = sql + " GROUP BY " + this.groupBy + " "
	}
	if len(this.orders) > 0 {
		sql = sql + "ORDER BY "
		for idx, v := range this.orders {
			if idx > 0 {
				sql = sql + ", "
			}
			sql = sql + v + " "
		}
	}
	if this.limit > 0 {
		sql = sql + fmt.Sprintf("LIMIT %d ", this.limit)
	}
	if this.offset > 0 {
		sql = sql + fmt.Sprintf("OFFSET %d ", this.offset)
	}
	return sql
}

func (this *SqlBuilder) GetFrom() string {
	tableStr := strings.Join(this.table, ",")
	return fmt.Sprintf("FROM %s ", tableStr)
}

func (this *SqlBuilder) GetValues() []interface{} {
	values := make([]interface{}, len(this.conditions))
	for idx, con := range this.conditions {
		values[idx] = con.Value
	}
	return values
}

func NewSqlBuilder() *SqlBuilder {
	o := &SqlBuilder{}
	o.limit = 0
	o.offset = 0
	return o
}

func CreateSqlBuilder(queryParams t.Params, limitParams t.LimitParams, orderBy t.Params, groupBy string) *SqlBuilder{
	o := &SqlBuilder{}
	o.Filters(queryParams)
	o.GroupBy(groupBy)
	o.Limits(limitParams)
	return o
}