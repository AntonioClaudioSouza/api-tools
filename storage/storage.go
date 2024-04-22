package storage

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	bucketName string
)

// *** Init storage with bucket name default
func Init() {
	bucketName = os.Getenv("MINIOBUCKETNAME")
	if len(bucketName) == 0 {
		log.Fatal("MINIOBUCKETNAME not defined")
	}
}

// *** sendFile with expiration link if expireIn != nil
func sendFile(filePathSource, filePathDest string, subPath string, expireIn *time.Duration) (string, error) {

	urlUpload := ""
	var fileObject string
	var err error

	clientMinio, err := NewClient()
	if err != nil {
		return "", err
	}

	if len(subPath) > 0 {
		filePathDest = subPath + "/" + filePathDest
	}

	if expireIn != nil {
		urlUpload, err = clientMinio.UploadFileExpirationLink(bucketName, filePathDest, filePathSource, *expireIn)
		if err != nil {
			return "", err
		}
		return urlUpload, nil
	}

	fileObject, err = clientMinio.UploadFile(bucketName, filePathDest, filePathSource)
	if err != nil {
		return "", err
	}

	return fileObject, nil
}

// *** Send file by multipart with prefix and randon name
func SendFileByMultiPartPrefix(fileImport *multipart.FileHeader, c *gin.Context, subPath string, expireIn *time.Duration) (string, error) {

	now := time.Now()
	namefile := fmt.Sprintf("%d-%s", now.UnixNano(), fileImport.Filename)
	filePathFull := "/tmp/" + namefile
	if err := c.SaveUploadedFile(fileImport, filePathFull); err != nil {
		return "", err
	}
	defer os.Remove(filePathFull)
	return sendFile(filePathFull, namefile, subPath, expireIn)
}

// *** Send file by multipart
func SendFileByMultiPart(fileImport *multipart.FileHeader, c *gin.Context, subPath string, expireIn *time.Duration) (string, error) {
	filePathFull := "/tmp/" + fileImport.Filename
	if err := c.SaveUploadedFile(fileImport, filePathFull); err != nil {
		return "", err
	}
	defer os.Remove(filePathFull)
	return sendFile(filePathFull, fileImport.Filename, subPath, expireIn)
}

// *** Send file by local path
func SendFile(filePath string, subPath string, expireIn *time.Duration) (string, error) {
	return sendFile(filePath, filePath, subPath, expireIn)
}

// *** Get file bytes from storage
func GetFileBytes(fileName string) ([]byte, error) {

	fileName, err := GetFile(fileName, "")
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fileName)
}

// *** Get file from storage
func GetFile(fileName string, savePath string) (string, error) {

	// *** Create client
	clientMinio, err := NewClient()
	if err != nil {
		return "", err
	}

	if len(savePath) == 0 {
		savePath = "/tmp/"
	}

	fileNameDst := savePath + "/" + fileName

	//*** Check if file exists, if yes remove
	if _, err := os.Stat(fileNameDst); err == nil {
		if err := os.Remove(fileNameDst); err != nil {
			log.Println("Error remove file: ", err.Error())
			return "", err
		}
	}

	// *** Download and save file
	if err := clientMinio.DownloadFile(bucketName, fileName, fileNameDst); err != nil {
		return "", err
	}

	return fileNameDst, nil
}

// *** Delete file from storage
func DeleteFile(fileName string) error {

	// *** Create client
	clientMinio, err := NewClient()
	if err != nil {
		return err
	}

	// *** Delete file
	if err := clientMinio.RemoveFile(bucketName, fileName); err != nil {
		return err
	}
	return nil
}

func GetLinkDownload(fileName string, expireIn *time.Duration) (map[string]interface{}, error) {

	// *** Create client
	clientMinio, err := NewClient()
	if err != nil {
		return nil, err
	}

	// *** Get link download
	if len(fileName) > 0 {
		if expireIn == nil {
			expireIn = new(time.Duration)
			*expireIn = 7 * time.Hour * 24 // 7 days max
		}

		link, err := clientMinio.GetLinkDownload(bucketName, fileName, *expireIn)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"link":     link,
			"expireIn": expireIn.String(),
		}, nil
	}

	return nil, fmt.Errorf("file name is empty")
}
