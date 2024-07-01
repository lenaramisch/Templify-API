package mjmlservice

import (
	"encoding/json"
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
	reg := regexp.MustCompile(`{{data:([a-zA-Z]+):""}}`)
	matches := reg.FindAllStringSubmatch(template, -1)

	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}

func createJSON(placeholders []string) (string, error) {
	data := make(map[string]string)
	for _, placeholder := range placeholders {
		data[placeholder] = ""
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// GET /templates/{template-name}
func (m *MJMLService) GetTemplatePlaceholders(template string) (string, error) {
	placeholders := extractPlaceholders(template)
	jsonPlaceholders, err := createJSON(placeholders)
	if err != nil {
		return "", fmt.Errorf("error creating JSON: %w", err)
	}

	return jsonPlaceholders, err
}
