package usecase

import (
	"context"
	"fmt"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddPDFTemplate(ctx context.Context, template *domain.Template) error {
	err := u.repository.AddPDFTemplate(ctx, template)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetPDFTemplateByName(ctx context.Context, templateName string) (*domain.Template, error) {
	templateDomain, err := u.repository.GetPDFTemplateByName(ctx, templateName)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return nil, err
	}
	return templateDomain, nil
}

func (u *Usecase) GetPDFPlaceholders(ctx context.Context, templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetPDFTemplateByName(ctx, templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GeneratePDF(ctx context.Context, templateName string, values map[string]string) ([]byte, error) {
	template, err := u.GetPDFTemplateByName(ctx, templateName)
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
