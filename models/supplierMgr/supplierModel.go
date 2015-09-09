package supplierMgr
import (
	"webo/models/svc"
	"webo/models/s"
	"strings"
	"webo/models/stat"
)


func Get(sn string)(map[string]interface{}, bool){
	params := svc.Params{
		s.Sn : sn,
	}
	status, retMap := svc.Get(s.Supplier, params)
	return retMap, strings.EqualFold(stat.Success, status)
}
