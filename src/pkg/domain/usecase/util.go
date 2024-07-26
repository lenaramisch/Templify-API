package usecase

import (
	"bytes"
	"log/slog"
	"regexp"
	domain "templify/pkg/domain/model"
	"text/template"
)

func ExtractPlaceholders(template string) []string {
	reg := regexp.MustCompile(`{{\s*\.([a-zA-Z]+)\s*}}`)
	matches := reg.FindAllStringSubmatch(template, -1)
	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}
	return placeholders
}

func FillTemplate(templateStr string, placeholderValues map[string]string) (string, error) {
	requiredPlaceholders := ExtractPlaceholders(templateStr)
	for _, requiredPlaceholder := range requiredPlaceholders {
		if _, ok := placeholderValues[requiredPlaceholder]; !ok {
			slog.With(
				"RequiredPlaceholder", requiredPlaceholder,
			).Debug("Missing placeholder")
			return "", domain.ErrorPlaceholderMissing{
				MissingPlaceholder: requiredPlaceholder,
			}
		}
	}
	templ, err := template.New("someName").Parse(templateStr)
	if err != nil {
		slog.With(
			"TemplateStr", templateStr,
		).Debug("Error parsing template")
		return "", err // domain error?
	}
	buf := &bytes.Buffer{}
	err = templ.Execute(buf, placeholderValues)
	if err != nil {
		slog.With(
			"Template", templateStr,
			"Values", placeholderValues,
		).Debug("Error executing template")
		return "", err
	}
	return buf.String(), nil
}
