package usecase

import (
	"fmt"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) SavePDF(fileName string, base64Content string) error {
	err := u.repository.SavePDF(fileName, base64Content)
	if err != nil {
		u.log.With("fileName", fileName).Debug("Error saving PDF")
		return err
	}
	return nil
}

func (u *Usecase) GetPDF(fileName string) (string, error) {
	pdf, err := u.repository.GetPDF(fileName)
	if err != nil {
		u.log.With("fileName", fileName).Debug("Error getting PDF")
		return "", err
	}
	return pdf, nil
}

func (u *Usecase) AddPDFTemplate(templateName string, typstString string) error {
	err := u.repository.AddPDFTemplate(templateName, typstString)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetPDFTemplateByName(templateName string) (*domain.Template, error) {
	templateDomain, err := u.repository.GetPDFTemplateByName(templateName)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return nil, err
	}
	return templateDomain, nil
}

func (u *Usecase) GetPDFPlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetPDFTemplateByName(templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GeneratePDF(templateName string, values map[string]string) ([]byte, error) {
	template, err := u.GetPDFTemplateByName(templateName)
	if err != nil {
		u.log.With(
			"TemplateName", templateName,
			"Error", err.Error(),
		).Debug("Error getting template during GeneratePDF")
		return nil, err
	}
	filledTemplate, err := FillTemplate(template.TemplateStr, values)
	if err != nil {
		u.log.With(
			"TemplateName", templateName,
			"Values", values,
			"Error", err.Error(),
		).Debug("Error filling template during GeneratePDF")
		return nil, err
	}

	u.log.With(
		"TemplateName", templateName,
		"FilledTemplate", filledTemplate,
	).Debug("Filled template")

	pdfFile, err := u.typstService.RenderTypst(filledTemplate)
	if err != nil {
		u.log.With(
			"Error", err.Error(),
		).Debug("Error using typst to generate PDF")
		return nil, err
	}
	return pdfFile, nil
}
