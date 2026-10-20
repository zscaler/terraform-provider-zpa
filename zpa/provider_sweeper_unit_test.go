package zpa

import (
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// sweptTenant mocks a tenant holding two test objects (t1, t2) next to
// production objects (prod, and the default policy rule) and returns the
// distinct IDs the sweeper deleted. Every list endpoint returns the same objects.
func sweptTenant(t *testing.T, sweep func(*testClient) error) []string {
	t.Helper()
	fx := newUnitFixture(t)
	object := func(id, name string, emails ...string) map[string]any {
		return map[string]any{
			"id": id, "name": name, "config": map[string]any{"name": name},
			"emailIds": emails, "associationType": "CONNECTOR_GRP",
			// the SDK lists browser access, inspection and PRA segments by their sub-apps
			"clientlessApps": []any{map[string]any{"id": "c-" + id}},
			"inspectionApps": []any{map[string]any{"id": "i-" + id}},
			"praApps":        []any{map[string]any{"id": "p-" + id}},
		}
	}
	items := []map[string]any{
		object("t1", "tf-acc-test-app", "pra_user_1@acme.com"),
		object("t2", "tf-updated-app", "jdoe@acme.com"),
		object("prod", "Production App", "jdoe@acme.com"),
		object("default", "Global_Policy"),
	}
	var deleted []string
	respond := func(r *http.Request, _ []byte) common.MockResponse {
		if r.Method == "DELETE" {
			parts := strings.Split(r.URL.Path, "/")
			deleted = append(deleted, parts[len(parts)-1])
			return common.NoContentResponse()
		}
		if strings.Contains(r.URL.Path, "/cbiconfig/") {
			return common.SuccessResponse(items) // CBI lists are plain arrays
		}
		// a paged list envelope that also decodes as a single object (policy set lookups)
		return common.SuccessResponse(map[string]any{"id": "ps1", "name": "policy set", "list": items, "totalPages": 1})
	}
	for depth := 1; depth <= 14; depth++ {
		pattern := strings.Repeat("/*", depth)
		for _, method := range []string{"GET", "POST", "DELETE"} {
			fx.api.OnFunc(method, pattern, respond)
		}
	}
	if err := sweep(&testClient{sdkV3Client: fx.api.Service.Client}); err != nil {
		t.Fatalf("sweeper failed: %v", err)
	}
	// IDs deleted, once each (provisioning keys are listed per association type)
	unique := map[string]bool{}
	for _, id := range deleted {
		unique[id] = true
	}
	ids := make([]string, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// Sweepers run against a live tenant: they must delete the objects created by
// the acceptance tests (tf-acc-test-*, tf-updated-*) and nothing else.
func TestSweepersDeleteOnlyTestObjects(t *testing.T) {
	sweepers := map[string]func(*testClient) error{
		"app connector group":    sweepTestAppConnectorGroup,
		"application server":     sweepTestApplicationServer,
		"application segment":    sweepTestApplicationSegment,
		"browser access segment": sweepTestApplicationSegmentBA,
		"inspection segment":     sweepTestApplicationInspection,
		"pra segment":            sweepTestApplicationPRA,
		"inspection custom ctl":  sweepTestInspectionCustomControl,
		"inspection profile":     sweepTestInspectionProfile,
		"lss config":             sweepTestLSSConfigController,
		"provisioning key":       sweepTestProvisioningKey,
		"segment group":          sweepTestSegmentGroup,
		"server group":           sweepTestServerGroup,
		"service edge group":     sweepTestServiceEdgeGroup,
		"cbi banner":             sweepTestCBIBanner,
		"cbi external profile":   sweepTestCBIExternalProfile,
		"cbi certificate":        sweepTestCBICertificate,
		"ba certificate":         sweepTestBaCertificate,
		"pra credential":         sweepTestPRACredentialController,
		"pra console":            sweepTestPRAConsoleController,
		"pra portal":             sweepTestPRAPortalController,
		"tag group":              sweepTestTagGroup,
		"tag namespace":          sweepTestTagNamespace,
	}
	for name, sweep := range sweepers {
		t.Run(name, func(t *testing.T) {
			deleted := sweptTenant(t, sweep)
			if strings.Join(deleted, ",") != "t1,t2" {
				t.Errorf("deleted %v, want only the test objects [t1 t2]", deleted)
			}
		})
	}
}

// Policy rules are swept per policy type; the default rule and production rules stay.
func TestSweepersDeleteOnlyTestPolicyRules(t *testing.T) {
	deleted := sweptTenant(t, sweepTestAccessPolicyRuleByType)
	for _, id := range deleted {
		if id != "t1" && id != "t2" {
			t.Errorf("deleted %q; only test rules may be deleted", id)
		}
	}
	if len(deleted) == 0 {
		t.Error("no test rule was deleted")
	}
}

// Tag keys are swept within each namespace.
func TestSweepersDeleteOnlyTestTagKeys(t *testing.T) {
	deleted := sweptTenant(t, sweepTestTagKey)
	for _, id := range deleted {
		if id != "t1" && id != "t2" {
			t.Errorf("deleted %q; only test tag keys may be deleted", id)
		}
	}
	if len(deleted) == 0 {
		t.Error("no test tag key was deleted")
	}
}

// Approvals have no name: the ones created by the acceptance tests are
// assigned to pra_user_* addresses.
func TestSweepersDeleteOnlyTestApprovals(t *testing.T) {
	deleted := sweptTenant(t, sweepTestPRAPrivilegedApprovalController)
	if strings.Join(deleted, ",") != "t1" {
		t.Errorf("deleted %v, want only the approval for pra_user_1 [t1]", deleted)
	}
}

// An LSS entry without a config block must not stop the sweep.
func TestSweepersLSSConfigWithoutConfig(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.pathV2("lssConfig"), common.SuccessResponse(common.ZPAList([]map[string]any{
		{"id": "noconfig"},
		{"id": "t1", "config": map[string]any{"name": "tf-acc-test-lss"}},
	})))
	fx.on("DELETE", fx.pathV2("lssConfig", "t1"), common.NoContentResponse())

	if err := sweepTestLSSConfigController(&testClient{sdkV3Client: fx.api.Service.Client}); err != nil {
		t.Fatalf("sweeper failed: %v", err)
	}
	if fx.calls("DELETE", fx.pathV2("lssConfig", "t1")) != 1 {
		t.Error("the test LSS config was not deleted")
	}
}
