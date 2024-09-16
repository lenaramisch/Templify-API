package domain

type WorkflowInfo struct {
	Name           string
	RequiredInputs []struct {
		ToEmail       string
		ToName        string
		EmailTemplate struct {
			TemplateName string
			Placeholders []string
		}
		PdfTemplates []struct {
			TemplateName string
			Placeholders []string
		}
	}
}
