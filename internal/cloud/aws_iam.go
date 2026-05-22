package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// iamClient defines the IAM operations used by driftwatch.
type iamClient interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
}

// fetchIAMRole retrieves attributes of an IAM role by name.
func fetchIAMRole(ctx context.Context, client iamClient, resourceID string) (map[string]string, error) {
	out, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(resourceID),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrResourceNotFound{ResourceID: resourceID}
		}
		return nil, fmt.Errorf("iam get role %q: %w", resourceID, err)
	}

	return roleAttributes(out), nil
}

func roleAttributes(out *iam.GetRoleOutput) map[string]string {
	attrs := map[string]string{}
	if out.Role == nil {
		return attrs
	}
	r := out.Role
	if r.RoleName != nil {
		attrs["name"] = aws.ToString(r.RoleName)
	}
	if r.Arn != nil {
		attrs["arn"] = aws.ToString(r.Arn)
	}
	if r.Path != nil {
		attrs["path"] = aws.ToString(r.Path)
	}
	if r.Description != nil {
		attrs["description"] = aws.ToString(r.Description)
	}
	if r.MaxSessionDuration != nil {
		attrs["max_session_duration"] = fmt.Sprintf("%d", aws.ToInt32(r.MaxSessionDuration))
	}
	return attrs
}
