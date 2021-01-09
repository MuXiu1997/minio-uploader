package uploader

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

func (u Uploader) Upload(folder, filePath string) (string, error) {
	objectName, isTemp, err := u.buildObjectName(folder, &filePath)
	if isTemp {
		defer func() {
			_ = os.Remove(filePath)
		}()
	}

	if err != nil {
		return "", err
	}

	_, err = u.minioClient.FPutObject(context.Background(), u.bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return u.buildReturnUrl(objectName), nil
}

func (u Uploader) checkLocal(filePath *string) (filename string, isLocal bool, err error) {
	if strings.HasPrefix(*filePath, "http://") || strings.HasPrefix(*filePath, "https://") {
		fileUrl := *filePath
		filename, err = httpFilename(fileUrl)
		if err != nil {
			return "", false, err
		}
		tempFilePath, err := downloadFile(fileUrl)
		if err != nil {
			return "", false, err
		}
		*filePath = tempFilePath
		return filename, false, nil
	}
	return "", true, nil
}

func (u Uploader) buildObjectName(folder string, filePath *string) (objectName string, isTemp bool, err error) {
	filename, isLocal, err := u.checkLocal(filePath)
	if err != nil {
		return "", false, err
	}
	if isLocal {
		filename = filepath.Base(*filePath)
	}

	fileExt := filepath.Ext(filename)
	filename = filename[:len(filename)-len(fileExt)]

	if u.useUUID {
		filename = uuid.Must(uuid.NewRandom()).String()
	}

	return path.Join(folder, filename+fileExt), !isLocal, nil
}

func (u Uploader) buildReturnUrl(objectName string) string {
	return strings.Replace(u.returnUrlTemplate, "{objectName}", objectName, -1)
}

func httpFilename(url string) (string, error) {
	reg, _ := regexp.Compile(`http[s]?://.*/([^/]*)`)
	subMatch := reg.FindStringSubmatch(url)
	if len(subMatch) < 2 {
		return "", fmt.Errorf("can not get file name")
	}
	return subMatch[1], nil
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	file, err := ioutil.TempFile("", "minio-uploader-temp-")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(body)
	if err != nil {
		return "", err
	}
	return filepath.Abs(file.Name())
}
