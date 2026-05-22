package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type rdsDescribeAPI interface {
	DescribeDBInstances(ctx context.Context, params *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
}

func fetchRDSInstance(ctx context.Context, client rdsDescribeAPI, resourceID string) (map[string]string, error) {
	out, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(resourceID),
	})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrResourceNotFound{ResourceID: resourceID}
		}
		return nil, fmt.Errorf("describe rds instance %s: %w", resourceID, err)
	}

	if len(out.DBInstances) == 0 {
		return nil, ErrResourceNotFound{ResourceID: resourceID}
	}

	return rdsInstanceAttributes(out.DBInstances[0]), nil
}

func rdsInstanceAttributes(db interface{ GetDBInstanceClass() *string }) map[string]string {
	return rdsAttrsFromInstance(db)
}

func rdsAttrsFromInstance(db interface{}) map[string]string {
	type rdsDB interface {
		GetDBInstanceClass() *string
	}
	// Use the concrete type from the SDK
	type concreteDB struct {
		DBInstanceClass         *string
		Engine                  *string
		EngineVersion           *string
		DBInstanceStatus        *string
		MultiAZ                 bool
		StorageEncrypted        bool
		AllocatedStorage        int32
		DBInstanceIdentifier    *string
	}

	attrs := map[string]string{}

	switch v := db.(type) {
	case interface {
		GetFields() (id, class, engine, version, status *string, multiAZ, encrypted bool, storage int32)
	}:
		id, class, engine, version, status, multiAZ, encrypted, storage := v.GetFields()
		setAttr(attrs, "db_instance_identifier", id)
		setAttr(attrs, "instance_class", class)
		setAttr(attrs, "engine", engine)
		setAttr(attrs, "engine_version", version)
		setAttr(attrs, "status", status)
		attrs["multi_az"] = fmt.Sprintf("%t", multiAZ)
		attrs["storage_encrypted"] = fmt.Sprintf("%t", encrypted)
		attrs["allocated_storage"] = fmt.Sprintf("%d", storage)
	}

	_ = concreteDB{}
	return attrs
}

func setAttr(attrs map[string]string, key string, val *string) {
	if val != nil {
		attrs[key] = *val
	}
}
