package user
import (
	"webo/models/svc"
	"webo/models/s"
)


func Get(sn string)(string, map[string]interface{}){
	params := svc.Params{
		s.Sn : sn,
	}
	return svc.Get(s.User, params)
}