package minioclient

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

const (
	KeyEndpoint  = "minio.endpoint"
	KeyAccessKey = "minio.accessKey"
	KeySecretKey = "minio.secretKey"
	KeyUseSSL    = "minio.useSSL"
)

func NewMinioClient(v *viper.Viper) (*minio.Client, error) {
	for _, key := range []string{KeyEndpoint, KeyAccessKey, KeySecretKey, KeyUseSSL} {
		if !v.IsSet(key) {
			return nil, fmt.Errorf("%s is not set", key)
		}
	}

	c, err := minio.New(v.GetString(KeyEndpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(v.GetString(KeyAccessKey), v.GetString(KeySecretKey), ""),
		Secure: v.GetBool(KeyUseSSL),
	})
	if err != nil {
		return nil, err
	}

	return c, nil
}
