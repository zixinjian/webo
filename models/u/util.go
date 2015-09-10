package u

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"time"
)

var gId uint32
var gOldTime time.Time

func ToInt64(value interface{}) (d int64) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		panic(fmt.Errorf("ToInt64 need numeric not `%T`", value))
	}
	return
}

//func Str2Float(value string) float64{
//	f, err = ParseFloat(s, 32)
//}

//func ToString(value interface{})string{
//	val := reflect.ValueOf(value)
//	int64
//	switch value.(type) {
//	case int, int8, int16, int32, int64:
//		d = val.Int()
//	case uint, uint8, uint16, uint32, uint64:
//		d = int64(val.Uint())
//	default:
//		panic(fmt.Errorf("ToInt64 need numeric not `%T`", value))
//	}
//	return
//}

func TUId() string {
	now := time.Now()
	if gOldTime.After(now) {
		now = gOldTime
	}
	if gId > 999 {
		gId = 0
		now.Add(time.Second)
	}
	gOldTime = now
	//    fmt.Println(now.Format("20060102150405"))
	ret := fmt.Sprintf("%s%03d", now.Format("20060102150405"), gId)
	gId = gId + 1
	return ret
}

func init() {
	gId = 0
	gOldTime = time.Now()
}

func ToStr(v interface{}) string {
	return orm.ToStr(v)
}
