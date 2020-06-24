package ctrl

import (
	"strings"
	"wb/cs"
	"wb/ii"
	"wb/lg"
	"wb/st"
)

type ItemDynamicController struct {
	ItemController
	isCtxItemInfoOk bool
}

func (this *ItemDynamicController) NestPrepare() {
	this.ItemController.NestPrepare()
	lg.Debug("ItemController.ItemDynamicController", " version: ", this.CtxVersion)
	this.isCtxItemInfoOk = false
	if svc, ok := UiInjectionMap[this.Ctx.Input.URL()]; ok {
		itemName := svc.Item.Name
		if itemName != "" {
			oItemInfo, ok := ii.ItemInfoMap[itemName]
			if ok {
				this.CtxItemInfo = oItemInfo
				this.isCtxItemInfoOk = true
			}
		}
	}
	if this.isCtxItemInfoOk == false {
		oItemInfo, code := this.GetItemInfoFromParamHi()
		if code == st.Success {
			this.CtxItemInfo = oItemInfo
			this.isCtxItemInfoOk = true
		} else {
			this.SendJson(&cs.JsonResult{st.ItemNotDefine, st.ParamItemIsNone_})
			this.StopRun()
		}
	}
}
func (this *ItemDynamicController) GetItemNameFromParamHi() (string, string) {
	itemName, ok := this.Ctx.Input.Params()[":hi"]
	if !ok {
		lg.Error(st.ParamItemIsNone_, this.Ctx.Input.Params())
		return "", st.ParamItemError
	}
	return itemName, st.Success
}

func (this *BaseController) GetItemInfoFromParamHi() (ii.ItemInfo, string) {
	itemName, ok := this.Ctx.Input.Params()[":hi"]
	if !ok {
		lg.Error(st.ParamItemIsNone_, this.Ctx.Input.Params())
		return ii.ItemInfo{}, st.ParamItemError
	}
	return ii.GetItemInfoByName(itemName)
}
func (this *ItemDynamicController) UiList() {
	// 如果可以查到ItemInfo则初始化为默认值
	if this.isCtxItemInfoOk {
		this.UiListItem(this.CtxItemInfo)
	} else {
		this.Data["item"] = "other"
		this.Data["thlist"] = ""
		this.TplName = "item/list.tpl"
	}
	this.injectUiList()
}
func (this *ItemDynamicController) injectUiList() {
	if injetion, ok := UiInjectionMap[this.Ctx.Input.URL()]; ok {
		for _, v := range injetion.Data {
			this.Data[v.Key] = v.Data
		}
		if strings.TrimSpace(injetion.Tpl) != "" {
			this.TplName = injetion.Tpl
		}
	}
}
func (this *ItemDynamicController) UiAdd() {
	if this.isCtxItemInfoOk {
		this.UiAddItem(this.CtxItemInfo)
	} else {
		this.Data["Form"] = ""
		this.Data["UrlService"] = ""
		this.Data["Onload"] = ""
		this.TplName = "item/edit.tpl"
	}
	this.injectUiList()
}
func (this *ItemDynamicController) UiUpdate() {
	if this.isCtxItemInfoOk {
		this.UiUpdateItem(this.CtxItemInfo)
	} else {
		this.Data["Form"] = ""
		this.Data["UrlService"] = ""
		this.Data["Onload"] = ""
		this.TplName = "item/edit.tpl"
	}
	this.injectUiList()
}
