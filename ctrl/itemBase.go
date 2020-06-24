package ctrl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"wb/cc"
	"wb/cfg"
	"wb/cs"
	"wb/ctrl/uibuilder"
	"wb/ii"
	"wb/lg"
	"wb/md"
	"wb/om"
	"wb/st"
	"wb/svc"
	"wb/ut"
)

type ItemBaseController struct {
	BaseController

	TplList   string
	TplAdd    string
	TplUpdate string
	TplView   string

	UrlItemAdd    string
	UrlItemUpdate string

	UrlItemList string
	UrlUiAdd    string
	UrlUiUpdate string
}

func (ibc *ItemBaseController) NestPrepare() {
	lg.Debug("ItemBaseController.NestPrepare", " version: ", ibc.CtxVersion)

	ibc.TplList = "item/list.tpl"
	ibc.TplAdd = "item/edit.tpl"
	ibc.TplUpdate = "item/edit.tpl"
	ibc.TplView = "item/view.tpl"

	ibc.UrlItemList = ""
	ibc.UrlUiAdd = ""
	ibc.UrlUiUpdate = ""
}
func (ibc *ItemBaseController) GetItem(oItemInfo ii.ItemInfo) (string, cs.MObject) {
	queryParams, _, _, _ := ibc.GetListParams(oItemInfo)
	return svc.Svc(ibc.CtxUser).Get(oItemInfo.Name, queryParams)
}
func (ibc *ItemBaseController) GetTableList(oItemInfo ii.ItemInfo, extParam om.Params) cs.TableResult {
	listParams := ibc.GetSvcListParams(oItemInfo)
	for k, v := range extParam {
		listParams.Query[k] = v
	}
	return svc.Svc(ibc.CtxUser).GetTableList(oItemInfo.Name, listParams)
}

