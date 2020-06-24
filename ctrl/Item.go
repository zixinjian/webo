package ctrl

import (
	"wb/cs"
	"wb/ii"
	"wb/om"
)

type ItemController struct {
	ItemBaseController
	CtxItemInfo ii.ItemInfo
}

func (this *ItemController) List() {
	tableResult := this.GetTableList(this.CtxItemInfo, om.Params{})
	this.SendJson(&tableResult)
}
func (this *ItemController) Get() {
	ret, res := this.GetItem(this.CtxItemInfo)
	this.SendJson(cs.JsonResult{ret, res})
}
func (this *ItemController) Add() {
	this.AddItem(this.CtxItemInfo)
}
func (this *ItemController) Adds() {
	this.AddItems(this.CtxItemInfo)
}
func (this *ItemController) AddWithSub() {
	this.AddItemWithSub(this.CtxItemInfo)
}
func (this *ItemController) Update() {
	this.UpdateItem(this.CtxItemInfo)
}
func (ic *ItemController) UpdateWithSub() {
	ic.UpdateItemWithSub(ic.CtxItemInfo)
}
func (this *ItemController) UpdateWithAddSub() {
	this.UpdateItemWithAddSub(this.CtxItemInfo)
}

//func (this *ItemController) BatchUpdate(){
//    this.UpdateBatchItem(this.CtxItemInfo)
//}
func (this *ItemController) Delete() {
	this.DeleteItem(this.CtxItemInfo)
}
func (this *ItemController) Upload() {
	this.UploadItem(this.CtxItemInfo)
}
func (this *ItemController) Autocomplete() {
	this.AutocompleteItem(this.CtxItemInfo)
}
func (this *ItemController) UiList() {
	this.UiListItem(this.CtxItemInfo)
}
func (this *ItemController) UiAdd() {
	this.UiAddItem(this.CtxItemInfo)
}
func (this *ItemController) UiUpdate() {
	this.UiUpdateItem(this.CtxItemInfo)
}
func (this *ItemController) UiView() {
	this.UiViewItem(this.CtxItemInfo)
}
func (this *ItemController) UiAttachment() {
	this.UiAttachmentItem(this.CtxItemInfo)
}
