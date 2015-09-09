package productMgr
import (
	"webo/models/stat"
	"webo/models/svc"
	"strings"
	"webo/models/s"
	"webo/models/t"
)

func Get(sn string)(map[string]interface{}, bool){
	params := t.Params{
		s.Sn : sn,
	}
	status, retMap := svc.Get(s.Product, params)
	return retMap, strings.EqualFold(stat.Success, status)
}