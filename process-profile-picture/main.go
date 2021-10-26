package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"
)

const moderationConfidenceThreshold = 0.9
const maxFileSize = 5 * 1024 * 1024

func contains(slice []string, item string) bool {
	for _, candidate := range slice {
		if candidate == item {
			return true
		}
	}
	return false
}

func checkIsContentOK(clientRekognition *rekognition.Rekognition, bucketName *string, key *string) (bool, *string, error) {
	mlRes, err := clientRekognition.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: bucketName,
				Name:   key,
			},
		},
	})
	if err != nil {
		return false, nil, err
	}
	// Refer to https://docs.aws.amazon.com/rekognition/latest/dg/moderation.html#moderation-api
	unacceptableLabels := []string{"Explicit Nudity", "Violence", "Visually Disturbing", "Hate Symbols"}
	for _, label := range mlRes.ModerationLabels {
		if *label.Confidence >= moderationConfidenceThreshold {
			if contains(unacceptableLabels, *label.Name) || contains(unacceptableLabels, *label.ParentName) {
				reason := fmt.Sprintf("%s (%f%% confidence)", *label.Name, *label.Confidence)
				return false, &reason, nil
			}
		}
	}
	return true, nil, nil
}

func processImage(imgBytes []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	var f bytes.Buffer
	err = png.Encode(&f, img)
	if err != nil {
		return nil, err
	}
	return f.Bytes(), nil
}

func handler(e events.S3Event) {
	outputBucketName := os.Getenv("OUTPUT_BUCKET_NAME")
	sess := session.Must(session.NewSession())
	clientS3 := s3.New(sess)
	clientRekognition := rekognition.New(sess)

	process_object := func(inputBucketName string, key string) {
		var err error

		defer func() {
			// Delete user uploaded image
			_, err = clientS3.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(inputBucketName),
				Key:    aws.String(key),
			})
			if err != nil {
				log.Println(err.Error())
			}
		}()

		// Download image
		dlRes, err := clientS3.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(inputBucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
		if *dlRes.ContentLength > int64(maxFileSize) {
			log.Printf("Size of file %s exceedes %d bytes!\n", key, maxFileSize)
			return
		}
		var b bytes.Buffer
		_, err = b.ReadFrom(dlRes.Body)
		defer dlRes.Body.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Analyze content
		isContentOK, reason, err := checkIsContentOK(clientRekognition, &inputBucketName, &key)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if !isContentOK {
			log.Printf("Rejected file %s because %s!\n", key, *reason)
			return
		}

		// Process image
		imgBytes, err := processImage(b.Bytes())
		if err != nil {
			log.Println(err.Error())
			return
		}

		// Upload image
		_, err = clientS3.PutObject(&s3.PutObjectInput{
			Body:        bytes.NewReader(imgBytes),
			Bucket:      aws.String(outputBucketName),
			Key:         aws.String(key),
			ContentType: aws.String("image/png"),
			ACL:         aws.String("public-read"),
		})
		if err != nil {
			log.Println(err.Error())
		}
	}

	for _, r := range e.Records {
		bucketName := r.S3.Bucket.Name
		// TODO Maybe check if the bucket is the correct one
		// TODO Should it check that the target and source bucket are not the same?
		key := r.S3.Object.Key
		if strings.HasSuffix(key, ".png") {
			process_object(bucketName, key)
		}
	}
}

func main() {
	lambda.Start(handler)
}
