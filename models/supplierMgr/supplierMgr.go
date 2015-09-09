package supplierMgr
import (
	"webo/models/svc"
	"webo/models/s"
	"strings"
	"webo/models/stat"
	"webo/models/t"
)


func Get(sn string)(map[string]interface{}, bool){
	params := t.Params{
		s.Sn : sn,
	}
	status, retMap := svc.Get(s.Supplier, params)
	return retMap, strings.EqualFold(stat.Success, status)
}
