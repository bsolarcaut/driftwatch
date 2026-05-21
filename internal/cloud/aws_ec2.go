package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// ec2Client defines the subset of the EC2 API used for drift detection.
// Using an interface here allows for easy mocking in tests.
type ec2Client interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// fetchEC2Instance retrieves live attributes for an aws_instance resource
// identified by the given resource ID (instance ID).
func fetchEC2Instance(ctx context.Context, client ec2Client, resourceID string) (map[string]string, error) {
	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{resourceID},
	})
	if err != nil {
		return nil, fmt.Errorf("ec2 DescribeInstances %s: %w", resourceID, err)
	}

	instance, err := extractInstance(out, resourceID)
	if err != nil {
		return nil, err
	}

	return instanceAttributes(instance), nil
}

// extractInstance pulls a single Instance out of a DescribeInstancesOutput,
// returning ErrResourceNotFound when no matching instance exists.
func extractInstance(out *ec2.DescribeInstancesOutput, resourceID string) (types.Instance, error) {
	for _, reservation := range out.Reservations {
		for _, inst := range reservation.Instances {
			if aws.ToString(inst.InstanceId) == resourceID {
				return inst, nil
			}
		}
	}
	return types.Instance{}, ErrResourceNotFound{ResourceID: resourceID}
}

// instanceAttributes converts the fields we care about on an EC2 instance into
// a flat string map that can be compared against Terraform state attributes.
func instanceAttributes(inst types.Instance) map[string]string {
	attrs := map[string]string{
		"instance_type": string(inst.InstanceType),
		"ami":           aws.ToString(inst.ImageId),
		"instance_state": string(inst.State.Name),
		"subnet_id":      aws.ToString(inst.SubnetId),
		"vpc_id":         aws.ToString(inst.VpcId),
		"private_ip":     aws.ToString(inst.PrivateIpAddress),
		"public_ip":      aws.ToString(inst.PublicIpAddress),
		"key_name":       aws.ToString(inst.KeyName),
	}

	// Flatten tags into "tags.<key>" entries so they match Terraform's
	// tags map representation in state.
	for _, tag := range inst.Tags {
		if tag.Key != nil && tag.Value != nil {
			attrs["tags."+*tag.Key] = *tag.Value
		}
	}

	// Remove empty values to avoid false-positive diffs against state
	// attributes that Terraform omits when unset.
	for k, v := range attrs {
		if v == "" {
			delete(attrs, k)
		}
	}

	return attrs
}
