package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type AAA struct {
	a string
}

func (aaa *AAA) fuA() {
	fmt.Println("aaa")
}

type BBB struct {
	aaa AAA
	b   string
}

func (bbb *BBB) funB() {
	fmt.Println("bbb")
}

func main() {
	b := BBB{}
	b.funB()
	//	orm.RegisterDriver("sqlite", orm.DR_Sqlite)
	//	orm.RegisterDataBase("default", "sqlite3", "test4.sq")
	//	o := orm.NewOrm()
	//	o.Using("default")
	//	res, err := o.Raw("CREATE TABLE test(id INTEGER, ts INTEGER)").Exec()
	//	if err != nil{
	//		fmt.Println("create", err)
	//	}
	//	fmt.Println("res", res)
	//	res1, err1 := o.Raw("INSERT INTO test (ts) VALUES(?)", "1092941466").Exec()
	//	if err1 != nil{
	//		fmt.Println("update", err1)
	//	}
	//	fmt.Println("res1", res1)
	//	var maps []orm.Params
	//	num, err := o.Raw("SELECT date('now','start of year','+9 months','weekday 2')").Values(&maps)
	//	fmt.Println(num, err)
	//	if err == nil && num > 0 {
	//		fmt.Println(maps[0]["ts"]) // slene
	//	}

}
