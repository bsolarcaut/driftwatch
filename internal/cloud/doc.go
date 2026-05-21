// Package cloud provides the Fetcher interface and implementations for
// retrieving live cloud resource attributes. Fetcher implementations are
// consumed by the drift detector to compare live state against Terraform
// state snapshots.
//
// The MockFetcher is provided for use in tests and local development
// scenarios where real cloud credentials are unavailable.
package cloud
