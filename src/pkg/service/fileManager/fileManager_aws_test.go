package filemanager

import (
	"log/slog"
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
	fm := createNewTestFileManagerAWS(t)

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

func Test_listBuckets(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerAWS(t)
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

func createNewTestFileManagerAWS(t *testing.T) *FileManagerAWS {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		t.Logf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()
	fileManagerConfig := &FileManagerConfig{
		BucketName:  viper.GetString("FILE_MANAGER_BUCKET_NAME"),
		Region:      viper.GetString("FILE_MANAGER_REGION"),
		AccessKeyID: viper.GetString("FILE_MANAGER_ACCESS_KEY_ID"),
		SecretKeyID: viper.GetString("FILE_MANAGER_SECRET_KEY_ID"),
	}

	return NewFileManagerAWSService(fileManagerConfig, slog.Default())
}

func Test_UploadFile(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerAWS(t)

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
	fm := createNewTestFileManagerAWS(t)

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

func Test_AWSGetObjectPresignedURL(t *testing.T) {
	// create new file manager
	fm := createNewTestFileManagerAWS(t)

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

func Test_AWSPostObjectPresignedURL(t *testing.T) {
	fm := createNewTestFileManagerAWS(t)

	t.Log("Getting presigned URL")

	// get presigned URL
	url, err := fm.GetFileUploadURL("example.pdf")
	if err != nil {
		t.Errorf("Error getting presigned URL: %v", err)
	}
	t.Logf("Presigned URL: %s", url.UploadURL)
	t.Logf("Presigned URL values: %v", url.Values)

	t.Log("Presigned URL retrieved")
	t.FailNow()
}
