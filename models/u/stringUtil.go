package u

import "strings"

func IsNullStr(str interface{}) bool {
	return strings.EqualFold(str.(string), "")
}

