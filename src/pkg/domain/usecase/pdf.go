package usecase

import (
	"context"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddPDFTemplate(ctx context.Context, template *domain.Template) error {
	err := u.repository.AddPDFTemplate(ctx, template)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetPDFTemplateByName(ctx context.Context, templateName string) (*domain.Template, error) {
	templateDomain, err := u.repository.GetPDFTemplateByName(ctx, templateName)
	if err != nil {
		return nil, err
	}
	return templateDomain, nil
}

func (u *Usecase) GetPDFPlaceholders(ctx context.Context, templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetPDFTemplateByName(ctx, templateName)
	if err != nil {
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GeneratePDF(ctx context.Context, templateName string, values map[string]string) ([]byte, error) {
	template, err := u.GetPDFTemplateByName(ctx, templateName)
	if err != nil {
		return nil, err
	}
	filledTemplate, err := FillTemplate(template.TemplateStr, values)
	if err != nil {
		return nil, domain.ErrorFillingTemplate{Reason: err.Error()}
	}

	pdfFile, err := u.typstService.RenderTypst(filledTemplate)
	if err != nil {
		return nil, domain.ErrorRenderingTypst{Reason: err.Error()}
	}
	return pdfFile, nil
}
