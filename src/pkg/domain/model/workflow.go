package domain

type Workflow struct {
	Name              string
	EmailSubject      string
	EmailTemplateName string
	PDFTemplateNames  string
	StaticAttachments string
}
