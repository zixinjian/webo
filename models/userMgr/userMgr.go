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


type User struct {
	Sn          string
	UserName	string
	Name 		string
	Role		string
	Department	string
	Flag        string
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


func GetUsersByDepartment(department string) []User{
	queryParams := t.Params{
		s.Department: department,
	}
	orderParams := t.Params{
		s.Name: s.Asc,
	}
	// TODO
	if code, userMaps := svc.GetItems(s.User, queryParams, orderParams, t.LimitParams{}); strings.EqualFold(code, stat.Success) {
		userList := make([]User, 0)
		for _, user := range userMaps {
			f, _ := user[s.Flag].(string)
			if f != s.FlagAvailable{
				continue
			}
			v, _ := user[s.Sn].(string)
			u, _ := user[s.UserName].(string)
			l, _ := user[s.Name].(string)
			r, _ := user[s.Role].(string)
			d, _ := user[s.Department].(string)
			userList = append(userList, User{v, u, l, r, d, f})
		}
		return userList
	} else {
		return make([]User, 0)
	}
}
