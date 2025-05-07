package utils

import (
	"bytes"
	"html/template"
	"log"
)

func ParseHtmlTemplate(path string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Fatalf("Error loading template: %v", err)
	}

	// Create a buffer to hold the rendered template
	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
		return "", err
	}

	// Print the email body
	return body.String(), nil
}
