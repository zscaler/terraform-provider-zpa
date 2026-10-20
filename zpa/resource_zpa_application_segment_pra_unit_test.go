package zpa

import (
	"strings"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func praApp(domain, port, protocol string) map[string]any {
	app := map[string]any{
		"name":                 domain,
		"domain":               domain,
		"application_port":     port,
		"application_protocol": protocol,
		"app_types":            []any{"SECURE_REMOTE_ACCESS"},
		"enabled":              true,
	}
	if protocol == "RDP" {
		app["connection_security"] = "ANY"
	}
	return app
}

func praSegmentConfig(apps ...map[string]any) map[string]any {
	domains := []any{}
	configs := []any{}
	for _, a := range apps {
		domains = append(domains, a["domain"])
		configs = append(configs, a)
	}
	return map[string]any{
		"name":             "PRA servers",
		"enabled":          true,
		"domain_names":     domains,
		"segment_group_id": "sg1",
		"tcp_port_ranges":  []any{"22", "22", "3389", "3389"},
		"tcp_port_range": []any{
			map[string]any{"from": "22", "to": "22"}, map[string]any{"from": "3389", "to": "3389"},
		},
		"common_apps_dto": []any{map[string]any{"apps_config": configs}},
	}
}

func praSegmentResponse(apps ...applicationsegmentpra.PRAApps) applicationsegmentpra.AppSegmentPRA {
	domains := []string{}
	for _, a := range apps {
		domains = append(domains, a.Domain)
	}
	return applicationsegmentpra.AppSegmentPRA{
		ID: "pra1", Name: "PRA servers", Enabled: true, SegmentGroupID: "sg1", DomainNames: domains,
		TCPAppPortRange: []zpacommon.NetworkPorts{{From: "22", To: "22"}, {From: "3389", To: "3389"}},
		PRAApps:         apps,
	}
}

var (
	rdpApp = applicationsegmentpra.PRAApps{ID: "pa1", AppID: "pra1", Name: "rdp.acme.com", Domain: "rdp.acme.com", ApplicationPort: "3389", ApplicationProtocol: "RDP", ConnectionSecurity: "ANY", Enabled: true}
	sshApp = applicationsegmentpra.PRAApps{ID: "pa2", AppID: "pra1", Name: "ssh.acme.com", Domain: "ssh.acme.com", ApplicationPort: "22", ApplicationProtocol: "SSH", Enabled: true}
)

func TestResourceApplicationSegmentPRA_CRUD(t *testing.T) {
	testCRUD(t, crudSpec{
		resource:   resourceApplicationSegmentPRA(),
		config:     praSegmentConfig(praApp("rdp.acme.com", "3389", "RDP")),
		id:         "pra1",
		response:   praSegmentResponse(rdpApp),
		createPath: v1Path("application"),
		itemPath:   v1Path("application", "pra1"),
		wantBody: map[string]any{
			"segmentGroupId": "sg1",
			"commonAppsDto": map[string]any{"appsConfig": []any{map[string]any{
				"domain": "rdp.acme.com", "applicationPort": "3389", "applicationProtocol": "RDP",
				"appTypes": []any{"SECURE_REMOTE_ACCESS"},
			}}},
		},
		wantState: map[string]any{
			"common_apps_dto": []any{map[string]any{"apps_config": []any{map[string]any{
				"domain": "rdp.acme.com", "pra_app_id": "pa1",
			}}}},
		},
	})
}

// Removing an app from apps_config sends its ID in deletedPraApps and drops its domain.
func TestResourceApplicationSegmentPRA_UpdateDeletesRemovedApps(t *testing.T) {
	fx := newUnitFixture(t)
	fx.on("GET", fx.path("application", "pra1"), common.SuccessResponse(praSegmentResponse(rdpApp, sshApp)))
	fx.on("PUT", fx.path("application", "pra1"), common.NoContentResponse())

	state := praSegmentConfig(praApp("rdp.acme.com", "3389", "RDP"), praApp("ssh.acme.com", "22", "SSH"))
	config := praSegmentConfig(praApp("rdp.acme.com", "3389", "RDP"))
	d := fx.updateData(resourceApplicationSegmentPRA(), "pra1", state, config)
	fx.run(resourceApplicationSegmentPRAUpdate, d)

	assertSubset(t, "update request", fx.body("PUT", fx.path("application", "pra1")), map[string]any{
		"domainNames": []any{"rdp.acme.com"},
		"commonAppsDto": map[string]any{
			"deletedPraApps": []any{"pa2"},
			"appsConfig":     []any{map[string]any{"domain": "rdp.acme.com", "appId": "pra1", "praAppId": "pa1"}},
		},
	})
}

func TestResourceApplicationSegmentPRA_RDPRequiresConnectionSecurity(t *testing.T) {
	fx := newUnitFixture(t)
	app := praApp("rdp.acme.com", "3389", "RDP")
	delete(app, "connection_security")

	err := fx.diff(resourceApplicationSegmentPRA(), "", map[string]any{}, praSegmentConfig(app))
	if err == nil || !strings.Contains(err.Error(), "connection_security is required when application_protocol is RDP") {
		t.Fatalf("expected an RDP app without connection_security to be rejected at plan time, got %v", err)
	}
}
