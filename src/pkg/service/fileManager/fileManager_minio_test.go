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

func Test_listBucketContentMinio(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerMinio(t)

	t.Log("Listing bucket content")

	// list bucket content
	objects, err := fm.ListFiles(fm.config.BucketName)
	if err != nil {
		t.Errorf("Error listing bucket content: %v", err)
	}
	t.Log("Bucket content:")
	if len(objects) > 0 {
		for _, object := range objects {
			t.Log(object)
		}
	} else {
		t.Log("No objects found")
	}
	t.FailNow()
}

func Test_listBucketsMinio(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerMinio(t)
	t.Log("Listing buckets")

	// list buckets
	buckets, err := fm.ListBuckets()
	if err != nil {
		t.Errorf("Error listing buckets: %v", err)
	}
	t.Log("Buckets:")
	if len(buckets) > 0 {
		for _, bucket := range buckets {
			t.Log(bucket)
		}
	} else {
		t.Log("No buckets found")
	}
	t.FailNow()
}

func createNewTestFileManagerMinio(t *testing.T) *FileManagerMinio {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Logf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()
	fileManagerConfig := &FileManagerConfig{
		BaseURL:     "localhost",
		Port:        viper.GetString("FILE_MANAGER_PORT"),
		BucketName:  viper.GetString("FILE_MANAGER_BUCKET_NAME"),
		Region:      viper.GetString("FILE_MANAGER_REGION"),
		AccessKeyID: viper.GetString("FILE_MANAGER_ACCESS_KEY_ID"),
		SecretKeyID: viper.GetString("FILE_MANAGER_SECRET_KEY_ID"),
	}

	return NewFileManagerMinioService(fileManagerConfig, slog.Default())
}

func Test_UploadFileMinio(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerMinio(t)

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

func Test_DownloadFileMinio(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerMinio(t)

	t.Log("Downloading file")

	// download file
	file, err := fm.DownloadFile(domain.FileDownloadRequest{
		FileName:  "example",
		Extension: "pdf",
	})
	if err != nil {
		t.Errorf("Error downloading file: %v", err)
	}
	t.Logf("File size: %d", len(file))
	t.Log("File downloaded")
	t.FailNow()
}

func Test_GetObjectPresignedURLMinio(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerMinio(t)

	t.Log("Getting presigned URL")

	// get presigned URL
	url, err := fm.GetFileDownloadURL("example.pdf")
	if err != nil {
		t.Errorf("Error getting presigned URL: %v", err)
	}
	t.Logf("Bucket name: %s", fm.config.BucketName)
	t.Logf("AWS credentials: ")
	t.Logf("Access key ID: %s", fm.config.AccessKeyID)
	t.Logf("Secret access key: %s", fm.config.SecretKeyID)
	t.Logf("Region: %s", fm.config.Region)
	t.Logf("====================")
	t.Logf("Presigned URL: %s", url)
	t.FailNow()
}

func Test_PostObjectPresignedURLMinio(t *testing.T) {
	fm := createNewTestFileManagerMinio(t)

	t.Log("Getting presigned URL")

	// get presigned URL
	url, err := fm.GetFileUploadURL("example.pdf")
	if err != nil {
		t.Errorf("Error getting presigned URL: %v", err)
	}
	t.Logf("Presigned URL retrieved: %s", url.UploadURL)

	t.Log("Presigned URL retrieved")

	// new request using example.pdf with presigned URL
	req, err := http.NewRequest("POST", url.UploadURL, bytes.NewReader(example_pdf))
	if err != nil {
		t.Errorf("Error creating new request: %v", err)
	}

	// do request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error doing request: %v", err)
	}
	defer resp.Body.Close()

	t.Logf("Response status: %s", resp.Status)
	t.Logf("Response content length: %d", resp.ContentLength)
	reqBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}
	t.Logf("Response: %s", reqBody)

	t.FailNow()
}
