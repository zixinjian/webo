package userMgr

import (
	"webo/models/s"
	"webo/models/t"
	"webo/models/svc"
)

func Get(sn string) (string, map[string]interface{}) {
	params := t.Params{
		s.Sn: sn,
	}
	return svc.Get(s.User, params)
}
