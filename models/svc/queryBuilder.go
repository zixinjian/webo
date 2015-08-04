package svc

import (
	"fmt"
	"strings"
	"webo/models/itemDef"
)

type queryBuilder struct {
	oEntityDef itemDef.ItemDef
	table      string
	limit      int64
	offset     int64
	orders     []string
	conditions []condition
}

type condition struct {
	FieldName string
	Value     interface{}
	Opt       string
}

func (this *queryBuilder) QueryTable(table string) {
	oEntityDef, ok := itemDef.EntityDefMap[table]
	if !ok {
		panic(fmt.Errorf("<queryBuilder.QueryTable>no such entry table define: %v", table))
	}
	this.table = table
	this.oEntityDef = oEntityDef
}

func (this *queryBuilder) Filter(fieldName string, value interface{}) {
	if this.oEntityDef.IsValidField(fieldName) {
		this.conditions = append(this.conditions, condition{fieldName, value, "="})
	}
}

func (this *queryBuilder) Limit(limit int64) {
	if limit > 0 {
		this.limit = limit
	}
}

func (this *queryBuilder) Offset(offset int64) {
	if offset > 0 {
		this.offset = offset
	}
}

func (this *queryBuilder) OrderBy(fieldName string, value interface{}) {
	if !this.oEntityDef.IsValidField(fieldName) {
		return
	}
	if strings.EqualFold(value.(string), "DESC") {
		this.orders = append(this.orders, fieldName+" ASC")
	}
	if strings.EqualFold(value.(string), "DESC") {
		this.orders = append(this.orders, fieldName+" DESC")
	}
}
func (this *queryBuilder) GetWhere() string {
	sql := "WHERE "
	for idx, cond := range this.conditions {
		if idx > 0 {
			sql = sql + "AND "
		}
		sql = sql + cond.FieldName + " = ? "
	}
	return sql
}

func (this *queryBuilder) GetCountSql() string {
	sql := fmt.Sprintf("SELECT COUNT(id) FROM %s ", this.table)
	if len(this.conditions) > 0 {
		sql = sql + this.GetWhere()
	}
	return sql
}
func (this *queryBuilder) GetSql() string {
	sql := fmt.Sprintf("SELECT * FROM %s ", this.table)
	if len(this.conditions) > 0 {
		sql = sql + this.GetWhere()
	}
	if this.limit > 0 {
		sql = sql + fmt.Sprintf("LIMIT %d ", this.limit)
	}
	if this.offset > 0 {
		sql = sql + fmt.Sprintf("OFFSET %d ", this.offset)
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
	return sql
}

func (this *queryBuilder) GetValues() []interface{} {
	values := make([]interface{}, len(this.conditions))
	for idx, con := range this.conditions {
		values[idx] = con.Value
	}
	return values
}

func NewQueryBUilder() *queryBuilder {
	o := &queryBuilder{}
	o.limit = 0
	o.offset = 0
	return o
}
