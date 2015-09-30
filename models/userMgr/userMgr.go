package userMgr

import (
	"webo/models/s"
	"webo/models/t"
	"webo/models/svc"
	"webo/models/itemDef"
	"strings"
	"webo/models/stat"
)

func Get(sn string) (string, map[string]interface{}) {
	params := t.Params{
		s.Sn: sn,
	}
	return svc.Get(s.User, params)
}


func GetUserEnum(queryParams t.Params, orderParams t.Params, limitParams t.LimitParams) []itemDef.EnumValue {
	if len(orderParams) <= 0{
		orderParams[s.Name]= s.Asc
	}
	if code, userMaps := svc.GetItems(s.User, queryParams, orderParams, limitParams); strings.EqualFold(code, stat.Success) {
		EnumList := make([]itemDef.EnumValue, len(userMaps))
		for idx, user := range userMaps {
			v, _ := user[s.Sn]
			u, _ := user[s.UserName]
			l, _ := user[s.Name]
			EnumList[idx] = itemDef.EnumValue{v.(string), u.(string), l.(string)}
		}
		return EnumList
	} else {
		return make([]itemDef.EnumValue, 0)
	}
}