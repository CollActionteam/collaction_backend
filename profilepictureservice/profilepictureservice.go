package profilepictureservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetUploadUrl(ext string, userID string) (string, error) {
	var (
		bucket  = "profile-url-picture-uploadbucket-5668889"
		filekey = userID
		region  = "eu-central-1"
	)

	mime, err := GetMime(ext)
	if err != nil {
		return "", err
	}

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return "", err
	}

	// Create S3 service client
	svc := s3.New(sess)
	reqs, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filekey),
		ContentType: aws.String(mime),
	})

	str, err := reqs.Presign(15 * time.Minute)

	// log.Println("The URL is:", str, " err:", err)
	if err != nil {
		return "", err
	}
	return str, nil
}

func GetMime(ext string) (string, error) {
	ext = strings.ToLower(ext)
	if ext == "png" {
		return "image/png", nil
	} else if ext == "jpg" {
		return "image/jpg", nil
	} else if ext == "jpeg" {
		return "image/jpeg", nil
	}
	return "", fmt.Errorf("file format not supported")

}
