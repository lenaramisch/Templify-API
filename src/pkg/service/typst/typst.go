package typst

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"templify/pkg/domain"
)

type TypstConfig struct {
}

type TypstService struct {
	config TypstConfig
}

func NewTypstService(config TypstConfig) *TypstService {
	return &TypstService{
		config: config,
	}
}

// TODO make general write String to file method that uses path instead of filename
func writeStringToFile(filledTemplStr string, fileName string) (string, error) {
	dir := "/tmp"
	typstFileName := filepath.Join(dir, fileName+".typ")
	// Create the file
	f, err := os.Create(typstFileName)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", typstFileName, err)
		return "", err
	}

	// Write the string to the file
	l, err := f.WriteString(filledTemplStr)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", typstFileName, err)
		f.Close()
		return "", err
	}
	fmt.Printf("%d bytes written successfully to %s\n", l, typstFileName)

	// Close the file
	err = f.Close()
	if err != nil {
		fmt.Printf("Error closing file %s: %v\n", typstFileName, err)
		return "", err
	}

	return typstFileName, nil
}

func (t *TypstService) RenderTypst(typstString string) ([]byte, error) {
	randomName := "djaijsdh" // TODO use UUID here
	writeStringToFile(typstString, randomName)
	cmd := exec.Command("typst", "compile", randomName)
	cmd.Dir = "/tmp"
	err := cmd.Run()
	if err != nil {
		slog.With(
			"Error", err.Error(),
		).Debug("Error executing typst command")
		return nil, err
	}
	// open the generated PDF file
	filePath := fmt.Sprintf("/tmp/%s.pdf", randomName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.With(
			"filePath", filePath,
		).Debug("Error reading PDF file")
	}
	return bytes, nil
}
