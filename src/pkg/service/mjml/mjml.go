package mjmlservice

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"text/template"
)

//go:embed template_test.mjml
var hardCodedTemplate string

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

// TODO add template post to DB request functionality
// POST /templates
// read MJML template given from user and save it to DB
func (m *MJMLService) TemplatePostRequest(MJMLTemplate string) error {
	panic("unimplemented")
}

// TODO get template with template-name from DB
func extractPlaceholders(template string) []string {
	reg := regexp.MustCompile(`{{\s*\.([a-zA-Z]+)\s*}}`)
	matches := reg.FindAllStringSubmatch(template, -1)

	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}

// GET /templates/{template-name}
func (m *MJMLService) GetTemplatePlaceholders(template string) ([]string, error) {
	placeholders := extractPlaceholders(template)
	if len(placeholders) == 0 {
		fmt.Printf("No placeholders found")
		return []string{}, nil
	}

	return placeholders, nil
}

func (m *MJMLService) FillTemplatePlaceholders(templateName string, values map[string]interface{}) (string, error) {
	fmt.Print("Got values in MJML Service: ", values)
	placeholders, err := m.GetTemplatePlaceholders(hardCodedTemplate)
	if err != nil {
		fmt.Print("Getting placeholders for template failed")
		return "Getting placeholders for template failed", err
	}
	for _, placeholder := range placeholders {
		if _, ok := values[placeholder]; !ok {
			return "", fmt.Errorf("missing placeholder: %s", placeholder)
		}
	}
	//TODO Get template by name from DB
	//For now use hard coded template
	templ, err := template.New("someName").Parse(hardCodedTemplate)
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
