package helpers

import (
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

// IsObjectNotFoundError centralizes detection of "not found" responses returned by the ZPA APIs.
// The Zscaler SDK exposes ErrorResponse.IsObjectNotFound(), but some endpoints return custom
// reason strings (for example "resource.type.not.found") even when the HTTP status is 400.
// This helper normalizes those variations so resources can share the same pattern.
func IsObjectNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	respErr, ok := err.(*errorx.ErrorResponse)
	if !ok || respErr == nil {
		return false
	}

	if respErr.IsObjectNotFound() {
		return true
	}

	if respErr.Parsed == nil {
		return false
	}

	if isKnownNotFoundIdentifier(respErr.Parsed.ID) {
		return true
	}

	if isKnownNotFoundIdentifier(respErr.Parsed.Reason) {
		return true
	}

	return false
}

func isKnownNotFoundIdentifier(value string) bool {
	if value == "" {
		return false
	}

	lower := strings.ToLower(value)

	switch lower {
	case "resource.not.found", "resource.type.not.found":
		return true
	}

	if strings.Contains(lower, "resource type not found") {
		return true
	}

	return false
}
