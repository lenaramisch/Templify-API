package domain

type WorkflowCreateRequest struct {
	Name              string
	EmailTemplateName string
	StaticAttachments []string
	TemplatedPDFs     []string
	EmailSubject      string
}
