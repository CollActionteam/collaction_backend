package main

import (
	"bytes"
	"errors"
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
	"github.com/nfnt/resize"
)

const (
	moderationConfidenceThreshold = 0.9
	maxFileSize                   = 5 * 1024 * 1024
	minWidth                      = 250
	maxWidth                      = 1024
	preferedSize                  = 300
)

func checkContent(clientRekognition *rekognition.Rekognition, bucketName *string, key *string) error {
	contains := func(slice []string, item string) bool {
		for _, candidate := range slice {
			if candidate == item {
				return true
			}
		}
		return false
	}

	mlRes, err := clientRekognition.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: bucketName,
				Name:   key,
			},
		},
	})
	if err != nil {
		return err
	}
	// Refer to https://docs.aws.amazon.com/rekognition/latest/dg/moderation.html#moderation-api
	unacceptableLabels := []string{"Explicit Nudity", "Violence", "Visually Disturbing", "Hate Symbols"}
	for _, label := range mlRes.ModerationLabels {
		if *label.Confidence >= moderationConfidenceThreshold {
			if contains(unacceptableLabels, *label.Name) || contains(unacceptableLabels, *label.ParentName) {
				reason := fmt.Sprintf("Rejected file %s because is might contain %s (%f%% confidence)!\n", *key, *label.Name, *label.Confidence)
				return errors.New(reason)
			}
		}
	}
	return nil
}

func processImage(imgBytes []byte) ([]byte, error) {
	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(imgBytes))

	if imgCfg.Width != imgCfg.Height {
		return nil, errors.New("image does not have an aspect ratio of 1.00")
	}

	if imgCfg.Width < minWidth || imgCfg.Width > maxWidth {
		return nil, fmt.Errorf("image is not between %d and %d pixels wide", minWidth, maxWidth)
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	img = resize.Thumbnail(preferedSize, preferedSize, img, resize.Lanczos3)

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

		// Analyze content
		err = checkContent(clientRekognition, &inputBucketName, &key)
		if err != nil {
			log.Println(err.Error())
			return
		}

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
