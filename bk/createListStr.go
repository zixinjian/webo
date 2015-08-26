package main
import (
	"webo/models/itemDef"
	"fmt"
	"webo/controllers/ui"
)

func main(){
	fmt.Println("ok")
	fmt.Println(itemDef.EntityDefMap)
	for k, oItemDef := range itemDef.EntityDefMap{
		fmt.Println(k)
		if k != "user"{
			continue
		}
		fmt.Println(ui.BuildColums(oItemDef))
	}
}
