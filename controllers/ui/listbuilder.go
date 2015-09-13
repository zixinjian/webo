package ui

import (
	"fmt"
	"strings"
	"webo/models/itemDef"
)

const thFormat = `                <th data-field="%s" %s %s>%s</th>
`

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
					sortable = sortable + ` data-order="desc"`
				}
			}
			th = th + fmt.Sprintf(thFormat, field.Name, visible, sortable, field.Label)
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

func BuildColums(itemDef itemDef.ItemDef) string {
	var columns = []string{}
	for _, field := range itemDef.Fields {
		if !field.UiList.Shown {
			continue
		}
		column := fmt.Sprintf(columnFormat, field.Name, field.Label)
		if !field.UiList.Visiable {
			column = column + ",visible:false\n"
		}
		if field.UiList.Sortable {
			column = column + ",sortable:\"true\""
			if field.UiList.Order == "desc" {
				column = column + ",order:\"desc\""
			}
		}
		if field.IsEditable() {
			//			column = column + ",editable: {"
			switch field.Input {
			case "text", "textarea":
				column = column + fmt.Sprintf(",editable:{type:\"%s\"}\n", field.Input)
				//				fmt.Println()
			case "select":
				srcs := []string{}
				for _, v := range field.Enum {
					srcs = append(srcs, fmt.Sprintf(`{value: "%s",text: "%s"}`, v.Sn, v.Label))
				}
				column = column + fmt.Sprintf(`,editable:{type:"%s",source:[%s]}`, field.Input, strings.Join(srcs, ",\n"))
			}
		}
		columns = append(columns, fmt.Sprintf(`{%s}`, column))
	}
	return fmt.Sprintf(`<script>columns = [%s]</script>`, strings.Join(columns, ","))
}
