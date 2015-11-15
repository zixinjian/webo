package controllers
import (
	"webo/models/svc"
	"webo/models/s"
	"webo/models/itemDef"
)


func getCategoryEnumList()[]itemDef.EnumValue{
	_, categorys := svc.GetAll(s.Category)
	enumList := make([]itemDef.EnumValue, len(categorys))
	for idx, category := range categorys{
		key := category[s.Key].(string)
		label := category[s.Name].(string)
		enumList[idx] = itemDef.EnumValue{key, key, label}
	}
	return enumList
}