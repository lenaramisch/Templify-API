package mjmlservice

import (
	"fmt"
	"regexp"
)

type MJMLConfig struct {
	//no config variables for now
}

type MJMLService struct {
	config MJMLConfig
}

// TODO add template post request functionality
// POST /templates
// read MJML template given from user and save it
func (m *MJMLService) TemplatePostRequest(MJMLTemplate string) error {
	panic("unimplemented")
}

func NewMJMLService(config MJMLConfig) *MJMLService {
	return &MJMLService{
		config: config,
	}
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
	//TODO Get template by name from DB
	//TODO use go template execute to fill template with values
	//TODO return filled mjml template as string or error
	return "", nil
}
