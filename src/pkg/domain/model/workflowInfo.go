package domain

type WorkflowInfo struct {
	Name           string
	EmailSubject   string
	RequiredInputs []struct {
		ToEmail           string
		ToName            string
		EmailTemplate     TemplateInfo
		PdfTemplates      []TemplateInfo
		StaticAttachments []string
	}
}
