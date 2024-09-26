package domain

type WorkflowCreateRequest struct {
	Name                string
	EmailTemplateName   string
	EmailTemplateString string
	IsMJML              bool
	StaticAttachments   []PDF
	TemplatedPDFs       []Template
	EmailSubject        string
}