func (ibc *ItemBaseController) AddItem(oItemInfo ii.ItemInfo) {
	svcParams := ibc.GetFormValues(oItemInfo)
	lg.Debug("ItemBaseController.AddItem:", svcParams)
	subItemMaps, _ := ibc.GetRequestMapFromJson()
	for k, v := range subItemMaps {
		svcParams[k] = v
	}
	lg.Debug("ItemBaseController.AddItem:", svcParams)
	if sNum := ibc.GetString("_num"); sNum != "" {
		if num, ok := ut.ToInt64(sNum); ok {
			rows := make([]map[string]interface{}, 0)
			for i := 0; i < int(num); i++ {
				row := ut.Maps.Copy(svcParams)
				row[cc.Col.Sn] = ut.TUId()
				rows = append(rows, row)
			}
			ret, res := svc.Svc(ibc.CtxUser).AddItems(oItemInfo, rows)
			ibc.SendJson(&cs.JsonResult{ret, res})
		} else {
			ibc.SendJson(&cs.JsonResult{st.FieldFormatError, "_num"})
		}
	} else {
		status, reason := svc.Svc(ibc.CtxUser).AddItem(oItemInfo.Name, svcParams)
		ibc.SendJson(&cs.JsonResult{status, reason})
	}
}
func (ibc *ItemBaseController) AddItems(oItemInfo ii.ItemInfo) {
	rows, _ := ibc.getRequestRowsFromJson()
	if len(rows) <= 0 {
		ibc.SendJson(&cs.JsonResult{st.RowsNotFound, st.RowsNotFound})
		return
	}
	result, desc := svc.Svc(ibc.CtxUser).AddItems(oItemInfo, rows)
	ibc.SendJson(&cs.JsonResult{result, desc})
}
func (ibc *ItemBaseController) AddItemWithSub(oItemInfo ii.ItemInfo) {
	svcParams := ibc.GetFormValues(oItemInfo)
	subItemMaps, _ := ibc.GetRequestMapFromJson()
	for k, v := range svcParams {
		subItemMaps[k] = v
	}
	result, desc := svc.Svc(ibc.CtxUser).AddItem(oItemInfo.Name, subItemMaps)
	ibc.SendJson(&cs.JsonResult{result, desc})
}
func (ibc *ItemBaseController) UpdateItem(oItemInfo ii.ItemInfo) {
	ibc.UpdateItemEx(oItemInfo, om.Params{})
}
func (ibc *ItemBaseController) UpdateItemEx(oItemInfo ii.ItemInfo, exParam om.Params) {
	svcParams := ibc.GetFormValues(oItemInfo)
	for k, v := range exParam {
		svcParams[k] = v
	}
	sn := ut.Maps.GetString(svcParams, cc.Col.Sn)
	if strings.EqualFold(sn, "") {
		lg.Error(fmt.Sprintf("Item %s Sn not define", oItemInfo.Name))
		ibc.SendJson(&cs.JsonResult{st.ItemNotDefine, st.ItemNotDefine})
		return
	}
	sns := strings.Split(sn, ",")
	if len(sns) <= 1 {
		status, reason := om.Table(oItemInfo.Name).Update(svcParams)
		ibc.SendJson(&cs.JsonResult{status, reason})
	}
	delete(svcParams, cc.Col.Sn)
	status, reason := om.Table(oItemInfo.Name).BatchUpdate(svcParams, cc.Col.Sn, sns)
	ibc.SendJson(&cs.JsonResult{status, reason})
}
func (ibc *ItemBaseController) UpdateItemWithSub(oItemInfo ii.ItemInfo) {
	svcParams := ibc.GetFormValues(oItemInfo)
	subItemMaps, _ := ibc.GetRequestMapFromJson()
	for k, v := range subItemMaps {
		svcParams[k] = v
	}
	sn := ut.Maps.GetString(svcParams, cc.Col.Sn)
	if strings.EqualFold(sn, "") {
		lg.Error(fmt.Sprintf("Item %s Sn not define", oItemInfo.Name))
		ibc.SendJson(&cs.JsonResult{st.SnNotFound, oItemInfo.Name})
		return
	}
	ret, res := svc.Svc(ibc.CtxUser).UpdeteItem(oItemInfo.Name, svcParams)
	ibc.SendJson(&cs.JsonResult{ret, res})
}
func (ibc *ItemBaseController) UpdateItemWithAddSub(oItemInfo ii.ItemInfo) {
	svcParams := ibc.GetFormValues(oItemInfo)
	subItemMaps, _ := ibc.GetRequestMapFromJson()
	for k, v := range subItemMaps {
		svcParams[k] = v
	}
	sn := ut.Maps.GetString(svcParams, cc.Col.Sn)
	if strings.EqualFold(sn, "") {
		lg.Error(fmt.Sprintf("Item %s Sn not define", oItemInfo.Name))
		ibc.SendJson(&cs.JsonResult{st.SnNotFound, oItemInfo.Name})
		return
	}
	ret, des := svc.Svc(ibc.CtxUser).DeleteItem(oItemInfo.Name, sn)
	if ret != st.Success {
		lg.Error("svc.UpdateItem error:", des)
		ibc.SendJson(&cs.JsonResult{ret, des})
	}
	status, reason := svc.Svc(ibc.CtxUser).AddItem(oItemInfo.Name, svcParams)
	ibc.SendJson(&cs.JsonResult{status, reason})
}
func (ibc *ItemBaseController) DeleteItem(oItemInfo ii.ItemInfo) {
	svcParams := ibc.GetFormValues(oItemInfo)
	sn := ut.Maps.GetString(svcParams, cc.Col.Sn)
	if sn == "" {
		lg.Error(fmt.Sprintf("Item %s Sn not define", oItemInfo.Name))
		ibc.SendJson(&cs.JsonResult{st.SnNotFound, oItemInfo.Name})
		return
	}
	status, reason := svc.Svc(ibc.CtxUser).Delete(oItemInfo.Name, sn)
	ibc.SendJson(&cs.JsonResult{status, reason})
}
func (ibc *ItemBaseController) UploadItem(oItemInfo ii.ItemInfo) {
	sn := ibc.GetString(cc.Col.Sn)
	if sn == "" {
		lg.Error("ItemController.Upload error: ", st.SnNotFound)
		ibc.Ctx.WriteString(st.SnNotFound)
	}
	f, h, e := ibc.GetFile("fileUpload")

	if e != nil {
		lg.Error("Upload error", e.Error())
		ibc.Ctx.WriteString("{error}")
		return
	}
	f.Close()
	saveDir := fmt.Sprintf("%s/static/files/%s/%s/", cfg.WbConfig.DataPath, oItemInfo.Name, sn)
	err := os.MkdirAll(saveDir, 0777)
	if err != nil {
		lg.Error("ItemController.Upload error: ", st.UploadErrorCreateDir)
		ibc.Ctx.WriteString(st.UploadErrorCreateDir)
		return
	}
	ibc.SaveToFile("fileUpload", saveDir+h.Filename)
	ibc.Ctx.WriteString("{}")
}

