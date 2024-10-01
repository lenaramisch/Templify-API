package domain

type WorkflowUseRequest struct {
	Name          string
	EmailTemplate TemplateToFill
	PdfTemplates  []TemplateToFill
	ToEmail       string
	ToName        string
}
