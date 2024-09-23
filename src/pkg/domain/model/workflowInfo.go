package domain

type WorkflowInfo struct {
	Name           string
	RequiredInputs []struct {
		ToEmail           string
		ToName            string
		EmailTemplate     TemplateInfo
		PdfTemplates      []TemplateInfo
		StaticAttachments []string
	}
}
