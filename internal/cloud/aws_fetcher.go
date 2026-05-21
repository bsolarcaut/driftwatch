package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// AWSFetcher fetches live resource attributes from AWS.
type AWSFetcher struct {
	ec2Client *ec2.Client
}

// NewAWSFetcher creates an AWSFetcher using the default AWS credential chain.
func NewAWSFetcher(ctx context.Context, region string) (*AWSFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &AWSFetcher{ec2Client: ec2.NewFromConfig(cfg)}, nil
}

// Fetch returns the live attributes for the given resource type and ID.
// Currently supports: aws_instance.
func (f *AWSFetcher) Fetch(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	switch resourceType {
	case "aws_instance":
		return f.fetchEC2Instance(ctx, resourceID)
	default:
		return nil, fmt.Errorf("%w: type %q id %q", ErrResourceNotFound, resourceType, resourceID)
	}
}

func (f *AWSFetcher) fetchEC2Instance(ctx context.Context, instanceID string) (map[string]string, error) {
	out, err := f.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instance %q: %w", instanceID, err)
	}

	var instance *types.Instance
	for _, r := range out.Reservations {
		for i := range r.Instances {
			instance = &r.Instances[i]
			break
		}
		if instance != nil {
			break
		}
	}

	if instance == nil {
		return nil, fmt.Errorf("%w: aws_instance %q", ErrResourceNotFound, instanceID)
	}

	attrs := map[string]string{
		"instance_type": string(instance.InstanceType),
		"instance_state": string(instance.State.Name),
		"ami":           aws.ToString(instance.ImageId),
		"subnet_id":     aws.ToString(instance.SubnetId),
		"vpc_id":        aws.ToString(instance.VpcId),
		"private_ip":    aws.ToString(instance.PrivateIpAddress),
		"public_ip":     aws.ToString(instance.PublicIpAddress),
	}
	return attrs, nil
}
