package uibuilder

import (
	"fmt"
	"wb/ii"
)

const thFormat = `                <th data-field="%s" %s %s>%s</th>
`

func BuildListThs(oItemInfo ii.ItemInfo) string {
	th := ""
	for _, field := range oItemInfo.Fields {
		if field.Input != ii.INone {
			visible := ""
			//			if !field.UiList.Visiable {
			//				visible = `data-visible="false"`
			//			}
			sortable := ""
			//			if field.UiList.Sortable {
			//				sortable = `data-sortable="true"`
			//				if field.UiList.Order == "desc" {
			//					sortable = sortable + ` data-order="desc"`
			//				}
			//			}
			var fieldName string
			if field.Model == ii.MRelation {
				fieldName = field.Name + "_" + "name"
			} else {
				fieldName = field.Name
			}
			th = th + fmt.Sprintf(thFormat, fieldName, visible, sortable, field.Label)
		}
	}
	return th
}

const columnFormat = `	field:"%s",
	title:"%s"`
const optFormat = `	{field:"action",
	align:"center",
	formatter:"actionFormatter",
	events:"actionEvents",
	width:"75px"}
	`
