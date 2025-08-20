package goat

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Template represents an email template
type Template struct {
	Name       string      // Name of the template
	ContentRaw string      // Raw HTML content with potential template variables
	Data       interface{} // Data containing the template variables values
}

// Render renders a template with the given data
func (t Template) Render() (string, error) {

	// Render the parent template
	tmpl, err := template.New(t.Name).Parse(t.ContentRaw)
	if err != nil {
		return "", err
	}

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, t.Data)
	if err != nil {
		return "", err
	}

	// Check if the rendered template contains <no value>
	if strings.Contains(buf.String(), "<no value>") {
		return "", fmt.Errorf("rendered template contains <no value>")
	}

	// Return the rendered template
	return buf.String(), nil
}
