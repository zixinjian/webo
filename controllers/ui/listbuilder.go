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
		if field.UiList.Shown {
			visible := ""
			if !field.UiList.Visiable {
				visible = `data-visible="false"`
			}
			sortable := ""
			if field.UiList.Sortable {
				sortable = `data-sortable="true"`
				if field.UiList.Order == "desc" {
					sortable = sortable + `data-order="desc"`
				}
			}
			th = th + fmt.Sprintf(`<th data-field="%s" %s %s>%s</th>`, field.Name, visible, sortable, field.Label)
		}
	}
	return th
}
