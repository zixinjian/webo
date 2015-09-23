package wborm
import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strconv"
)


func GetCount(sql string, values []interface{})int64{
	o := orm.NewOrm()
	var maps []orm.Params
	if _, err := o.Raw(sql, values...).Values(&maps); err == nil {
		beego.Debug("GetCount :", maps)
		if len(maps) <= 0 {
			beego.Error("GetCount error: len(maps) < 0")
			return 0
		}
		if total, ok := maps[0]["count"]; ok {
			total64, err := strconv.ParseInt(total.(string), 10, 64)
			if err != nil {
				beego.Error("GetCount error: ", err)
				return 0
			}
			return total64
		}
	}else{
		beego.Error("GetCount error: ", err)
		return 0
	}
	beego.Error("GetCount error:")
	return 0
}