package u

import (
	"strings"
	"text/template"
	"bytes"
)

func IsNullStr(str interface{}) bool {
	return strings.EqualFold(str.(string), "")
}


func TemplateFormat(temp string, data interface{}) string{
	t := template.New("temp")
	t, _ = t.Parse(temp)
	b := new(bytes.Buffer)
	t.Execute(b, data)
	return b.String()
}