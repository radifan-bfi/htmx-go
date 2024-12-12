package utils

import (
	"errors"
	"html/template"
	"reflect"
	"unicode"
)

// GetTemplateFuncMap returns the common template functions used across templates
func GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"last": func(x int, a interface{}) bool {
			maxLength := reflect.ValueOf(a).Len() - 1
			return x == maxLength
		},
		"toCapitalCase": func(s string) string {
			if s == "" {
				return s
			}
			var result []rune
			capitalize := true
			for i, r := range s {
				if i > 0 && unicode.IsUpper(r) {
					result = append(result, ' ')
				}
				if capitalize {
					r = unicode.ToUpper(r)
					capitalize = false
				}
				result = append(result, r)
			}
			return string(result)
		},
	}
}

// initializeTemplate is a helper function to create and register a template
func initializeTemplate(templates map[string]*template.Template, funcMap template.FuncMap, name string, files ...string) {
	tmpl := template.New(name).Funcs(funcMap)
	templates[name] = template.Must(tmpl.ParseFiles(files...))
}

// InitializeTemplates sets up all the templates with the common function map
func InitializeTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)
	funcMap := GetTemplateFuncMap()

	// Main templates with layout
	mainTemplates := map[string][]string{
		"form.html":           {"views/form.html", "views/layout.html"},
		"form-list.html":      {"views/form-list.html", "views/partials/form-table-partial.html", "views/layout.html"},
		"form-not-found.html": {"views/form-not-found.html", "views/layout.html"},
		"submissions.html":    {"views/submissions.html", "views/partials/submissions-partial.html", "views/layout.html"},
		"create-form.html":    {"views/create-form.html", "views/layout.html"},
	}

	// Partial templates
	partialTemplates := map[string][]string{
		"partials/form-table-partial.html":      {"views/partials/form-table-partial.html"},
		"partials/submissions-partial.html":     {"views/partials/submissions-partial.html"},
		"partials/alerts/error.html":            {"views/partials/alerts/error.html"},
		"partials/alerts/success.html":          {"views/partials/alerts/success.html"},
		"partials/alerts/success-redirect.html": {"views/partials/alerts/success-redirect.html"},
	}

	// Initialize main templates
	for name, files := range mainTemplates {
		initializeTemplate(templates, funcMap, name, files...)
	}

	// Initialize partial templates
	for name, files := range partialTemplates {
		initializeTemplate(templates, funcMap, name, files...)
	}

	return templates
}
