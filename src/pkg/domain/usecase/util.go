package usecase

import (
	"bytes"
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
			return "", domain.ErrorPlaceholderMissing{
				MissingPlaceholder: requiredPlaceholder,
			}
		}
	}
	templ, err := template.New("someName").Parse(templateStr)
	if err != nil {
		return "", err // domain error?
	}
	buf := &bytes.Buffer{}
	err = templ.Execute(buf, placeholderValues)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ConvertPlaceholdersToSlice(placeholders map[string]*string) map[string]string {
	result := make(map[string]string)
	for key, value := range placeholders {
		result[key] = *value
	}
	return result
}
