package zpa

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The provider block bounds max_retries and request_timeout.
func TestProviderValidateSettings(t *testing.T) {
	cases := map[string]struct {
		config  map[string]any
		wantErr string
	}{
		"defaults":                 {map[string]any{}, ""},
		"max retries at limit":     {map[string]any{"max_retries": 100}, ""},
		"max retries over limit":   {map[string]any{"max_retries": 101}, "at most (100)"},
		"request timeout in range": {map[string]any{"request_timeout": 300}, ""},
		"request timeout too long": {map[string]any{"request_timeout": 301}, "in the range (0 - 300)"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := ""
			for _, d := range ZPAProvider().Validate(terraform.NewResourceConfigRaw(tc.config)) {
				if d.Severity == diag.Error {
					got = d.Summary
					break
				}
			}
			if tc.wantErr == "" && got != "" || tc.wantErr != "" && !strings.Contains(got, tc.wantErr) {
				t.Errorf("got %q, want %q", got, tc.wantErr)
			}
		})
	}
}