func (ibc *ItemBaseController) AutocompleteItem(oItemInfo ii.ItemInfo) {
	keyword := ibc.GetString(cc.Keyword)
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		keyword = cc.Col.Name
	}
	term := ibc.GetString(cc.Term)
	if strings.EqualFold(term, "") {
		ibc.Data["json"] = "[]"
		ibc.ServeJSON()
		return
	}
	limitParams := om.LimitParams{
		om.Limit: om.LimitDefault,
	}

	orderByParams := om.OrderParams{
		keyword: om.Asc,
	}
	addParams := ibc.GetFormValues(oItemInfo)
	delete(addParams, cc.Keyword)
	addParams["%"+keyword] = term
	_, _, resultMaps := om.Table(oItemInfo.Name).ListWithParams(addParams, om.Params{}, limitParams, orderByParams)
	retList := md.TransAutocompleteList(resultMaps, keyword)
	ibc.SendJson(&retList)
}

func (ibc *ItemBaseController) UiListItem(oItemInfo ii.ItemInfo) {
	itemName := oItemInfo.Name
	ibc.Data["item"] = itemName

	ibc.SetDataWithoutDefaultStr(cc.UrlItemList, ibc.UrlItemList, "/item/list/"+itemName)
	if ibc.Ctx.Request.URL.RawQuery != "" {
		ibc.Data["UrlItemList"] = ibc.Data["UrlItemList"].(string) + "?" + ibc.Ctx.Request.URL.RawQuery
	}
	ibc.SetDataWithoutDefaultStr(cc.UrlUiAdd, ibc.UrlUiAdd, "/ui/add/"+itemName)
	ibc.SetDataWithoutDefaultStr(cc.UrlUiUpdate, ibc.UrlUiUpdate, "/ui/update/"+itemName)
	ibc.Data["thlist"] = uibuilder.BuildListThs(oItemInfo)
	customTpl := cfg.WebConfig.ViewsPath + "/" + oItemInfo.Name + "/" + "list.tpl"
	if ut.IsPathExist(customTpl) {
		ibc.TplName = oItemInfo.Name + "/" + "list.tpl"
	} else {
		ibc.TplName = ibc.TplList
	}
}

