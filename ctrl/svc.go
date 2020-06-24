package ctrl

import (
	"fmt"
	"strings"
	"wb/cs"
	"wb/ii"
	"wb/lg"
	"wb/om"
	"wb/st"
	"wb/svc"
)

type SvcController struct {
	ItemBaseController
}

func (this *SvcController) UiSvc() {
	svc, ok := UiInjectionMap[this.Ctx.Input.URL()]
	if !ok {
		this.Ctx.WriteString(fmt.Sprintln("No such config:", this.Ctx.Input.URL()))
		return
	}
	if svc.Item.Name != "" {
		oItemInfo, ok := ii.ItemInfoMap[svc.Item.Name]
		if !ok {
			lg.Error(st.ItemNotDefine_, svc.Item.Name)
			this.Ctx.WriteString(st.ItemNotDefine_)
			return
		}
		this.Data["item"] = oItemInfo.Name
		// for _, injection := range svc.Item.Injections{
		//    fmt.Println("Injection", injection)
		//    switch injection {
		//    case "From":
		//        this.FillFormElement(uibuilder.BuildFormElement(oItemInfo, oldValueMap, statusMap))
		//    }
		// }
	}
	for _, param := range svc.Params {
		key := param.Name
		this.Data["Param_"+key] = this.GetString(key)
		// fmt.Println("param", param, key, this.GetString(key))
	}
	for _, v := range svc.Data {
		this.Data[v.Key] = v.Data
	}
	this.TplName = svc.Tpl
}
func (this *SvcController) Sqs() (string, cs.MObject) {
	sqlKey := this.Ctx.Input.URL()
	sqlKey = strings.TrimPrefix(sqlKey, "/sqs/")
	params := om.Params{}
	for k, v := range this.Ctx.Request.PostForm {
		params[k] = v
	}
	status, sqList := svc.SqsExec(sqlKey, params)
	return status, sqList
}
