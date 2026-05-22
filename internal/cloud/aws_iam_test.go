package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIAMClient struct {
	out *iam.GetRoleOutput
	err error
}

func (m *mockIAMClient) GetRole(_ context.Context, _ *iam.GetRoleInput, _ ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	return m.out, m.err
}

func TestFetchIAMRole_Success(t *testing.T) {
	client := &mockIAMClient{
		out: &iam.GetRoleOutput{
			Role: &types.Role{
				RoleName:           aws.String("my-role"),
				Arn:                aws.String("arn:aws:iam::123456789012:role/my-role"),
				Path:               aws.String("/"),
				Description:        aws.String("test role"),
				MaxSessionDuration: aws.Int32(3600),
			},
		},
	}

	attrs, err := fetchIAMRole(context.Background(), client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["name"] != "my-role" {
		t.Errorf("expected name=my-role, got %q", attrs["name"])
	}
	if attrs["max_session_duration"] != "3600" {
		t.Errorf("expected max_session_duration=3600, got %q", attrs["max_session_duration"])
	}
	if attrs["description"] != "test role" {
		t.Errorf("expected description='test role', got %q", attrs["description"])
	}
}

func TestFetchIAMRole_NotFound(t *testing.T) {
	client := &mockIAMClient{
		err: &notFoundError{},
	}

	_, err := fetchIAMRole(context.Background(), client, "missing-role")
	var nf ErrResourceNotFound
	if !errors.As(err, &nf) {
		t.Fatalf("expected ErrResourceNotFound, got %v", err)
	}
	if nf.ResourceID != "missing-role" {
		t.Errorf("expected ResourceID=missing-role, got %q", nf.ResourceID)
	}
}

func TestFetchIAMRole_APIError(t *testing.T) {
	client := &mockIAMClient{
		err: errors.New("api failure"),
	}

	_, err := fetchIAMRole(context.Background(), client, "my-role")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchIAMRole_NoDescription(t *testing.T) {
	client := &mockIAMClient{
		out: &iam.GetRoleOutput{
			Role: &types.Role{
				RoleName: aws.String("bare-role"),
				Arn:      aws.String("arn:aws:iam::000000000000:role/bare-role"),
				Path:     aws.String("/"),
			},
		},
	}

	attrs, err := fetchIAMRole(context.Background(), client, "bare-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := attrs["description"]; ok {
		t.Error("description should be absent when not set")
	}
}
