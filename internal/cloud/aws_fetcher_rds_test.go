package cloud

import (
	"context"
	"testing"
)

func TestAWSFetcher_RDSInstance_UnsupportedType(t *testing.T) {
	f := NewAWSFetcher()
	_, err := f.Fetch(context.Background(), "aws_unsupported_rds", "some-id")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestMockFetcher_SimulatesRDSInstance(t *testing.T) {
	m := NewMockFetcher()
	m.Resources[resourceTypeRDS+"/mydb"] = map[string]string{
		"db_instance_identifier": "mydb",
		"instance_class":         "db.t3.micro",
		"engine":                 "mysql",
		"engine_version":         "8.0.32",
		"status":                 "available",
		"multi_az":               "false",
		"storage_encrypted":      "true",
		"allocated_storage":      "100",
	}

	attrs, err := m.Fetch(context.Background(), resourceTypeRDS, "mydb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["engine"] != "mysql" {
		t.Errorf("expected engine=mysql, got %s", attrs["engine"])
	}
	if attrs["storage_encrypted"] != "true" {
		t.Errorf("expected storage_encrypted=true, got %s", attrs["storage_encrypted"])
	}
}

func TestMockFetcher_MissingRDSInstanceReturnsNotFound(t *testing.T) {
	m := NewMockFetcher()
	_, err := m.Fetch(context.Background(), resourceTypeRDS, "ghost-db")
	if err == nil {
		t.Fatal("expected ErrResourceNotFound")
	}
	var nf ErrResourceNotFound
	if !isErrNotFound(err, &nf) {
		t.Errorf("expected ErrResourceNotFound, got %T", err)
	}
}

func isErrNotFound(err error, target *ErrResourceNotFound) bool {
	e, ok := err.(ErrResourceNotFound)
	if ok {
		*target = e
	}
	return ok
}
