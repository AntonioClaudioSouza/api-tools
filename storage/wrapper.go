package storage

import (
	"context"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioWrapper struct {
	*minio.Client
}

func NewClient() (*MinioWrapper, error) {

	urlConnection := os.Getenv("MINIOURL")
	accessKey := os.Getenv("MINIOACCESSKEY")
	secretKey := os.Getenv("MINIOSECRETKEY")

	// *** Initialize minio client object.
	minioClient, err := minio.New(urlConnection, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioWrapper{minioClient}, nil
}

// *** Create a new bucket is not exists.
func (m *MinioWrapper) CreateNewBucketIfNotExists(bucketName string) error {

	exists, err := m.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return m.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

// *** Upload a new object no external link download.
func (m *MinioWrapper) UploadFile(bucketName, objectName, filePath string) (string, error) {

	var err error
	err = m.CreateNewBucketIfNotExists(bucketName)
	if err != nil {
		return "", err
	}

	// Upload the zip file with FPutObject
	_, err = m.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

// *** Upload a new object and create external link download.
func (m *MinioWrapper) UploadFileExpirationLink(bucketName, objectName, filePath string, expireIn time.Duration) (string, error) {

	var err error
	if err = m.CreateNewBucketIfNotExists(bucketName); err != nil {
		return "", err
	}

	// Upload the zip file with FPutObject
	_, err = m.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	// Generate a presigned url which expires in expireIn variable.
	presignedURL, err := m.PresignedGetObject(context.Background(), bucketName, objectName, expireIn, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// *** Remove object from bucket.
func (m *MinioWrapper) RemoveFile(bucketName, objectName string) error {

	if err := m.CreateNewBucketIfNotExists(bucketName); err != nil {
		return err
	}

	//*** Check if object exists
	_, err := m.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	return m.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
}

// *** Download object from bucket.
// *** Only no public objects
func (m *MinioWrapper) DownloadFile(bucketName, objectName, filePath string) error {

	if err := m.CreateNewBucketIfNotExists(bucketName); err != nil {
		return err
	}

	//*** Check if object exists
	_, err := m.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return err
	}

	return m.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
}

func (m *MinioWrapper) GetLinkDownload(bucketName, objectName string, expireIn time.Duration) (string, error) {

	if err := m.CreateNewBucketIfNotExists(bucketName); err != nil {
		return "", err
	}

	//*** Check if object exists
	_, err := m.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", err
	}

	// Generate a presigned url which expires in expireIn variable.
	presignedURL, err := m.PresignedGetObject(context.Background(), bucketName, objectName, expireIn, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
