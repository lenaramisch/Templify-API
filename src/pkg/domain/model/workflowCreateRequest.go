package domain

type WorkflowCreateRequest struct {
	Name                string
	EmailTemplateName   string
	EmailTemplateString string
	IsMJML              bool
	StaticAttachments   []struct {
		Content  string
		FileName string
	}
	TemplatedPDFs []struct {
		TemplateName   string
		TemplateString string
	}
	EmailSubject string
}
