package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// mockS3Client implements s3Client for tests.
type mockS3Client struct {
	locationOut  *s3.GetBucketLocationOutput
	locationErr  error
	versioningOut *s3.GetBucketVersioningOutput
	versioningErr error
	taggingOut   *s3.GetBucketTaggingOutput
	taggingErr   error
}

func (m *mockS3Client) GetBucketLocation(_ context.Context, _ *s3.GetBucketLocationInput, _ ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return m.locationOut, m.locationErr
}
func (m *mockS3Client) GetBucketTagging(_ context.Context, _ *s3.GetBucketTaggingInput, _ ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return m.taggingOut, m.taggingErr
}
func (m *mockS3Client) GetBucketVersioning(_ context.Context, _ *s3.GetBucketVersioningInput, _ ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error) {
	return m.versioningOut, m.versioningErr
}

func TestFetchS3Bucket_Success(t *testing.T) {
	client := &mockS3Client{
		locationOut:  &s3.GetBucketLocationOutput{LocationConstraint: "eu-west-1"},
		versioningOut: &s3.GetBucketVersioningOutput{Status: types.BucketVersioningStatusEnabled},
		taggingOut: &s3.GetBucketTaggingOutput{
			TagSet: []types.Tag{{Key: aws.String("env"), Value: aws.String("prod")}},
		},
	}
	attrs, err := fetchS3Bucket(context.Background(), client, "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["region"] != "eu-west-1" {
		t.Errorf("region = %q, want eu-west-1", attrs["region"])
	}
	if attrs["versioning"] != "Enabled" {
		t.Errorf("versioning = %q, want Enabled", attrs["versioning"])
	}
	if attrs["tag:env"] != "prod" {
		t.Errorf("tag:env = %q, want prod", attrs["tag:env"])
	}
}

func TestFetchS3Bucket_DefaultRegion(t *testing.T) {
	client := &mockS3Client{
		locationOut:  &s3.GetBucketLocationOutput{}, // empty == us-east-1
		versioningOut: &s3.GetBucketVersioningOutput{},
		taggingErr:   errors.New("NoSuchTagSet"),
	}
	attrs, err := fetchS3Bucket(context.Background(), client, "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["region"] != "us-east-1" {
		t.Errorf("region = %q, want us-east-1", attrs["region"])
	}
}

func TestFetchS3Bucket_LocationError(t *testing.T) {
	client := &mockS3Client{
		locationErr: errors.New("NoSuchBucket"),
	}
	_, err := fetchS3Bucket(context.Background(), client, "missing-bucket")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchS3Bucket_VersioningError(t *testing.T) {
	client := &mockS3Client{
		locationOut:  &s3.GetBucketLocationOutput{LocationConstraint: "us-west-2"},
		versioningErr: errors.New("AccessDenied"),
	}
	_, err := fetchS3Bucket(context.Background(), client, "my-bucket")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
