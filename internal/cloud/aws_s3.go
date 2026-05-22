package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// s3Client defines the subset of the S3 API used by fetchS3Bucket.
type s3Client interface {
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}

// fetchS3Bucket retrieves live attributes for an aws_s3_bucket resource.
func fetchS3Bucket(ctx context.Context, client s3Client, bucketName string) (map[string]string, error) {
	attrs := make(map[string]string)

	// Location / region
	locOut, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("get bucket location %q: %w", bucketName, err)
	}
	region := string(locOut.LocationConstraint)
	if region == "" {
		region = "us-east-1" // AWS returns empty string for us-east-1
	}
	attrs["region"] = region

	// Versioning
	verOut, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("get bucket versioning %q: %w", bucketName, err)
	}
	attrs["versioning"] = string(verOut.Status)

	// Tags
	tagOut, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		for _, tag := range tagOut.TagSet {
			if tag.Key != nil && tag.Value != nil {
				attrs["tag:"+*tag.Key] = *tag.Value
			}
		}
	}
	// Ignore tagging errors — bucket may simply have no tags.

	return attrs, nil
}
