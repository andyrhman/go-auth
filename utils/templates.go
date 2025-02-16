package utils

import (
	"bytes"
	"html/template"
)

func ParseTemplate(file string, data interface{}) (string, error) {
	tpl, err := template.ParseFiles(file)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}