func (ibc *ItemBaseController) UiAddItem(oItemInfo ii.ItemInfo) {
	itemName := oItemInfo.Name
	ibc.SetDataWithoutDefaultStr(cc.UrlService, ibc.UrlItemAdd, "/item/add/"+itemName)
	sn := ut.TUId()
	ibc.Data["Sn"] = sn
	ibc.Data["Form"] = uibuilder.BuildAddForm(oItemInfo, sn)
	ibc.Data["Onload"] = uibuilder.BuildAddOnLoadJs(oItemInfo)
	customTpl := cfg.WebConfig.ViewsPath + "/" + oItemInfo.Name + "/" + "add.tpl"
	if ut.IsPathExist(customTpl) {
		vMap := map[string]interface{}{}
		for _, field := range oItemInfo.Fields {
			vMap[field.Name] = field.GetDefaultValue()
		}
		vMap[cc.Col.Sn] = ut.TUId()
		ibc.FillFormElement(uibuilder.BuildFormElement(oItemInfo, vMap, map[string]string{}))
		ibc.TplName = oItemInfo.Name + "/" + "add.tpl"
	} else {
		ibc.TplName = ibc.TplAdd
	}
	lg.Debug("UiAddItem tpl", ibc.TplName)
}
func (ibc *ItemBaseController) UiUpdateItemWithStatus(oItemInfo ii.ItemInfo, statusMap map[string]string) cs.MObject {
	itemName := oItemInfo.Name
	sn := ibc.GetString(cc.Col.Sn)
	if sn == "" {
		lg.Error("ui.Update", st.ParamSnIsNone)
		ibc.Ctx.WriteString(st.ParamSnIsNone)
		return cs.MObject{}
	}
	params := om.Params{cc.Col.Sn: sn}
	c, oldValueMap := om.Table(itemName).Get(params)
	if c == st.Success {
		oItemInfo.FormatDateTimeInMObject(oldValueMap)
		ibc.SetDataWithoutDefaultStr(cc.UrlService, ibc.UrlItemUpdate, "/item/update/"+itemName)
		ibc.Data["Sn"] = sn
		ibc.Data["Form"] = uibuilder.BuildUpdatedForm(oItemInfo, oldValueMap)
		ibc.Data["Onload"] = uibuilder.BuildUpdateOnLoadJs(oItemInfo)
		customTpl := cfg.WebConfig.ViewsPath + "/" + oItemInfo.Name + "/" + "update.tpl"
		if ut.IsPathExist(customTpl) {
			ibc.FillFormElement(uibuilder.BuildFormElement(oItemInfo, oldValueMap, statusMap))
			ibc.FillValue(oldValueMap)
			ibc.TplName = oItemInfo.Name + "/" + "update.tpl"
		} else {
			ibc.TplName = ibc.TplUpdate
		}
		lg.Debug("UiUpdateItem tpl", ibc.TplName)
	} else {
		lg.Error("UiUpdateItem error", st.ItemNotFound)
		ibc.Ctx.WriteString(st.ItemNotFound)
	}
	return oldValueMap
}
func (ibc *ItemBaseController) UiUpdateItem(oItemInfo ii.ItemInfo) {
	ibc.UiUpdateItemWithStatus(oItemInfo, map[string]string{})
}
func (ibc *ItemBaseController) UiViewItem(oItemInfo ii.ItemInfo) {
	itemName := oItemInfo.Name
	sn := ibc.GetString(cc.Col.Sn)
	if sn == "" {
		lg.Error("UiViewItem:", st.ParamSnIsNone)
		ibc.Ctx.WriteString(st.ParamSnIsNone)
		return
	}
	params := om.Params{cc.Col.Sn: sn}
	c, oldValueMap := om.Table(itemName).Get(params)
	if c == st.Success {
		ibc.SetDataWithoutDefaultStr(cc.UrlService, ibc.UrlItemUpdate, "/item/update/"+itemName)
		ibc.Data["Sn"] = sn
		oldValueMap = oItemInfo.FormatMObject(oldValueMap)
		ibc.Data["ViewPanels"] = uibuilder.BuildViewPanelFromItemInfo(oItemInfo, oldValueMap, "")
		customTpl := cfg.WebConfig.ViewsPath + "/" + oItemInfo.Name + "/" + "view.tpl"
		if ut.IsPathExist(customTpl) {
			ibc.TplName = oItemInfo.Name + "/" + "view.tpl"
		} else {
			ibc.TplName = ibc.TplView
		}
	} else {
		ibc.Ctx.WriteString(st.ItemNotFound)
	}
}
func (ibc *ItemBaseController) UiAttachmentItem(oItemInfo ii.ItemInfo) {
	sn := ibc.GetString(cc.Col.Sn)
	dir := fmt.Sprintf("%s/static/files/%s/%s/", cfg.WbConfig.DataPath, oItemInfo.Name, sn)
	urlDir := fmt.Sprintf("static/files/%s/%s/", oItemInfo.Name, sn)
	fileList := ut.ListDir(dir, true)
	str := ""
	for _, path := range fileList {
		str += "<a href=\"/" + urlDir + path + "\">" + path + "</a><br>"
		pathLow := strings.ToLower(path)
		if strings.HasSuffix(pathLow, ".png") || strings.HasSuffix(pathLow, ".jpg") || strings.HasSuffix(pathLow, ".gif") {
			str += `<img src="/` + urlDir + path + ` "alt="` + path + `" /><br>`
		}
	}
	ibc.Data["FileList"] = str
	ibc.TplName = "item/attachment.tpl"
}

