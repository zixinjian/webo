package main

import (
	"fmt"
	"webo/models/itemDef"
)

func main() {
	//	maps := itemDef.ReadDefFromCsv()
	fmt.Println("m", itemDef.EntityDefMap)
	//	for _, v := range itemDef.EntityDefMap{
	//		if v.Name == "supplier"{
	//			for _, f := range v.Fields{
	////				name,type,label,require,unique,input,model,default,enum
	//				fmt.Println(fmt.Sprintf("field,%s,%s,%s,%s,%s,%s,%s,%s,%s", f.Name, f.Type, f.Label, f.Require, f.Unique, f.Input, f.Model, "", ""))
	//			}
	//		}
	//	}
}
