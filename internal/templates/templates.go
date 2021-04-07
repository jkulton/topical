package templates

import (
	"html/template"
	"log"
)

// GenerateTemplates generates and returns templates instance
func GenerateTemplates(templatesGlob string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"noescape": func(str string) template.HTML { return template.HTML(str) },
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob(templatesGlob)

	if err != nil {
		log.Print("Error generating templates:")
		log.Print(err)
		return nil, err
	}

	return templates, nil
}
