package typst

import (
	"log/slog"
	"testing"
)

// test render method
func TestRenderTypst(T *testing.T) {
	typstService := NewTypstService(&TypstConfig{}, &slog.Logger{})

	// test render method
	PDFbytes, err := typstService.RenderTypst("= Client Name: {{.Name}}, Age: {{.Age}}")
	if err != nil {
		T.Errorf("Error rendering typst: %v", err)
	}
	if len(PDFbytes) == 0 {
		T.Errorf("Empty PDF file")
	}
}
