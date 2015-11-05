package wblog
import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

func LogItemDelete(item interface{}, deleter string) bool{
	log := logs.NewLogger(10000)
	log.SetLogger("file", `{"filename":"db/deletedItems.log"}`)
	if b, err := json.Marshal(item); err == nil{
		fmt.Println(b)
		beego.Warn("deleter: ", deleter, " item: ", string(b))
		log.Info("deleter: ", deleter, " item: ", string(b))
		return true
	}
	return false
}

