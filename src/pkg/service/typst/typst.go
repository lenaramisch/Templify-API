package typst

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type TypstConfig struct {
}

type TypstService struct {
	config *TypstConfig
	log    *slog.Logger
}

func NewTypstService(config *TypstConfig, log *slog.Logger) *TypstService {
	return &TypstService{
		config: config,
	}
}

func writeStringToFile(filledTemplStr string, filePath string) (string, error) {
	// Create the file
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return "", err
	}

	// Write the string to the file
	l, err := f.WriteString(filledTemplStr)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", filePath, err)
		f.Close()
		return "", err
	}
	fmt.Printf("%d bytes written successfully to %s\n", l, filePath)

	// Close the file
	err = f.Close()
	if err != nil {
		fmt.Printf("Error closing file %s: %v\n", filePath, err)
		return "", err
	}

	return filePath, nil
}

func (t *TypstService) RenderTypst(typstString string) ([]byte, error) {
	randomName := uuid.New().String()
	filePath := fmt.Sprintf("/tmp/%s.typ", randomName)
	writeStringToFile(typstString, filePath)
	defer os.Remove(filePath)
	completeTypstFileName := fmt.Sprintf("%s.typ", randomName)
	cmd := exec.Command("typst", "compile", completeTypstFileName)
	cmd.Dir = "/tmp"
	err := cmd.Run()
	if err != nil {
		t.log.With(
			"Error", err.Error(),
		).Debug("Error executing typst command")
		return nil, err
	}
	defer os.Remove(fmt.Sprintf("/tmp/%s.pdf", randomName))
	// open the generated PDF file
	filePath = fmt.Sprintf("/tmp/%s.pdf", randomName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		t.log.With(
			"filePath", filePath,
		).Debug("Error reading PDF file")
	}

	return bytes, nil
}
