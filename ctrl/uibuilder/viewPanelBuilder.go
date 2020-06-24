package uibuilder

import (
	"fmt"
	"wb/ii"
	"wb/lang"
	"wb/ut"
)

const viewPanelFormat = `    <a href="#%s"/>
<div class="panel panel-default m-t">
        <div class="panel-heading font-bold b-t">%s</div>
        <table class="table hover m-b-none bg-white">
            %s
        </table>
    </div>
`
const viewPanelTrFormat = `            <tr>
                <td class="bg-light lter b-r" width="30%%">%s</td>
                <td width="70%%">%s</td>
            </tr>
`

type attr struct {
	Label string
	Value string
}
type cate struct {
	Cate  string
	Attrs []attr
}

func BuildViewPanelFromItemInfo(oItemInfo ii.ItemInfo, oldValueMap map[string]interface{}, showCate string) string {
	cates := make([]cate, 0)
	for _, field := range oItemInfo.Fields {
		if field.Cate == "" {
			continue
		}
		if showCate != "" && field.Cate != showCate {
			continue
		}
		if field.Input == ii.INone {
			continue
		}
		valueString := ""
		if v, ok := oldValueMap[field.Name]; ok {
			valueString = ut.ToStr(v) + field.Unit
		}
		bExist := false
		for i, c := range cates {
			if c.Cate == field.Cate {
				c.Attrs = append(c.Attrs, attr{field.Label, valueString})
				bExist = true
				cates[i] = c
				break
			}
		}
		if !bExist {
			attrs := []attr{{field.Label, valueString}}
			c := cate{field.Cate, attrs}
			cates = append(cates, c)
		}
	}
	s := ""
	for _, c := range cates {
		s = s + "\n" + buildViewPanel(c)
	}
	return s
}
func buildViewPanel(c cate) string {
	viewPanelTrStr := ""
	for _, a := range c.Attrs {
		viewPanelTrStr += BuildViewTr(a.Label, a.Value)
	}
	title := lang.GetLabel(c.Cate)
	viewPanelStr := fmt.Sprintf(viewPanelFormat, c.Cate, title, viewPanelTrStr)
	return viewPanelStr
}
func BuildViewTr(key, attr string) string {
	return fmt.Sprintf(viewPanelTrFormat, key, attr)
}
