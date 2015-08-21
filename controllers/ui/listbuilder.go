package ui

import (
	"fmt"
	"webo/models/itemDef"
)

func BuildListThs(itemDef itemDef.ItemDef) string {
	th := ""
	for _, field := range itemDef.Fields {
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
