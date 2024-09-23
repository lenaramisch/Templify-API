package domain

type WorkflowUseRequest struct {
	Name          string
	EmailTemplate TemplateToFill
	PdfTemplate   TemplateToFill
	ToEmail       string
	ToName        string
}
