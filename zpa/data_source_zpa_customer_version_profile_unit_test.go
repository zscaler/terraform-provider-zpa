package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customerversionprofile"
)

// Version profiles are looked up by name only.
func TestDataSourceCustomerVersionProfile(t *testing.T) {
	testDataSource(t, dataSourceSpec{
		resource: dataSourceCustomerVersionProfile(),
		listPath: v1Path("visible", "versionProfiles"),
		items: []customerversionprofile.CustomerVersionProfile{
			{ID: "0", Name: "Default"},
			{ID: "2", Name: "New Release", Description: "Latest release", Versions: []customerversionprofile.Versions{{ID: "v1", Platform: "el8", Role: "zpn_assistant"}}},
		},
		name: "New Release",
		id:   "2",
		wantState: map[string]any{
			"description": "Latest release",
			"versions":    []any{map[string]any{"platform": "el8", "role": "zpn_assistant"}},
		},
	})
}
