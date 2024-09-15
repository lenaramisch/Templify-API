package domain

type WorkflowCreateRequest struct {
	Name              string
	EmailTemplateName string
	StaticAttachments []struct {
		Content  string
		FileName string
	}
	TemplatedPDFs []struct {
		TemplateName   string
		TemplateString string
	}
	EmailSubject string
}
