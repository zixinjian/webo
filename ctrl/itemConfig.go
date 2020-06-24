package ctrl

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"sort"
	"wb/ii"
	"wb/lg"
	"wb/st"
)

type ItemConfigController struct {
	BaseController
}

func (this *ItemConfigController) Get() {
	itemNames := make([]string, 0, len(ii.ItemInfoMap))
	for k, _ := range ii.ItemInfoMap {
		itemNames = append(itemNames, k)
	}
	sort.Strings(itemNames)
	this.Data["ItemNames"] = itemNames
	this.TplName = "item/itemconf.tpl"
}
func (this *ItemConfigController) GetItemConfig() {
	name := this.GetString("name")
	if oItemInfo, ok := ii.ItemInfoMap[name]; ok {
		this.SendJson(oItemInfo)
	} else {
		this.SendJson("no such name")
	}
}
func (this *ItemConfigController) UpdateItemConfig() {
	requestBody := this.Ctx.Input.RequestBody

	var oItemInfo ii.ItemInfo
	json.Unmarshal(requestBody, &oItemInfo)
	lg.Debug("ItemConfigController.UpdateItemConfig", string(requestBody), oItemInfo)
	for idx, f := range oItemInfo.Fields {
		if f.Require == "t" {
			f.Require = "true"
		} else {
			f.Require = "false"
		}
		if f.Unique == "t" {
			f.Unique = "true"
		} else {
			f.Unique = "false"
		}
		if f.Ext.Item == "" {
			f.Ext.Fields = nil
		}
		if f.Name == "creater" {
			f.Unique = "false"
		}
		if f.Name == "createtime" {
			f.Unique = "false"
		}
		oItemInfo.Fields[idx] = f
	}
	xmlOutPut, outPutErr := xml.MarshalIndent(oItemInfo, "  ", "    ")
	if outPutErr == nil {
		//加入XML头
		headerBytes := []byte(xml.Header)
		//拼接XML头和实际XML内容
		xmlOutPutData := append(headerBytes, xmlOutPut...)
		//写入文件
		ioutil.WriteFile("newItemConfig/"+oItemInfo.Name+".xml", xmlOutPutData, os.ModeAppend)
		ii.CreateSql(oItemInfo, "newItemConfig/sql/"+oItemInfo.Name+".sql")
		this.SendJson(map[string]string{"status": st.Success})
	} else {
		lg.Error(outPutErr)
		this.SendJson(map[string]string{"status": st.Failed})
	}
}
