package profilepictureservice

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetUploadUrl(ext string, userID string) (string, error) {

	var (
		bucket  = os.Getenv("BUCKET")
		filekey = userID + "." + ext
	)

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	// Create S3 service client
	svc := s3.New(session.Must(session.NewSession()))
	reqs, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filekey),
	})

	str, err := reqs.Presign(15 * time.Minute)

	if err != nil {
		return "", err
	}
	return str, nil
}
