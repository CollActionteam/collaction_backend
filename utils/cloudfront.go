package utils

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func CreateCFClient() *cloudfront.CloudFront {
	sess := session.Must(session.NewSession())
	return cloudfront.New(sess)
}

func InvalidateCache(distributionId string, pattern string) error {
	svc := CreateCFClient()
	_, err := svc.CreateInvalidation(&cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionId),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(
				fmt.Sprintf("invalidate p=\"%s\" t=\"%s\"", pattern, time.Now().Format("2006-01-02 15:04:05"))),
			Paths: &cloudfront.Paths{
				Quantity: aws.Int64(1),
				Items: []*string{
					aws.String(pattern),
				},
			},
		},
	})
	return err
}
