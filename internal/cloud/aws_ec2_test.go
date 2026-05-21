package cloud

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// mockEC2Client implements a minimal EC2 API for testing fetchEC2Instance.
type mockEC2Client struct {
	output *ec2.DescribeInstancesOutput
	err    error
}

func (m *mockEC2Client) DescribeInstances(
	_ context.Context,
	_ *ec2.DescribeInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.DescribeInstancesOutput, error) {
	return m.output, m.err
}

// singleInstanceOutput builds a DescribeInstancesOutput with one instance.
func singleInstanceOutput(instanceID, instanceType, state, az string, tags []types.Tag) *ec2.DescribeInstancesOutput {
	return &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId:       aws.String(instanceID),
						InstanceType:     types.InstanceType(instanceType),
						Placement:        &types.Placement{AvailabilityZone: aws.String(az)},
						State:            &types.InstanceState{Name: types.InstanceStateName(state)},
						Tags:             tags,
					},
				},
			},
		},
	}
}

func TestFetchEC2Instance_Success(t *testing.T) {
	client := &mockEC2Client{
		output: singleInstanceOutput(
			"i-0abc123",
			"t3.micro",
			"running",
			"us-east-1a",
			[]types.Tag{{Key: aws.String("Name"), Value: aws.String("web-server")}},
		),
	}

	attrs, err := fetchEC2Instance(context.Background(), client, "i-0abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cases := []struct {
		key  string
		want string
	}{
		{"instance_type", "t3.micro"},
		{"availability_zone", "us-east-1a"},
		{"instance_state", "running"},
		{"tag_Name", "web-server"},
	}
	for _, tc := range cases {
		got, ok := attrs[tc.key]
		if !ok {
			t.Errorf("missing attribute %q", tc.key)
			continue
		}
		if got != tc.want {
			t.Errorf("attribute %q: got %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestFetchEC2Instance_NotFound(t *testing.T) {
	client := &mockEC2Client{
		output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{},
		},
	}

	_, err := fetchEC2Instance(context.Background(), client, "i-missing")
	if err == nil {
		t.Fatal("expected error for missing instance, got nil")
	}
}

func TestFetchEC2Instance_NoTags(t *testing.T) {
	client := &mockEC2Client{
		output: singleInstanceOutput("i-notags", "t2.small", "stopped", "us-west-2b", nil),
	}

	attrs, err := fetchEC2Instance(context.Background(), client, "i-notags")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attrs["instance_state"] != "stopped" {
		t.Errorf("expected instance_state=stopped, got %q", attrs["instance_state"])
	}
	// No tag_ keys should be present.
	for k := range attrs {
		if len(k) > 4 && k[:4] == "tag_" {
			t.Errorf("unexpected tag attribute %q for instance with no tags", k)
		}
	}
}

func TestExtractInstance_MissingFields(t *testing.T) {
	// Instance with nil placement and nil state should not panic.
	inst := types.Instance{
		InstanceId:   aws.String("i-bare"),
		InstanceType: types.InstanceTypeT3Micro,
	}

	attrs := instanceAttributes(inst)
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("expected instance_type=t3.micro, got %q", attrs["instance_type"])
	}
	// availability_zone and instance_state should be empty strings, not missing.
	if _, ok := attrs["availability_zone"]; !ok {
		t.Error("expected availability_zone key to be present even when nil")
	}
}
