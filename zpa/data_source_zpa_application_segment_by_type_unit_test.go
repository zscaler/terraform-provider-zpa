package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbytype"
)

// The search is not an exact match: the segment with exactly the given name wins.
func TestDataSourceApplicationSegmentByType(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "getAppsByType"), common.SuccessResponse(common.ZPAList([]applicationsegmentbytype.AppSegmentBaseAppDto{
		{ID: "ba0", AppID: "app0", Name: "crm-legacy"},
		{ID: "ba1", AppID: "app1", Name: "CRM", Domain: "crm.acme.com", ApplicationPort: "443", ApplicationProtocol: "HTTPS"},
	})))

	d := fx.resourceData(dataSourceApplicationSegmentByType(), "", map[string]any{"application_type": "BROWSER_ACCESS", "name": "crm"})
	fx.run(dataSourceApplicationSegmentByTypeRead, d)
	if d.Id() != "ba1" {
		t.Errorf("id = %q, want ba1 (the exact name match)", d.Id())
	}
	assertState(t, d, map[string]any{"app_id": "app1", "application_port": "443"})
	reqs := fx.api.Server.AllRequests()
	if len(reqs) == 0 || !strings.Contains(reqs[0].Query, "applicationType=BROWSER_ACCESS") {
		t.Errorf("the application type must be passed as a query parameter, got %v", reqs)
	}
}

func TestDataSourceApplicationSegmentByType_InvalidTypeRejected(t *testing.T) {
	fx := newUnitFixture(t)
	d := fx.resourceData(dataSourceApplicationSegmentByType(), "", map[string]any{"application_type": "UNKNOWN", "name": "crm"})
	if msg := fx.runErr(dataSourceApplicationSegmentByTypeRead, d); !strings.Contains(msg, "invalid application_type") {
		t.Fatalf("expected an invalid application_type to be rejected, got %q", msg)
	}
	if n := len(fx.api.Server.AllRequests()); n != 0 {
		t.Errorf("no API call may be made, got %d", n)
	}
}
