package productMgr

import (
	"strings"
	"webo/models/s"
	"webo/models/stat"
	"webo/models/svc"
	"webo/models/t"
)

func Get(sn string) (map[string]interface{}, bool) {
	params := t.Params{
		s.Sn: sn,
	}
	status, retMap := svc.Get(s.Product, params)
	return retMap, strings.EqualFold(stat.Success, status)
}
