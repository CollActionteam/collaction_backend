package main

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(e events.S3Event) {
	outputBucketName := os.Getenv("OUTPUT_BUCKET_NAME")
	svc := s3.New(session.Must(session.NewSession()))

	process_object := func(inputBucketName string, key string) {
		var err error

		// Download image
		dlRes, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(inputBucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			panic(err.Error())
		}
		var b bytes.Buffer
		_, err = b.ReadFrom(dlRes.Body)
		if err != nil {
			panic(err.Error())
		}
		dlRes.Body.Close()

		// Process image
		img, _, err := image.Decode(bytes.NewReader(b.Bytes()))
		if err != nil {
			panic(err.Error())
		}
		var f bytes.Buffer
		err = png.Encode(&f, img)
		if err != nil {
			panic(err.Error())
		}

		// Upload image
		_, err = svc.PutObject(&s3.PutObjectInput{
			Body:        bytes.NewReader(f.Bytes()),
			Bucket:      aws.String(outputBucketName),
			Key:         aws.String(key),
			ContentType: aws.String("image/png"),
			ACL:         aws.String("public-read"),
		})
		if err != nil {
			panic(err.Error())
		}

		// Delete user uploaded image
		_, err = svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(inputBucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			panic(err.Error())
		}
	}

	for _, r := range e.Records {
		bucketName := r.S3.Bucket.Name
		// TODO maybe check if the bucket is the correct one
		// Should it check that the target and source bucket are not the same?
		key := r.S3.Object.Key
		if strings.HasSuffix(key, ".png") {
			process_object(bucketName, key)
		}
	}
}

func main() {
	lambda.Start(handler)
}
