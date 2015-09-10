package user

import (
	"webo/models/s"
	"webo/models/svc"
)

func Get(sn string) (string, map[string]interface{}) {
	params := svc.Params{
		s.Sn: sn,
	}
	return svc.Get(s.User, params)
}
