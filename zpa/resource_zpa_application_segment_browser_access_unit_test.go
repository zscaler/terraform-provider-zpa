package zpa

import (
	"reflect"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
)

func clientlessApp(name, port string) map[string]any {
	return map[string]any{
		"name":                 name,
		"application_port":     port,
		"application_protocol": "HTTPS",
		"certificate_id":       "c1",
		"domain":               name,
		"enabled":              true,
	}
}

func browserAccessConfig() map[string]any {
	return map[string]any{
		"name":             "CRM portal",
		"enabled":          true,
		"domain_names":     []any{"crm.acme.com"},
		"segment_group_id": "sg1",
		"tcp_port_ranges":  []any{"443", "443"},
		// tcp_port_range is computed: state always holds the blocks Read echoed back
		"tcp_port_range":  []any{map[string]any{"from": "443", "to": "443"}},
		"clientless_apps": []any{clientlessApp("crm.acme.com", "443")},
	}
}

func browserAccessResponse() applicationsegmentbrowseraccess.BrowserAccess {
	return applicationsegmentbrowseraccess.BrowserAccess{
		ID: "ba1", Name: "CRM portal", Enabled: true, SegmentGroupID: "sg1",
		DomainNames: []string{"crm.acme.com"}, TCPPortRanges: []string{"443", "443"},
		ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{{
			ID: "ca1", AppID: "ba1", Name: "crm.acme.com", ApplicationPort: "443", ApplicationProtocol: "HTTPS",
			CertificateID: "c1", Domain: "crm.acme.com", Enabled: true,
		}},
	}
}

func TestResourceApplicationSegmentBrowserAccess_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceApplicationSegmentBrowserAccess(),
		config:     browserAccessConfig(),
		id:         "ba1",
		response:   browserAccessResponse(),
		createPath: v1Path("application"),
		itemPath:   v1Path("application", "ba1"),
		wantBody: map[string]any{
			"name": "CRM portal", "segmentGroupId": "sg1", "domainNames": []any{"crm.acme.com"},
			"tcpPortRanges": []any{"443", "443"},
			"clientlessApps": []any{map[string]any{
				"name": "crm.acme.com", "applicationPort": "443", "applicationProtocol": "HTTPS", "certificateId": "c1",
			}},
		},
		wantState: map[string]any{
			"segment_group_id": "sg1",
			"clientless_apps":  []any{map[string]any{"name": "crm.acme.com", "id": "ca1", "application_port": "443"}},
		},
	})
}

// An update keeps the IDs of existing clientless apps, matched by name.
func TestResourceApplicationSegmentBrowserAccess_UpdateKeepsClientlessAppIDs(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "ba1"), common.SuccessResponse(browserAccessResponse()))
	fx.on("PUT", fx.path("application", "ba1"), common.NoContentResponse())

	config := browserAccessConfig()
	config["description"] = "updated"
	d := fx.updateData(resourceApplicationSegmentBrowserAccess(), "ba1", browserAccessConfig(), config)
	fx.run(resourceApplicationSegmentBrowserAccessUpdate, d)

	assertSubset(t, "update request", fx.body("PUT", fx.path("application", "ba1")), map[string]any{
		"id":             "ba1",
		"clientlessApps": []any{map[string]any{"id": "ca1", "appId": "ba1", "name": "crm.acme.com"}},
	})
}

// A clientless app on a port the application does not have yet needs the port
// ranges written first: the API rejects the app otherwise (invalid.clientless.port).
func TestResourceApplicationSegmentBrowserAccess_NewPortWrittenBeforeClientlessApp(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "ba1"), common.SuccessResponse(browserAccessResponse()))
	fx.on("PUT", fx.path("application", "ba1"), common.NoContentResponse())

	config := browserAccessConfig()
	config["tcp_port_ranges"] = []any{"443", "443", "8443", "8443"}
	config["clientless_apps"] = []any{clientlessApp("crm.acme.com", "443"), clientlessApp("admin.acme.com", "8443")}
	d := fx.updateData(resourceApplicationSegmentBrowserAccess(), "ba1", browserAccessConfig(), config)
	fx.run(resourceApplicationSegmentBrowserAccessUpdate, d)

	puts := fx.requestBodies("PUT", fx.path("application", "ba1"))
	if len(puts) != 2 {
		t.Fatalf("expected a port pre-write and the full update (2 PUTs), got %d", len(puts))
	}
	// 1st PUT: the new port range, with only the clientless apps that already exist
	if got := toStrings(puts[0]["tcpPortRanges"]); !reflect.DeepEqual(got, []string{"443", "443", "8443", "8443"}) {
		t.Errorf("pre-write tcpPortRanges = %v", got)
	}
	if apps, _ := puts[0]["clientlessApps"].([]any); len(apps) != 1 {
		t.Errorf("pre-write must keep only the existing clientless app, got %v", puts[0]["clientlessApps"])
	}
	// 2nd PUT: the full desired state
	if apps, _ := puts[1]["clientlessApps"].([]any); len(apps) != 2 {
		t.Errorf("final update must include both clientless apps, got %v", puts[1]["clientlessApps"])
	}
}
