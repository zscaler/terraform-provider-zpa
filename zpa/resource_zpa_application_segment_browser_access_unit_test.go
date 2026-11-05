package zpa

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// Every configurable clientless app field must reach the plan: Terraform diffs
// set elements by hash only, so a field left out of the hash plans no update.
func TestResourceApplicationSegmentBrowserAccess_ClientlessAppChangesArePlanned(t *testing.T) {
	r := resourceApplicationSegmentBrowserAccess()
	for field, value := range map[string]any{
		"local_domain":         "isolation1.acme.com",
		"enabled":              false,
		"trust_untrusted_cert": true,
	} {
		t.Run(field, func(t *testing.T) {
			fx := newUnitFixture(t)
			state := browserAccessConfig()
			app := clientlessApp("crm.acme.com", "443")
			app["local_domain"] = "isolation.acme.com"
			app["id"], app["app_id"] = "ca1", "ba1"
			state["clientless_apps"] = []any{app}

			config := browserAccessConfig()
			changed := clientlessApp("crm.acme.com", "443")
			changed["local_domain"] = "isolation.acme.com"
			changed[field] = value
			config["clientless_apps"] = []any{changed}

			if d := fx.updateData(r, "ba1", state, config); !d.HasChange("clientless_apps") {
				t.Fatalf("changing clientless_apps.%s planned no update", field)
			}
		})
	}
}

// A changed local_domain is sent in the update, with the IDs the API needs to
// match the clientless app (id) and its application segment (appId).
func TestResourceApplicationSegmentBrowserAccess_UpdateSendsLocalDomain(t *testing.T) {
	fx := newUnitFixture(t)
	resp := browserAccessResponse()
	resp.ClientlessApps[0].LocalDomain = "isolation.acme.com"
	fx.on("GET", fx.path("application", "ba1"), common.SuccessResponse(resp))
	fx.on("PUT", fx.path("application", "ba1"), common.NoContentResponse())

	state := browserAccessConfig()
	app := clientlessApp("crm.acme.com", "443")
	app["local_domain"] = "isolation.acme.com"
	state["clientless_apps"] = []any{app}
	config := browserAccessConfig()
	changed := clientlessApp("crm.acme.com", "443")
	changed["local_domain"] = "isolation1.acme.com"
	config["clientless_apps"] = []any{changed}

	d := fx.updateData(resourceApplicationSegmentBrowserAccess(), "ba1", state, config)
	fx.run(resourceApplicationSegmentBrowserAccessUpdate, d)

	assertSubset(t, "update request", fx.body("PUT", fx.path("application", "ba1")), map[string]any{
		"id": "ba1",
		"clientlessApps": []any{map[string]any{
			"id": "ca1", "appId": "ba1", "name": "crm.acme.com", "localDomain": "isolation1.acme.com",
		}},
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

// With a Zscaler-managed certificate, ZPA generates ext_id and ext_domain_name;
// the resource and the data source store both in clientless_apps.
func TestApplicationSegmentBrowserAccess_StoresGeneratedExternalURLParts(t *testing.T) {
	resp := &applicationsegmentbrowseraccess.BrowserAccess{
		ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{{
			Name: "app01.internal.acme.com", ApplicationPort: "443", ApplicationProtocol: "HTTPS",
			ExtLabel: "app01", ExtDomain: "acme.com", ExtID: "ab12", ExtDomainName: "acme-com.p.zsproxy.net",
		}},
	}
	for name, r := range map[string]*schema.Resource{
		"resource":    resourceApplicationSegmentBrowserAccess(),
		"data source": dataSourceApplicationSegmentBrowserAccess(),
	} {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{})
		if err := d.Set("clientless_apps", flattenBaClientlessApps(resp)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		apps, ok := d.Get("clientless_apps").([]any) // list in the data source
		if !ok {
			apps = d.Get("clientless_apps").(*schema.Set).List() // set in the resource
		}
		app := apps[0].(map[string]any)
		if app["ext_id"] != "ab12" || app["ext_domain_name"] != "acme-com.p.zsproxy.net" {
			t.Errorf("%s: ext_id = %q, ext_domain_name = %q", name, app["ext_id"], app["ext_domain_name"])
		}
	}
}
