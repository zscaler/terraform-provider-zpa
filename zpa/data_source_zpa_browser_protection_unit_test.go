package zpa

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

var browserProtectionProfiles = []browser_protection.BrowserProtection{
	{ID: "bp0", Name: "Default Browser Protection Legacy"},
	{ID: "bp1", Name: "Default Browser Protection", Description: "Monitors browsers", DefaultCSP: true},
}

// Browser protection profiles are looked up by name (there is no ID argument).
func TestDataSourceBrowserProtection(t *testing.T) {
	testDataSource(t, dataSourceSpec{
		resource:  dataSourceBrowserProtection(),
		listPath:  v1Path("browserProtectionProfile"),
		items:     browserProtectionProfiles,
		name:      "Default Browser Protection",
		id:        "bp1",
		wantState: map[string]any{"description": "Monitors browsers", "default_csp": true},
	})
}

// Without a name the tenant's default (first) profile is returned; the built-in
// profile has no ID, so one is derived from its name.
func TestDataSourceBrowserProtection_DefaultProfile(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("browserProtectionProfile"), common.SuccessResponse(common.ZPAList([]browser_protection.BrowserProtection{
		{Name: "Zs Recommended profile", DefaultCSP: true}, browserProtectionProfiles[1],
	})))

	d := fx.resourceData(dataSourceBrowserProtection(), "", map[string]any{})
	fx.run(dataSourceBrowserProtectionRead, d)
	if d.Id() == "" || d.Id() != generateShortID("Zs Recommended profile") {
		t.Errorf("id = %q, want the ID derived from the profile name", d.Id())
	}
	assertState(t, d, map[string]any{"name": "Zs Recommended profile"})
}
