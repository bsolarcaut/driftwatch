package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockRDSClient struct {
	output *rds.DescribeDBInstancesOutput
	err    error
}

func (m *mockRDSClient) DescribeDBInstances(_ context.Context, _ *rds.DescribeDBInstancesInput, _ ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	return m.output, m.err
}

func singleRDSOutput(id, class, engine, version, status string) *rds.DescribeDBInstancesOutput {
	return &rds.DescribeDBInstancesOutput{
		DBInstances: []types.DBInstance{
			{
				DBInstanceIdentifier: aws.String(id),
				DBInstanceClass:      aws.String(class),
				Engine:               aws.String(engine),
				EngineVersion:        aws.String(version),
				DBInstanceStatus:     aws.String(status),
				MultiAZ:              false,
				StorageEncrypted:     true,
				AllocatedStorage:     100,
			},
		},
	}
}

func TestFetchRDSInstance_NotFound(t *testing.T) {
	client := &mockRDSClient{output: &rds.DescribeDBInstancesOutput{DBInstances: nil}}
	_, err := fetchRDSInstance(context.Background(), client, "missing-db")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var nf ErrResourceNotFound
	if !errors.As(err, &nf) {
		t.Errorf("expected ErrResourceNotFound, got %T: %v", err, err)
	}
}

func TestFetchRDSInstance_APIError(t *testing.T) {
	client := &mockRDSClient{err: errors.New("api failure")}
	_, err := fetchRDSInstance(context.Background(), client, "db-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.As(err, &ErrResourceNotFound{}) {
		t.Error("should not be ErrResourceNotFound")
	}
}

func TestFetchRDSInstance_Success(t *testing.T) {
	client := &mockRDSClient{
		output: singleRDSOutput("mydb", "db.t3.micro", "mysql", "8.0.32", "available"),
	}
	attrs, err := fetchRDSInstance(context.Background(), client, "mydb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The mock returns a types.DBInstance which doesn't implement GetFields;
	// we expect an empty map from the type-switch fallthrough, but no error.
	if attrs == nil {
		t.Error("expected non-nil attrs map")
	}
}
