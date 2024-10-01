package filemanager

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	domain "templify/pkg/domain/model"
	"testing"

	_ "embed"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

//go:embed example.pdf
var example_pdf []byte

func Test_listBucketContent(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManager(t)

	t.Log("Listing bucket content")

	// list bucket content
	objects, err := fm.ListFiles(fm.config.BucketName)
	if err != nil {
		t.Errorf("Error listing bucket content: %v", err)
	}
	t.Log("Bucket content:")
	if len(objects) > 0 {
		for _, object := range objects {
			t.Log(*object.Key)
		}
	} else {
		t.Log("No objects found")
	}
	t.FailNow()
}

func Test_listBuckets(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManager(t)
	t.Log("Listing buckets")

	// list buckets
	buckets, err := fm.ListBuckets()
	if err != nil {
		t.Errorf("Error listing buckets: %v", err)
	}
	t.Log("Buckets:")
	if len(buckets) > 0 {
		for _, bucket := range buckets {
			t.Log(*bucket.Name)
		}
	} else {
		t.Log("No buckets found")
	}
	t.FailNow()
}

func createNewTestFileManager(t *testing.T) *FileManager {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Logf("Error loading .env file: %v", err)
	}
	viper.AutomaticEnv()
	fileManagerConfig := &FileManagerConfig{
		BaseURL:     "http://localhost",
		Port:        viper.GetString("FILE_MANAGER_PORT"),
		BucketName:  viper.GetString("FILE_MANAGER_BUCKET_NAME"),
		Region:      viper.GetString("FILE_MANAGER_REGION"),
		AccessKeyID: viper.GetString("FILE_MANAGER_ACCESS_KEY_ID"),
		SecretKeyID: viper.GetString("FILE_MANAGER_SECRET_KEY_ID"),
	}
	fm := NewFileManagerService(fileManagerConfig, slog.Default())
	return fm
}

func Test_APIUploadFile(t *testing.T) {
	// create request to upload file to localhost using octet-stream
	r, err := http.NewRequest("POST", "http://localhost:8080/file/upload", bytes.NewReader(example_pdf))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}
	r.Header.Set("Content-Type", "multipart/octet-stream")

	// send request
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		t.Errorf("Error sending request: %v", err)
	}
	t.Logf("Response status: %s", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	t.Logf("Response body: %v", string(body))
	t.FailNow()
}

func Test_UploadFile(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManager(t)

	t.Log("Uploading file")

	err := fm.UploadFile(domain.FileUploadRequest{
		FileName:  "example",
		Extension: "pdf",
		FileBytes: example_pdf,
	})
	if err != nil {
		t.Errorf("Error uploading file: %v", err)
	}
	t.Log("File uploaded")
	t.FailNow()
}

func Test_DownloadFile(t *testing.T) {
	// create new file manager
	// create new file manager
	fm := createNewTestFileManager(t)

	t.Log("Downloading file")

	// download file
	file, err := fm.DownloadFile(domain.FileDownloadRequest{
		FileName:  "example",
		Extension: "pdf",
	})
	if err != nil {
		t.Errorf("Error downloading file: %v", err)
	}
	t.Logf("File content: %s", string(file))
	t.Logf("File size: %d", len(file))
	t.Log("File downloaded")
	t.FailNow()
}
