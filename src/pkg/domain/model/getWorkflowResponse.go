package domain

type GetWorkflowResponse struct {
	Name              string
	EmailSubject      string
	EmailTemplate     TemplateInfo
	PDFTemplates      []TemplateInfo
	StaticAttachments []string
}
