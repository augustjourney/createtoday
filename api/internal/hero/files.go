package hero

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/logger"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFileToS3(bucket string, fileName string, fileBytes io.ReadSeeker, config *config.Config) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.S3Region),
		Endpoint:         aws.String(config.S3Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewStaticCredentials(
			config.S3AccessKeyId,
			config.S3SecretAccessKey,
			"",
		),
	})
	if err != nil {
		logger.Log.Error("could not create session in aws sdk", "err", err)
		return "", err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   fileBytes,
	})

	if err != nil {
		logger.Log.Error("could not upload file to s3", "err", err, "fileName", fileName, "bucket", bucket)
		return "", err
	}

	fileUrlOnCdn := config.CdnUrl + "/" + bucket + "/" + fileName

	return fileUrlOnCdn, nil
}

func RemoveLocalFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		logger.Log.Error(err.Error(), "filePath", path)
		return err
	}
	return nil
}

func GetExtensionFromFileName(fileName string) string {
	ext := strings.SplitAfter(fileName, ".")
	if len(ext) < 2 {
		logger.Log.Error("file does not contain extension", "file", fileName)
		return ""
	}
	return ext[len(ext)-1]
}

func GetMediaTypeFromMime(mime string) string {
	parts := strings.SplitN(mime, "/", 2)
	if len(parts) != 2 {
		return "unknown-type"
	}
	return parts[0]
}

func MakeFileHashName(text string, ext string) string {
	hash := md5.Sum([]byte(text))
	hashString := hex.EncodeToString(hash[:])
	return hashString + "." + ext
}
