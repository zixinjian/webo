package ui
import (
	"text/template"
	"bytes"
	"fmt"
)

type NavLiValue struct {
	Class   string
	Url 	string
	Text    string
}
const navLiTemplateStr = `                <li class="{{.Class}}"><a href="{{.Url}}" target="frame-content">{{.Text}}</a></li>
`
var navLiTemplate *template.Template

func BuildNavs(navValues []NavLiValue) string{
	navs := ""
	for _, navValue := range navValues{
		b := new(bytes.Buffer)
		fmt.Println(navLiTemplate.Execute(b, navValue))
		navs = navs + b.String()
	}
	return navs
}

func init(){
	t := template.New("controllers.ui.navLiTemplate")
	navLiTemplate, _ = t.Parse(navLiTemplateStr)
}