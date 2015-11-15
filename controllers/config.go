package controllers
import (
	"webo/models/wbconf"
)


type ConfigController struct {
	ItemController
}

func (this *ConfigController) Add() {
	this.ItemController.Add()
	wbconf.LoadCategory()
}

func (this *ConfigController) Update() {
	this.ItemController.Update()
	wbconf.LoadCategory()
}

