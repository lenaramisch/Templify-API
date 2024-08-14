package domain

type PDFTemplateFillRequest struct {
	Placeholders map[string]string `json:"placeholders"`
}
