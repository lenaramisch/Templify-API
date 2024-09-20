package domain

type WorkflowUseRequest struct {
	Name          string
	EmailTemplate struct {
		Placeholders map[string]*string
		TemplateName string
	}
	PdfTemplate *struct {
		Placeholders map[string]*string
		TemplateName string
	}
	ToEmail string
	ToName  string
}
