package mjmlservice

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed template_test.mjml
var testTemplate string

func TestExtractPlaceholders(t *testing.T) {
	// prepare
	expected := []string{"Message", "Banana", "Apple", "Juice"}

	// test
	placeholders := extractPlaceholders(testTemplate)

	// assert
	assert.Equal(t, expected, placeholders, "Expected %v is not the same as we got from the extractPlaceholders function: %v", expected, placeholders)
}

// func TestRenderMJML(t *testing.T) {
// 	//prepare
// 	expected:= ""
// }
