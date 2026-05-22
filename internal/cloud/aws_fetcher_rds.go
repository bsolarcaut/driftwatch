package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

const resourceTypeRDS = "aws_db_instance"

// fetchRDSInstanceViaAWS loads AWS credentials and fetches an RDS DB instance.
func fetchRDSInstanceViaAWS(ctx context.Context, resourceID string) (map[string]string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	client := rds.NewFromConfig(cfg)
	return fetchRDSInstance(ctx, client, resourceID)
}

// fetchForRDS dispatches RDS fetch calls from the AWSFetcher.
func (f *AWSFetcher) fetchForRDS(ctx context.Context, resourceID string) (map[string]string, error) {
	return fetchRDSInstanceViaAWS(ctx, resourceID)
}
