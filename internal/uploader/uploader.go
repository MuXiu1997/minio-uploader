package uploader

import (
	"context"
	"fmt"
	"minio-uploader/internal/transformer"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

const (
	KeyBucket            = "bucket"
	KeyUseUUID           = "useUUID"
	KeyReturnUrlTemplate = "returnUrlTemplate"
)

type Uploader struct {
	minioClient       *minio.Client
	bucket            string
	useUUID           bool
	returnUrlTemplate string
}

func NewUploader(v *viper.Viper, minioClient *minio.Client) (*Uploader, error) {
	for _, key := range []string{KeyBucket, KeyUseUUID, KeyReturnUrlTemplate} {
		if !v.IsSet(key) {
			return nil, fmt.Errorf("%s is not set", key)
		}
	}
	return &Uploader{
		minioClient:       minioClient,
		bucket:            v.GetString(KeyBucket),
		useUUID:           v.GetBool(KeyUseUUID),
		returnUrlTemplate: v.GetString(KeyReturnUrlTemplate),
	}, nil
}

func (u Uploader) Upload(folder string, files []string) ([]string, error) {
	filesCount := len(files)
	successCount := 0

	results := make([]string, filesCount)
	resultChan := make(chan struct {
		Index     int
		ReturnUrl string
	})
	errChan := make(chan error)

	for i, file := range files {
		go func(i int, file string) {
			t := transformer.Factory(file)
			result, err := t.Transform()
			if err != nil {
				errChan <- err
				return
			}
			objectName := u.buildObjectName(folder, result)
			_, err = u.minioClient.PutObject(
				context.Background(),
				u.bucket, objectName,
				result.Buffer,
				-1,
				minio.PutObjectOptions{},
			)
			if err != nil {
				errChan <- err
				return
			}
			returnUrl := u.buildReturnUrl(objectName)
			resultChan <- struct {
				Index     int
				ReturnUrl string
			}{Index: i, ReturnUrl: returnUrl}
		}(i, file)
	}

	for {
		select {
		case err := <-errChan:
			return nil, err
		case result := <-resultChan:
			results[result.Index] = result.ReturnUrl
			successCount++
			if successCount == filesCount {
				return results, nil
			}
		}
	}
}

func (u Uploader) buildObjectName(folder string, result *transformer.Result) string {
	filename := result.Filename
	if u.useUUID {
		filename = uuid.Must(uuid.NewRandom()).String()
	}
	objectName := path.Join(folder, filename+result.FileExt)
	return objectName
}

func (u Uploader) buildReturnUrl(objectName string) string {
	return strings.Replace(u.returnUrlTemplate, "{objectName}", objectName, -1)
}
