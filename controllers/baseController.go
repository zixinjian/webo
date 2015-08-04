package controllers

import (
	"github.com/astaxie/beego"
	"webo/models/itemDef"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) GetFormValues(itemD itemDef.ItemDef) map[string]interface{} {
	//fmt.Println("def", itemD)
	var retMap map[string]interface{}
	retMap = make(map[string]interface{})
	formValues := this.Input()
	for _, field := range itemD.Fields {
		if _, ok := formValues[field.Name]; ok {
			if v, fok := field.GetValue(this.GetString(field.Name)); fok {
				retMap[field.Name] = v
			}
		}
	}
	return retMap
}
