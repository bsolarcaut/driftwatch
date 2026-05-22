package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// fetchIAMRoleViaAWS wires up a real IAM client and delegates to fetchIAMRole.
// It is called by AWSFetcher when the resource type is "aws_iam_role".
func fetchIAMRoleViaAWS(ctx context.Context, resourceID string) (map[string]string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config for iam: %w", err)
	}
	client := iam.NewFromConfig(cfg)
	return fetchIAMRole(ctx, client, resourceID)
}
