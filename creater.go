package main

import (
	"fmt"
	"time"
	//	"webo/models/itemDef"
	//	"github.com/astaxie/beego/orm"
	//	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println(fmt.Sprintf("%v %v", "abc", 1), time.Now().Unix())
	//	orm.RegisterDriver("sqlite", orm.DR_Sqlite)
	//	orm.RegisterDataBase("default", "sqlite3", "db/frame.sqlite3")
	//	o := orm.NewOrm()
	//	qs := o.QueryTable("abc")
	//	qs.Limit(10, 0)
	//	var resultMaps []orm.Params
	//	qs.Values(&resultMaps)
	//	fmt.Println("res", len(resultMaps), resultMaps)
	//	var list []string
	//	fmt.Println(len(list))
	//	list = append(list, "abc")
	//	fmt.Println(len(list))
	//	fmt.Println(list)
	//	itemDefMap := itemDef.EntityDefMap
	//	for item, def := range itemDefMap {
	//		fmt.Println(item, def)
	//	}
	//	fmt.Println("start")
}
