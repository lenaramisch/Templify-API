package mjmlservice

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"regexp"

	"example.SMSService.com/pkg/domain"
)

type MJMLConfig struct {
	//no config variables for now
}

type MJMLService struct {
	config MJMLConfig
}

func NewMJMLService(config MJMLConfig) *MJMLService {
	return &MJMLService{
		config: config,
	}
}

func extractPlaceholders(MJMLTemplateString string) []string {
	reg := regexp.MustCompile(`{{\s*\.([a-zA-Z]+)\s*}}`)
	matches := reg.FindAllStringSubmatch(MJMLTemplateString, -1)

	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}

// GET /templates/{template-name}
func (m *MJMLService) GetTemplatePlaceholders(domainTemplate domain.Template) ([]string, error) {
	MJMLTemplateString := domainTemplate.MJMLString
	placeholders := extractPlaceholders(MJMLTemplateString)
	if len(placeholders) == 0 {
		fmt.Printf("No placeholders found")
		return []string{}, nil
	}

	return placeholders, nil
}

func (m *MJMLService) FillTemplatePlaceholders(domainTemplate domain.Template, values map[string]interface{}) (string, error) {
	MJMLTemplateString := domainTemplate.MJMLString
	placeholders, err := m.GetTemplatePlaceholders(domainTemplate)
	if err != nil {
		fmt.Print("Getting placeholders for template failed")
		return "Getting placeholders for template failed", err
	}
	for _, placeholder := range placeholders {
		if _, ok := values[placeholder]; !ok {
			return "", fmt.Errorf("missing placeholder: %s", placeholder)
		}
	}
	templ, err := template.New("someName").Parse(MJMLTemplateString)
	if err != nil {
		fmt.Print("Error at template.New.Parse")
		return "", err
	}
	buf := &bytes.Buffer{}
	err = templ.Execute(buf, values)
	if err != nil {
		fmt.Print("Error at execute")
		return "", err
	}
	filledTemplate := buf.String()

	return filledTemplate, nil
}
