package main

import (
	"html/template"
	"log"
)

// GenerateTemplates generates and returns templates instance
func GenerateTemplates(templatesGlob string) *template.Template {
	funcMap := template.FuncMap{
		"noescape": func(str string) template.HTML { return template.HTML(str) },
	}

	templates, err := template.New("").Funcs(funcMap).ParseGlob(templatesGlob)

	if err != nil {
		log.Print("Error generating templates:")
		panic(err)
	}

	return templates
}
