package usecase

import (
	"fmt"
	"log/slog"
	domain "templify/pkg/domain/model"
)

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

func (u *Usecase) FillPDFTemplatePlaceholders(templateName string, values map[string]string) (string, error) {
	pdfTempl, err := u.GetPDFTemplateByName(templateName)
	if err != nil {
		return "Getting template from db failed", err
	}
	return FillTemplate(pdfTempl.TemplateStr, values)
}

func (u *Usecase) GeneratePDF(templateName string, values map[string]string) ([]byte, error) {
	filledTemplate, err := u.FillPDFTemplatePlaceholders(templateName, values)
	if err != nil {
		slog.With(
			"TemplateName", templateName,
			"Values", values,
			"Error", err.Error(),
		).Debug("Error filling template during GeneratePDF")
		return nil, err
	}
	pdfFile, err := u.typstService.RenderTypst(filledTemplate)
	if err != nil {
		slog.With(
			"Error", err.Error(),
		).Debug("Error using typst to generate PDF")
		return nil, err
	}
	return pdfFile, nil
}
