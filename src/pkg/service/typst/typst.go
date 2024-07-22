package typst

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"text/template"

	"example.SMSService.com/pkg/domain"
)

type TypstConfig struct {
}

type TypstService struct {
	config TypstConfig
}

func NewTypstService(config TypstConfig) *TypstService {
	return &TypstService{
		config: config,
	}
}

func writeStringToFile(filledTemplStr string, templName string) (string, error) {
	typstFileName := templName + ".typ"
	f, err := os.Create(typstFileName)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	l, err := f.WriteString(filledTemplStr)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return "", err
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return typstFileName, nil
}

func renderTypst(typstFileName string) (string, error) {
	cmd := exec.Command("typst", "compile", typstFileName)
	cmd.Dir = "/Users/lenaramisch/Documents/GitHub/SMS-Service-API/src/pkg/service/typst/pdf-output"
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Print("Error rendering Typst File to PD ", err)
		return "", err
	}
	fmt.Println(out.String())
	return out.String(), nil
}

func extractPlaceholders(typstTemplString string) []string {
	reg := regexp.MustCompile(`{{\s*\.([a-zA-Z]+)\s*}}`)
	matches := reg.FindAllStringSubmatch(typstTemplString, -1)

	var placeholders []string
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}

func (ts *TypstService) GetPDFTemplatePlaceholders(typstString string) ([]string, error) {
	placeholders := extractPlaceholders(typstString)
	if len(placeholders) == 0 {
		fmt.Printf("No placeholders found")
		return []string{}, nil
	}

	return placeholders, nil
}

func (ts *TypstService) FillPDFTemplatePlaceholders(typstTempl *domain.PDFTemplate, values map[string]string) (string, error) {
	placeholders, err := ts.GetPDFTemplatePlaceholders(typstTempl.TypstString)
	if err != nil {
		fmt.Print("Getting placeholders for template failed")
		return "Getting placeholders for template failed", err
	}
	for _, placeholder := range placeholders {
		if _, ok := values[placeholder]; !ok {
			return "", fmt.Errorf("missing placeholder: %s", placeholder)
		}
	}
	templ, err := template.New("someName").Parse(typstTempl.TypstString)
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
	filledTemplateString := buf.String()
	typstFileName, err := writeStringToFile(filledTemplateString, typstTempl.Name)
	if err != nil {
		return "", err
	}
	result, err := renderTypst(typstFileName)
	if err != nil {
		return "", err
	}
	return result, nil
}
