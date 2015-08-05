package ui

import (
	"fmt"
	"webo/models/itemDef"
)

func BuildListThs(entity string) string {
	oEntityDef, ok := itemDef.EntityDefMap[entity]
	if !ok {
		fmt.Println("BuildAddForm none")
	}
	th := ""
	for _, field := range oEntityDef.Fields {
		if field.Input != "none" {
			th = th + fmt.Sprintf(`<th data-field="%s" data-sortable="true">%s</th>`, field.Name, field.Label)
		}
	}
	return th
}