// TODO deprecated
func (ibc *ItemBaseController) GetListPostParams(oItemInfos ...ii.ItemInfo) (queryParams om.Params, limitParams om.LimitParams, orderByParams om.OrderParams, searchText string) {
	requestMap, _ := ibc.GetRequestMapFromJson()
	lg.Debug("BaseController.GetListParams requestMap:", requestMap)
	if search, ok := requestMap["search"]; ok {
		searchText = ut.ToStr(search)
		delete(requestMap, "search")
	}
	limitParams = ibc.getLimitParamFromJsonMap(requestMap)
	delete(requestMap, om.Limit)
	delete(requestMap, om.Offset)
	orderByParams = ibc.getOrderParamFromJsonMap(requestMap)
	delete(requestMap, om.Order)
	delete(requestMap, om.Sort)
	filterParams := map[string]interface{}{}
	if filterStr, ok := requestMap["_filter"]; ok {
		filterMap := GetFilterMap(filterStr)
		delete(requestMap, "_filter")
		if len(oItemInfos) > 0 {
			filterParams = ibc.getFilteredQueryParamFromJsonMap(oItemInfos[0], filterMap)
		} else {
			filterParams = ibc.getQueryParamFromJsonMap(filterMap)
		}
	}
	if len(oItemInfos) > 0 {
		queryParams = ibc.getFilteredQueryParamFromJsonMap(oItemInfos[0], requestMap)
	} else {
		queryParams = ibc.getQueryParamFromJsonMap(requestMap)
	}
	for k, v := range filterParams {
		if _, ok := queryParams[k]; !ok {
			queryParams[k] = v
		}
	}
	lg.Debug("BaseController.GetListParams queryParams:", queryParams)
	return
}
func GetFilterMap(src interface{}) (p om.Params) {
	if s, ok := src.(string); !ok {
		return p
	} else {
		if err := json.Unmarshal([]byte(s), &p); err == nil {
			return p
		}
		return p
	}
}
func (ibc *ItemBaseController) GetSvcListParams(oItemInfo ii.ItemInfo, searchFields ...string) svc.ListParams {
	queryParams, searchParams, limitParams, orderByParams := ibc.GetListParams(oItemInfo, searchFields...)
	return svc.ListParams{queryParams, searchParams, limitParams, orderByParams}
}
func (ibc *ItemBaseController) GetListParams(oItemInfo ii.ItemInfo, searchFields ...string) (queryParams, searchParams om.Params, limitParams om.LimitParams, orderByParams om.OrderParams) {
	queryParams, limitParams, orderByParams, searchText := ibc.GetListPostParams(oItemInfo)
	searchParams = om.Params{}

	if len(searchFields) == 0 {
		searchFields = append(searchFields, "name")
	}
	for _, searchField := range searchFields {
		if _, ok := oItemInfo.GetField(searchField); ok {
			searchParams[searchField] = searchText
		}
	}
	fromValue := ibc.GetFormValues(oItemInfo)
	for k, v := range fromValue {
		queryParams[k] = v
	}
	return queryParams, searchParams, limitParams, orderByParams
}
func (ibc *ItemBaseController) GetParamItemList(addParam om.Params, oItemInfo ii.ItemInfo) (string, int64, []cs.MObject) {
	queryParams, searchParams, limitParams, orderByParams := ibc.GetListParams(oItemInfo)
	for k, v := range addParam {
		queryParams[k] = v
	}
	return om.Table(oItemInfo.Name).ListWithParams(queryParams, searchParams, limitParams, orderByParams)
}
func (ibc *ItemBaseController) GetFormValues(oItemInfo ii.ItemInfo) map[string]interface{} {
	fromValues, _ := ibc.getExFormValues(oItemInfo)
	lg.Debug("BaseController.GetFormValues formValues: ", fromValues)
	return fromValues
}
func (ibc *ItemBaseController) getExFormValues(oItemInfo ii.ItemInfo) (map[string]interface{}, map[string]interface{}) {
	retMap := make(map[string]interface{})
	exMap := make(map[string]interface{})
	formValues := ibc.Input()
	lg.Debug("BaseController.GetFormValues from values: ", formValues)
	for k, _ := range formValues {
		if field, ok := oItemInfo.GetField(k); ok {
			switch field.Model {
			case ii.MUpload, ii.MSubItem:
				continue
			}
			strValue := ibc.GetString(field.Name)
			rs, _ := ibc.replaceSpecialValues(strValue)
			if fromValue, fok := field.GetValidValue(rs); fok {
				retMap[field.Name] = fromValue
			}
		} else {
			exMap[k] = ibc.GetString(k)
		}
	}
	return retMap, exMap
}
func (ibc *ItemBaseController) GetItemInfoByName(itemName string) (ii.ItemInfo, string) {
	oItemInfo, ok := ii.ItemInfoMap[itemName]
	if !ok {
		lg.Error(st.ItemNotDefine_, itemName)
		return ii.ItemInfo{}, st.ItemNotDefine
	}
	return oItemInfo, st.Success
}

func (ibc *ItemBaseController) AddSpecialModelValues(oii ii.ItemInfo, m map[string]interface{}) map[string]interface{} {
	for _, f := range oii.Fields {
		switch f.Model {
		case ii.MCurUserName:
			m[f.Name] = ibc.CtxUser.GetName()
		case ii.MCurUserSn:
			m[f.Name] = ibc.CtxUser.GetSn()
		case ii.MCurTime:
			m[f.Name] = ut.GetCurDbTime()
		}
	}
	return m
}
func (ibc *ItemBaseController) GetAttachmentList(retList []cs.MObject, itemName string) {
	attList := ut.ListDir("static/files/"+itemName, false)
	for _, vm := range retList {
		sn := vm.GetString("sn")
		if ut.IsStrInSlice(attList, sn) {
			vm["atmt"] = true
		} else {
			vm["atmt"] = false
		}
	}
}
