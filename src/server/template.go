package server

import (
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

var tmpl Template

func init() {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	templates := template.Must(template.New("templates").Funcs(funcMap).ParseGlob("./views/*.html.tpl"))

	tmpl = Template{
		templates: templates,
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
