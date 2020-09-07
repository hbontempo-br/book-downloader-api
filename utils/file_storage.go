package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v6"
	"go.uber.org/zap"
)

const (
	ErrClientCreation    = FileStorageErr("error on file storage client creation")
	ErrFileSave          = FileStorageErr("error on saving file to file storage")
	ErrFileRetrieval     = FileStorageErr("error on retrieving file from file storage")
	ErrFileLinkRetrieval = FileStorageErr("error on retrieving file link from file storage")
)

type FileStorageErr string

func (e FileStorageErr) Error() string {
	return string(e)
}

type FileStorage interface {
	Save(reader io.Reader, bucket, location string) error
	Get(bucket, location string) (io.Reader, error)
}

func NewMinioFileStorage(endpoint, accessKey, secretKey string, ssl bool) (MinioFileStorage, error) {
	mfs := MinioFileStorage{endpoint: endpoint, accessKey: accessKey, secretKey: secretKey, ssl: ssl}
	err := mfs.createClient()
	return mfs, err
}

type MinioFileStorage struct {
	endpoint  string
	accessKey string
	secretKey string
	ssl       bool
	client    *minio.Client
}

func (mfs *MinioFileStorage) createClient() error {
	client, err := minio.NewV4(mfs.endpoint, mfs.accessKey, mfs.secretKey, mfs.ssl)
	if err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.createClient", "errors", err)
		return ErrClientCreation
	}
	mfs.client = client
	return nil
}

func (mfs *MinioFileStorage) Save(reader io.Reader, bucket, location string) error {
	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Save", "errors", err)
		return ErrFileSave
	}
	r := bytes.NewReader(buffer.Bytes())
	mt, err := mimetype.DetectReader(r)
	if err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Save", "errors", err)
		return ErrFileSave
	}
	mtStr := mt.String()

	if _, err := r.Seek(0, 0); err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Save", "errors", err)
		return ErrFileSave
	}

	if _, err := mfs.client.PutObject(bucket, location, r, -1, minio.PutObjectOptions{ContentType: mtStr}); err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Save", "errors", err)
		return ErrFileSave
	}

	return nil
}

func (mfs *MinioFileStorage) Get(bucket, location string) (io.Reader, error) {
	reader, err := mfs.client.GetObject(bucket, location, minio.GetObjectOptions{})
	if err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Get", "errors", err)
		return nil, ErrFileRetrieval
	}
	return reader, nil
}

func (mfs *MinioFileStorage) GetLink(filename, bucket, location string, expiry time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%v"`, filename))
	presignedURL, err := mfs.client.PresignedGetObject(bucket, location, expiry, reqParams)
	if err != nil {
		zap.S().Errorw("Generic error on MinioFileStorage.Get", "errors", err)
		return nil, ErrFileLinkRetrieval
	}

	return presignedURL, nil
}
