package u

import "strings"

func StrJoin(strs []string, joint string) string {
	ret := ""
	for idx, s := range strs {
		if idx > 0 {
			ret = ret + joint
		}
		ret = ret + s
	}
	return ret
}

func IsNullStr(str interface{}) bool {
	return strings.EqualFold(str.(string), "")
}
