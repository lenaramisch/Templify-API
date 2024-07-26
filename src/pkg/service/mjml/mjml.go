package mjmlservice

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type MJMLConfig struct {
	Host string
	Port int
}

type MJMLService struct {
	config *MJMLConfig
}

func NewMJMLService(config *MJMLConfig) *MJMLService {
	return &MJMLService{
		config: config,
	}
}

func (m *MJMLService) RenderMJML(MJMLString string) (string, error) {
	// Create a new POST request
	url := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(MJMLString)))
	if err != nil {
		slog.Debug("Error creating request")
		return "", err
	}

	r.Header.Set("Content-Type", "text/plain")

	//Perform the request
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		slog.With("err", err.Error()).Debug("Error sending request")
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Println("Error response from MJML Service:", string(body))
		return "", fmt.Errorf(fmt.Sprintf("MJML Service returned status code %d", resp.StatusCode))
	}

	htmlString := string(body)
	return htmlString, nil
}
