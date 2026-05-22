package cloud

import "strings"

// notFoundError is a test-only sentinel used in unit tests to simulate
// AWS "NoSuchEntity" / "NotFound" responses without importing real SDK types.
type notFoundError struct{}

func (e *notFoundError) Error() string { return "NoSuchEntityException: not found" }

// isNotFound returns true when the error message contains common AWS
// not-found indicators. Real AWS SDK errors (e.g. *types.NoSuchEntityException)
// also satisfy this because their Error() strings contain these substrings.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "nosuchentity") ||
		strings.Contains(msg, "notfound") ||
		strings.Contains(msg, "no such") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "404")
}
