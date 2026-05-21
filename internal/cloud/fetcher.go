package cloud

import (
	"context"
	"fmt"
)

// ResourceAttributes represents the live attributes of a cloud resource.
type ResourceAttributes map[string]interface{}

// Fetcher defines the interface for retrieving live cloud resource state.
type Fetcher interface {
	Fetch(ctx context.Context, resourceType, resourceID string) (ResourceAttributes, error)
}

// ErrResourceNotFound is returned when a resource cannot be located in the cloud provider.
type ErrResourceNotFound struct {
	ResourceType string
	ResourceID   string
}

func (e *ErrResourceNotFound) Error() string {
	return fmt.Sprintf("resource not found: type=%s id=%s", e.ResourceType, e.ResourceID)
}

// MockFetcher is an in-memory Fetcher implementation used for testing and local development.
type MockFetcher struct {
	// Resources maps "resourceType/resourceID" to its attributes.
	Resources map[string]ResourceAttributes
}

// NewMockFetcher creates a MockFetcher pre-populated with the given resource map.
func NewMockFetcher(resources map[string]ResourceAttributes) *MockFetcher {
	if resources == nil {
		resources = make(map[string]ResourceAttributes)
	}
	return &MockFetcher{Resources: resources}
}

// Fetch returns the attributes for the specified resource, or ErrResourceNotFound if absent.
func (m *MockFetcher) Fetch(_ context.Context, resourceType, resourceID string) (ResourceAttributes, error) {
	key := resourceType + "/" + resourceID
	attrs, ok := m.Resources[key]
	if !ok {
		return nil, &ErrResourceNotFound{ResourceType: resourceType, ResourceID: resourceID}
	}
	// Return a shallow copy to prevent mutation of the underlying map.
	copy := make(ResourceAttributes, len(attrs))
	for k, v := range attrs {
		copy[k] = v
	}
	return copy, nil
}
