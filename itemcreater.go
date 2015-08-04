package main

import (
	"fmt"
	"webo/models/itemDef"
)

func main() {
	createTh()
    createsql()
}
func createTh() {
	for itemName, oItemDef := range itemDef.EntityDefMap {
		fmt.Println("itemName:", itemName)
		for _, field := range oItemDef.Fields {
			if field.Input != "none" {
				th := fmt.Sprintf(`<th data-field="%s" data-sortable="true">%s</th>`, field.Name, field.Label)
				fmt.Println(th)
			}
		}
	}
}
func createsql() {
	for itemName, oItemDef := range itemDef.EntityDefMap {
		//        fmt.Println(itemName, oItemDef)
		fieldsql := "CREATE TABLE " + itemName + " (id integer NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE"
		for _, field := range oItemDef.Fields {
			//            fmt.Println("idx:", idx)
			fieldsql = fieldsql + "," + field.Name
			switch field.Model {
			case "sn", "text", "password", "enum", "curuser":
				fieldsql = fieldsql + " varchar"
			case "timestamp":
				fieldsql = fieldsql + " time"
			case "integer", "curtime":
				fieldsql = fieldsql + " integer"
			default:
				fmt.Println("no such modal", field.Name, field.Model)
			}
			if field.Require == "true" {
				fieldsql = fieldsql + " NOT NULL"
			}
			if field.Unique == "true" {
				fieldsql = fieldsql + " UNIQUE"
			}
			if field.Default == "" {
				continue
			}
			switch field.Model {
			case "sn", "text", "password", "enum", "integer":
				fieldsql = fieldsql + " DEFAULT " + field.Default.(string)
			default:
				fmt.Println("no default", field.Name)
			}
		}
		fieldsql = fieldsql + ")"
		fmt.Println(fieldsql)
	}
}
