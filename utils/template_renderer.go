package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
)

type templateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer(templates map[string]*template.Template) *templateRenderer {
	return &templateRenderer{
		templates: templates,
	}
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}

	withoutTemplate := tmpl.Lookup("layout.html") == nil
	if withoutTemplate {
		if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
			log.Errorf(err.Error())
		}

		return nil
	}

	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}
