package u

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
	"time"
	"math"
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
	if gId > 99 {
		gId = 0
		now.Add(time.Second)
	}
	gOldTime = now
	//    fmt.Println(now.Format("20060102150405"))
	ret := fmt.Sprintf("%s%02d", now.Format("20060102150405"), gId)
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

func GetStringValue(vMap map[string]interface{}, key string)string{
	if value, ok := vMap[key];ok {
		if v, vok := value.(string); vok {
			return v
		}
	}
	return ""
}

func GetToday() (totay string){
	t := time.Now()
	return t.Format("2006.01.02")
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
