package filemanager

import (
	"log/slog"
	"testing"
)

func Test_listBucket(t *testing.T) {
	// create new file manager
	_ = NewFileManagerService(&FileManagerConfig{
		BaseURL:     "http://localhost",
		Port:        "9000",
		BucketName:  "templify-static-files",
		Region:      "eu-central-1",
		AccessKeyID: "JJVRguTZGE6pBwxNh6gO",
		SecretKeyID: "gAqSlxClhXIWoLSkZeqdgKjW15lfO2KsUzJIRCXj",
	}, slog.Default())

	t.Log("Listing buckets")

	// // list buckets
	// buckets, err := fileManager.ListBuckets()
	// if err != nil {
	// 	t.Errorf("Error listing buckets: %v", err)
	// }
	// t.Log("Buckets:")
	// if len(buckets) > 0 {
	// 	for _, bucket := range buckets {
	// 		t.Log(bucket.Name)
	// 	}
	// } else {
	// 	t.Log("No buckets found")
	// }
}
