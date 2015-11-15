package wbconf
import (
	"webo/models/svc"
	"webo/models/s"
	"webo/models/t"
	"webo/models/lang"
)

func LoadCategory(){
	_, categorys := svc.Query(s.Category, t.Params{}, t.LimitParams{}, t.Params{})
	for _, cate := range categorys {
		lang.AddLabel(cate[s.Key].(string), cate[s.Name].(string))
	}
}