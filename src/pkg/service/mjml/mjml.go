package mjmlservice

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"example.SMSService.com/pkg/domain"
)

type MJMLConfig struct {
	Host string
	Port int
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

func (m *MJMLService) FillTemplatePlaceholders(domainTemplate domain.Template, values map[string]any) (string, error) {
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

func (m *MJMLService) RenderMJML(MJMLString string) (string, error) {
	//call MJML Service on Port 5000
	// Create a new POST request
	url := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(MJMLString)))
	if err != nil {
		slog.Debug("Error creating request")
		return "", err
	}

	req.Header.Set("Content-Type", "text/plain")

	//Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.With("err", err.Error()).Debug("Error sending request")
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Println("Error response from MJML Service:", string(body))
		return "", fmt.Errorf(fmt.Sprintf("MJML Service returned status code %d", resp.StatusCode))
	}

	htmlString := string(body)
	return htmlString, nil
}
