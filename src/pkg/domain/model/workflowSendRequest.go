package domain

type WorkflowSendRequest struct {
	toEmail           string
	toName            string
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
